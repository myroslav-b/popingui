package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func errorToStr(e error) string {
	if e == nil {
		return "nil"
	}
	return e.Error()
}

func TestParseConfig(t *testing.T) {
	cases := []struct {
		have     string
		wantConf tConfig
		wantErr  error
	}{
		{
			`
# TOML config file

field1 = "string 1"
runtime = 10

[hosts.host1]
address = "10.1.1.1"
actions = ["ping","http:8080"]
step = 1

[hosts.host2]
address = "ippo.if.ua"
actions = ["http", "http:8888", "https:443"]
step = 3

[hosts.host3]
address = "mon.gov.ua"
actions = ["https"]
step = 2
			`,
			tConfig{"string 1", 10, map[string]tHost{"host1": {"10.1.1.1", []string{"ping", "http:8080"}, 1}, "host2": {"ippo.if.ua", []string{"http", "http:8888", "https:443"}, 3}, "host3": {"mon.gov.ua", []string{"https"}, 2}}}, nil},
		{
			`
field1 = "field 1"
runtime = 0

[hosts.map1]
address = "1.1.1.1"
actions = ["ping","http:8080","https:4443"]
step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "https:4443"}, 1}}}, nil},
		{
			`
	field1 = "field 1"
	runtime = 0
	
	[hosts.map1]
	address = "1.1.1.1"
	step = 1
	actions = ["ping","http:8080","https:4443"]
	
	[hosts.map2]
	address = "127.0.0.1"
	actions = ["http", "https"]
	step = 100
	
	[hosts.map3]
	address = "255.255.255.255"
	actions = ["ping"]
	step = 10
				`, tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "https:4443"}, 1}, "map2": {"127.0.0.1", []string{"http", "https"}, 100}, "map3": {"255.255.255.255", []string{"ping"}, 10}}}, nil},
		{``,
			tConfig{}, errors.New("Host list is empty")},
		{`Any text`,
			tConfig{}, errors.New("Error decoding configuration file")},
		{
			`
field1 = "field 1"
runtime = 0

[hosts.map1]
address = "1.10.100.1000"
actions = ["ping","http:8080","https:4443"]
step = 1
		`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.10.100.1000", []string{"ping", "http:8080", "https:4443"}, 1}}}, errors.New("Invalid Address: 1.10.100.1000")},
		{
			`
	field1 = "field 1"
	runtime = 0
	
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:8080","httpss:4443"]
	step = 1
				`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "httpss:4443"}, 1}}}, errors.New("Invalid Action: action is not supported")},
		{
			`
	field1 = "field 1"
	runtime = -1
	
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:8080","https:4443"]
	step = 1
				`,
			tConfig{"field 1", -1, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "https:4443"}, 1}}}, errors.New("Invalid Runtime")},
		{
			`
	field1 = "field 1"
	runtime = 0
	
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping:80","http:8080","https:4443"]
	step = 1
				`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping:80", "http:8080", "https:4443"}, 1}}}, errors.New("Invalid Action: wrong structure")},
		{
			`
	field1 = "field 1"
	runtime = 0
		
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:8080:80","https:4443"]
	step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080:80", "https:4443"}, 1}}}, errors.New("Invalid Action: wrong structure")},
		{
			`
	field1 = "field 1"
	runtime = 0
			
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:8080","https:4443:443"]
	step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "https:4443:443"}, 1}}}, errors.New("Invalid Action: wrong structure")},
		{
			`
	field1 = "field 1"
	runtime = 0
			
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:8080","https:abrakadabra"]
	step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "https:abrakadabra"}, 1}}}, errors.New("Invalid Action: wrong structure")},
		{
			`
	field1 = "field 1"
	runtime = 0
				
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["pingtogo","http:8080","https:4443"]
	step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"pingtogo", "http:8080", "https:4443"}, 1}}}, errors.New("Invalid Action: action is not supported")},
		{
			`
	field1 = "field 1"
	runtime = 0
	
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:port","https:4443"]
	step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:port", "https:4443"}, 1}}}, errors.New("Invalid Action: wrong structure")},
		{
			`
	field1 = "field 1"
	runtime = 0
		
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:8080","https:port"]
	step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "https:port"}, 1}}}, errors.New("Invalid Action: wrong structure")},
		{
			`
	field1 = "field 1"
	runtime = 0
			
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:99999","https:4443"]
	step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:99999", "https:4443"}, 1}}}, errors.New("Invalid Action: wrong structure")},
		{
			`
	field1 = "field 1"
	runtime = 0
				
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:8080","https:99999"]
	step = 1
			`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "https:99999"}, 1}}}, errors.New("Invalid Action: wrong structure")},
		{
			`
	field1 = "field 1"
	runtime = 0
	
	[hosts.map1]
	address = "1.1.1.1"
	actions = []
	step = 1
				`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{}, 1}}}, errors.New("Invalid Action: no actions")},
		{
			`
	field1 = "field 1"
	runtime = 0
	
	[hosts.map1]
	address = "1.1.1.1"
	actions = ["ping","http:8080","https:4443"]
	step = -1
				`,
			tConfig{"field 1", 0, map[string]tHost{"map1": {"1.1.1.1", []string{"ping", "http:8080", "https:4443"}, -1}}}, errors.New("Invalid Ports: Step error")},
	}

	for _, c := range cases {
		gotConf, gotErr := parseConfig(strings.NewReader(c.have))
		if (!reflect.DeepEqual(gotConf, c.wantConf)) || (errorToStr(gotErr) != errorToStr(c.wantErr)) {
			t.Errorf("TestParseConfig for TOML \n %v \n\n gives the result \n\n %v \n\n and \n\n %v \n\n want \n\n %v \n\n and \n\n %v \n\n", c.have, gotConf, gotErr, c.wantConf, c.wantErr)
		}
	}
}
