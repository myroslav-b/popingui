package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func work(ctx context.Context, wg *sync.WaitGroup, worker tWorker) {
	ticker := time.NewTicker(worker.getTick())
	defer ticker.Stop()
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			worker.shoot()
		}
	}
}

func speacker(chMessage <-chan tMessage) {
	for {
		select {
		case message := <-chMessage:
			fmt.Println(message.text, ":  ", message.target, ",  ", message.err)
		}
	}
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("This is not Linux")
		os.Exit(0)
	}

	pidMark, err := initPidMark(os.TempDir(), cPidFileName, os.Getpid())
	if err != nil {
		fmt.Println("Program stopped: " + error.Error(err))
		return
	}

	var flags tFlags
	check := checkFlags(&flags, os.Args[1:])
	switch check {
	case cStart:
		{
			if pidMark.getPid() == 0 {
				err = pidMark.new(os.Getpid())
				if err != nil {
					fmt.Println("Program stopped: " + error.Error(err))
					return
				}
			}

			nameConfigFile := flags.c
			config, err := readConfig(nameConfigFile)
			if err != nil {
				fmt.Println("Program stopped: " + error.Error(err))
				return
			}
			fmt.Println("Program started")

			chMessage := make(chan tMessage, 100)
			go speacker(chMessage)

			brigade := bootFactory(config, chMessage)

			wg := &sync.WaitGroup{}
			ctx, stopPing := context.WithCancel(context.Background())
			for _, worker := range brigade {
				fmt.Println("Worker started. Target: ", worker.getTarget())
				wg.Add(1)
				go work(ctx, wg, worker)
			}

			chanStopSignal := make(chan os.Signal)
			signal.Notify(chanStopSignal, syscall.SIGUSR1)
			chanReadSignal := make(chan os.Signal)
			signal.Notify(chanReadSignal, syscall.SIGUSR2)
			timer := time.NewTimer(time.Second * time.Duration(config.Runtime))
			for {
				select {
				case <-timer.C:
					err = pidMark.remove()
					if err != nil {
						fmt.Println(err)
					}
					stopPing()
					wg.Wait()
					fmt.Println("Time is up. Program stopped")
					return
				case <-chanStopSignal:
					err = pidMark.remove()
					if err != nil {
						fmt.Println(err)
					}
					stopPing()
					wg.Wait()
					fmt.Println("Program stopped")
					return
				case <-chanReadSignal:
					stopPing()
					wg.Wait()

					//nameConfigFile := flags.c
					config, err = readConfig(nameConfigFile)
					if err != nil {
						fmt.Println("Program stopped: " + error.Error(err))
						return
					}
					brigade = bootFactory(config, chMessage)
					ctx, stopPing = context.WithCancel(context.Background())
					for _, worker := range brigade {
						wg.Add(1)
						go work(ctx, wg, worker)
					}

					fmt.Println("Config file is read. Configuration changed")
				}
			}
		}
	case cStop:
		{
			if pidMark.getPid() == 0 {
				fmt.Println("PID file not found (or not opening). The program is probably not running")
				return
			}

			process, err := os.FindProcess(pidMark.getPid())
			if err != nil {
				fmt.Println("Program stopped: problem with process PID")
				return
			}
			err = process.Signal(syscall.SIGUSR1)
			if err != nil {
				fmt.Println("Program stopped: failed to send signal")
				return
			}

			fmt.Println("Stop signal sent")
		}
	case cVerify:
		{
			if pidMark.getPid() > 0 {
				fmt.Print("PID file exist. ")
				_, err := os.FindProcess(pidMark.getPid())
				if err == nil {
					fmt.Println("The program is probably works")
				} else {
					fmt.Println("Error: no process found")
				}
			} else {
				fmt.Println("PID file not found (or not opening). The program is probably not running")
			}
			return
		}
	case cClean:
		{

		}
	case cRead:
		{
			if pidMark.getPid() == 0 {
				fmt.Println("PID file not found (or not opening). The program is probably not running")
				return
			}

			process, err := os.FindProcess(pidMark.getPid())
			if err != nil {
				fmt.Println("Program stopped: problem with process PID")
				return
			}
			err = process.Signal(syscall.SIGUSR2)
			if err != nil {
				fmt.Println("Program stopped: failed to send signal")
				return
			}
			fmt.Println("Read signal sent")
		}
	default:
		{

		}
	}

}
