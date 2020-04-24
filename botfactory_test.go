package main

import (
	"reflect"
	"testing"
	"time"
)

// Random test ;)  It needs to be refined
func TestBotfactory(t *testing.T) {
	type tBot struct {
		target string
		tick   time.Duration
	}
	cases := []struct {
		have tConfig
		vant []tBot
	}{
		{tConfig{"string 1", 10, map[string]tHost{"host1": {"10.1.1.1", []string{"ping", "http:8080"}, 1}, "host2": {"ippo.if.ua", []string{"http", "http:8888", "https:443"}, 3}, "host3": {"mon.gov.ua", []string{"https"}, 2}}},
			[]tBot{{"10.1.1.1:1", 1000000000}, {"10.1.1.1:8080", 1000000000}, {"ippo.if.ua:80", 3000000000}, {"ippo.if.ua:8888", 3000000000}, {"ippo.if.ua:443", 3000000000}, {"mon.gov.ua:443", 2000000000}}},
	}

	for _, c := range cases {
		brigade := bootFactory(c.have, nil)
		got := make([]tBot, len(brigade))
		for i, s := range brigade {
			//fmt.Println(i)
			//got = append(got, tBot{s.getTarget(), s.getTick()})
			got[i] = tBot{s.getTarget(), s.getTick()}
		}
		if !reflect.DeepEqual(got, c.vant) {
			t.Errorf("bootFactory(%v) == %v, vant %v", c.have, got, c.vant)
		}
	}
}
