package main

import (
	"flag"
	"github.com/mkmik/stringlist"
	"github.com/prometheus-community/pro-bing"
	"log"
	"os/exec"
	"time"
)

func IsNetworkUp(host string, pings int, timeout int) bool {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		panic(err)
	}

	pinger.Timeout = time.Duration(timeout) * time.Second
	pinger.Count = pings
	if err = pinger.Run(); err != nil {
		// fail to execute the ping, happens if network interface is down or other issues
		log.Println("Failed to execute the ping: ", err)
		return false
	}

	// considering the network up if the packet loss is less than 50%
	return pinger.Statistics().PacketLoss < 50.0
}

func ExecuteCommand(command []string) {
	log.Println("Executing command: ", command)

	c := command[0]
	args := command[1:]
	err := exec.Command(c, args...).Run()
	if err != nil {
		log.Println("Error while executing command:\n", err)
	}
}

func main() {
	failuresThreshold := flag.Int("f", 5, "Number of fails, after which executing the command")
	interval := flag.Int("t", 60, "Interval between test execution in seconds")
	pings := flag.Int("p", 4, "Number of pings to do in each test")
	timeout := flag.Int("S", 2, "Seconds after a ping goes in timeout")
	host := flag.String("h", "1.1.1.1", "Host to ping, can be ip address or hostname")

	var command stringlist.Value
	flag.Var(&command, "c", "Command with args, separated by a comma, to execute when the network is down\nDefault shutdown\nExample: /sbin/shutdown,-h,now")
	flag.Parse()

	if command == nil {
		command = []string{"/sbin/shutdown", "-h", "now"}
	}

	consecutiveFailures := 0

	for {
		if IsNetworkUp(*host, *pings, *timeout) {
			if consecutiveFailures > 0 {
				// connection comes back, resetting the timer
				consecutiveFailures = 0
			}
		} else {
			// connection is down, incrementing the count
			consecutiveFailures++

			// failure thresholds reached, triggering shutdown
			if consecutiveFailures >= *failuresThreshold {
				ExecuteCommand(command)
			}
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
