package main

import (
	"flag"
	"github.com/prometheus-community/pro-bing"
	"log"
	"os/exec"
	"time"
)

func IsOnline(host string, pings int, timeout int) bool {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		panic(err)
	}

	pinger.Timeout = time.Duration(timeout)
	pinger.Count = pings
	err = pinger.Run()
	if err != nil {
		log.Println("Ping fails", err)
		// ping fails, network is down
		return false
	}

	// ping is ok
	return true
}

func PowerOff() {
	log.Println("Shutting down the system")
	err := exec.Command("shutdown", "-h", "now").Run()
	if err != nil {
		log.Println("Error shutting down computer:", err)
	}
}

func main() {

	failuresThreshold := flag.Int("f", 5, "Number of fails, after which triggering shutdown")
	interval := flag.Int("t", 60, "Interval between test execution in seconds")
	pings := flag.Int("n", 3, "Number of pings to do in each test")
	timeout := flag.Int("T", 3, "Seconds after a ping goes in timeout")
	host := flag.String("h", "1.1.1.1", "Host to ping, can be ip address or hostname")
	flag.Parse()

	consecutiveFailures := 0

	for {
		if IsOnline(*host, *pings, *timeout) {
			if consecutiveFailures > 0 {
				// internet comes back, resetting the timer
				consecutiveFailures = 0
			}
		} else {
			// internet is down, incrementing the count
			consecutiveFailures++

			// failure thresholds reached, triggering shutdown
			if consecutiveFailures >= *failuresThreshold {
				PowerOff()
			}
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
