package main

import (
	"errors"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type tHost struct {
	Address string
	Actions []string
	Step    int
}
type tConfig struct {
	Field1  string
	Runtime int
	Hosts   map[string]tHost
}

type tMessage struct {
	text   string
	target string
	err    error
}

type tAction = struct {
	//net        string
	//defParam   string
	//message    tMessage
	validation func(st string) error
	newBot     func(target string, tick time.Duration, chMessage chan<- tMessage) tWorker
	defParam   func() string
}

type tActionMap = map[string]tAction

var actionMap = tActionMap{
	"ping": {
		//"ip",
		//"1",
		//tMessage{"Ping bell"},
		func(st string) error {
			split := strings.Split(st, ":")
			if len(split) > 1 { //|| split[0] != "ping" {
				return errors.New("Invalid Action: wrong structure")
			}
			return nil
		},
		func(target string, tick time.Duration, chMessage chan<- tMessage) tWorker {
			net := "ip"
			//defParam := "1"
			message := tMessage{"Ping bell", "", nil}
			worker := new(tPingBot)
			worker.init(net, target, tick, chMessage, message)
			return worker
		},
		func() string {
			return "1"
		},
	},
	"http": {
		//"tcp",
		//"80",
		//tMessage{"HTTP bell"},
		func(st string) error {
			split := strings.Split(st, ":")
			/*if split[0] != "http" {
				return errors.New("Invalid Action: wrong structure")
			}*/
			switch len(split) {
			case 1:
				return nil
			case 2:
				p, err := strconv.Atoi(split[1])
				if err != nil {
					return errors.New("Invalid Action: port is not a number")
				}
				if (p <= 0) || (p > 65535) {
					return errors.New("Invalid Action: port is not correct")
				}
				return nil
			default:
				return errors.New("Invalid Action: wrong structure")
			}
		},
		func(target string, tick time.Duration, chMessage chan<- tMessage) tWorker {
			net := "tcp"
			//defParam := "80"
			message := tMessage{"HTTP bell", "", nil}
			worker := new(tTCPBot)
			worker.init(net, target, tick, chMessage, message)
			return worker
		},
		func() string {
			return "80"
		},
	},
	"https": {
		//"tcp",
		//"443",
		//tMessage{"HTTPS bell"},
		func(st string) error {
			split := strings.Split(st, ":")
			/*if split[0] != "https" {
				return errors.New("Invalid Action: wrong structure")
			}*/
			switch len(split) {
			case 1:
				return nil
			case 2:
				p, err := strconv.Atoi(split[1])
				if err != nil {
					return errors.New("Invalid Action: port is not a number")
				}
				if (p <= 0) || (p > 65535) {
					return errors.New("Invalid Action: port is not correct")
				}
				return nil
			default:
				return errors.New("Invalid Action: wrong structure")
			}
		},
		func(target string, tick time.Duration, chMessage chan<- tMessage) tWorker {
			net := "tcp"
			//defParam := "443"
			message := tMessage{"HTTPS bell", "", nil}
			worker := new(tTCPBot)
			worker.init(net, target, tick, chMessage, message)
			return worker
		},
		func() string {
			return "443"
		},
	},
}

func matchIP(ip string) bool {
	//const cRegexpIP = `^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`
	const cRegexpIP = `^(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$`
	b, err := regexp.MatchString(cRegexpIP, ip)
	if !b || err != nil {
		return false
	}
	return true
}

func matchDN(dn string) bool {
	//const cRegexpDN = `^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\\.)+[A-Za-z]{2,6}$`
	//const cRegexpDN = `^(([A-Za-z0-9]\-*[A-Za-z0-9]*){1,63}\.)+[A-Za-z]{2,6}$`
	const cRegexpDN = `^(([A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9]|[A-Za-z0-9]){1,63}\.)+[A-Za-z]{2,6}$`
	b, err := regexp.MatchString(cRegexpDN, dn)
	if !b || err != nil {
		return false
	}
	return true
}

func matchAddress(address string) bool {
	return matchIP(address) || matchDN(address)
}

func parseConfig(r io.Reader) (tConfig, error) {
	var conf tConfig
	_, err := toml.DecodeReader(r, &conf)
	if err != nil {
		return conf, errors.New("Error decoding configuration file")
	}
	if (conf.Runtime < 0) || (conf.Runtime > math.MaxInt32) {
		return conf, errors.New("Invalid Runtime")
	}
	if len(conf.Hosts) == 0 {
		return conf, errors.New("Host list is empty")
	}
	for _, paramHost := range conf.Hosts {
		if !matchAddress(paramHost.Address) {
			return conf, errors.New("Invalid Address: " + paramHost.Address)
		}
		if len(paramHost.Actions) == 0 {
			return conf, errors.New("Invalid Action: no actions")
		}
		for _, action := range paramHost.Actions {
			actionSplit := strings.Split(action, ":")
			_, ok := actionMap[actionSplit[0]]
			if !ok {
				return conf, errors.New("Invalid Action: action is not supported")
			}
			err = actionMap[actionSplit[0]].validation(action)
			if err != nil {
				return conf, errors.New("Invalid Action: wrong structure")
			}
			/*if len(actionSplit) > 1 {
				p, err := strconv.Atoi(actionSplit[1])
				if err != nil {
					return conf, errors.New("Invalid Action: port is not a number")
				}
				if (p <= 0) || (p > 65535) {
					return conf, errors.New("Invalid Action: port is not correct")
				}
			}*/
		}
		if (paramHost.Step <= 0) || (paramHost.Step > math.MaxInt32) {
			return conf, errors.New("Invalid Ports: Step error")
		}
	}
	//fmt.Println(md, r, conf)
	return conf, nil
}

func normalizationConfig(conf tConfig) tConfig {
	if conf.Runtime == 0 {
		conf.Runtime = math.MaxInt32
	}
	return conf
}

func readConfig(nameConfigFile string) (tConfig, error) {
	var f *os.File
	f, err := os.Open(nameConfigFile)
	if err != nil {
		//fmt.Println("Config file not found (or not opening). Program stopped")
		return tConfig{}, errors.New("Config file not found (or not opening)")
	}
	defer f.Close()
	var config tConfig
	config, err = parseConfig(f)
	if err != nil {
		return config, errors.New("Config parsing error: " + error.Error(err))
	}
	////
	config = normalizationConfig(config)
	////
	return config, nil
}
