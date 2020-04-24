package main

import (
	"net"
	"strings"
	"time"
)

type tBrigade []tWorker

type tWorker interface {
	shoot() error
	init(net string, target string, tick time.Duration, chMessage chan<- tMessage, message tMessage)
	getTarget() string
	//setTick(tick time.Duration)
	getTick() time.Duration
}

type tPingBot struct {
	net     string
	target  string
	tick    time.Duration
	ch      chan<- tMessage
	message tMessage
}

func (pingBot *tPingBot) init(net string, target string, tick time.Duration, chMessage chan<- tMessage, message tMessage) {
	pingBot.net = net
	pingBot.target = target
	pingBot.tick = tick
	pingBot.ch = chMessage
	pingBot.message = message
}

func (pingBot *tPingBot) shoot() error {
	conn, err := net.Dial(pingBot.net, pingBot.target)
	pingBot.message.target = pingBot.target
	pingBot.message.err = err
	if err != nil {
		pingBot.ch <- pingBot.message
	} else {
		pingBot.ch <- pingBot.message
		conn.Close()
	}
	return err
	//fmt.Println("pingBot shoot")

}

func (pingBot *tPingBot) getTarget() string {
	return pingBot.target
}

/*func (pingBot *tPingBot) setTick(tick time.Duration) {
	pingBot.tick = tick
}*/

func (pingBot *tPingBot) getTick() time.Duration {
	return pingBot.tick
}

type tTCPBot struct {
	net     string
	target  string
	tick    time.Duration
	ch      chan<- tMessage
	message tMessage
}

func (tcpBot *tTCPBot) init(net string, target string, tick time.Duration, chMessage chan<- tMessage, message tMessage) {
	tcpBot.net = net
	tcpBot.target = target
	tcpBot.tick = tick
	tcpBot.ch = chMessage
	tcpBot.message = message
}

func (tcpBot *tTCPBot) shoot() error {
	conn, err := net.Dial(tcpBot.net, tcpBot.target)
	tcpBot.message.target = tcpBot.target
	tcpBot.message.err = err
	if err != nil {
		tcpBot.ch <- tcpBot.message
	} else {
		tcpBot.ch <- tcpBot.message
		conn.Close()
	}
	//fmt.Println("tcpBot shoot")
	return nil
}

func (tcpBot *tTCPBot) getTarget() string {
	return tcpBot.target
}

/*func (tcpBot *tTCPBot) setTick(tick time.Duration) {
	tcpBot.tick = tick
}*/

func (tcpBot *tTCPBot) getTick() time.Duration {
	return tcpBot.tick
}

func bootFactory(config tConfig, chMessage chan<- tMessage) tBrigade {
	brigade := make(tBrigade, 0)
	for _, host := range config.Hosts {
		for _, action := range host.Actions {
			split := strings.Split(action, ":")
			param := actionMap[split[0]].defParam()
			if len(split) == 2 {
				param = split[1]
			}
			target := strings.Join([]string{host.Address, param}, ":")
			tick := time.Second * time.Duration(host.Step)
			worker := actionMap[split[0]].newBot(target, tick, chMessage)
			//worker := actionMap[split[0]].newBot(actionMap[split[0]])
			brigade = append(brigade, worker)
		}
	}
	return brigade
}
