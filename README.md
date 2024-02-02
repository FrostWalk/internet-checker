# Internet Checker

Is a simple Go program that monitors the availability of the connection at specified intervals,
in order to check it ping a well-known address (e.g., 1.1.1.1) to determine whether the connection is up or
down. If the ping fails for a specified number checks, the program will shut down the system.
<br/>**Yes I know this is an ugly solution.**

## Features

- Monitors internet connectivity using ICMP ping.
- Configurable threshold for consecutive failures before taking action.
- Systemd service integration for automatic startup and management.

## Why?

The Internet Checker was created to address a specific scenario where a personal computer serves as a home server
connected to a UPS without a USB port for querying its power status. In this setup,
the router, is situated far from the server and cannot be connected to the UPS.

The motivation behind creating this program stems from the need to determine the power state of the UPS remotely.
By continuously pinging a defined host, the program monitors the internet connection. If the ping fails,
indicating a power outage, and the connection doesn't recover within a specified time frame, the program takes proactive
action by shutting down the PC just a few moments before the battery run out. This approach helps ensure that the PC is 
powered off in the event of an extended power outage. When the power comes back, there is a raspberry pi 
(which by design boot as soon as it's connected to power) that when turned on checks that the server is up and if not, 
uses wake on lan to turn it on

## Usage

```
./internet-checker

-T int
    Seconds after a ping goes in timeout (default 3)
-f int
    Number of fails, after which triggering shutdown (default 5)
-h string
    Host to ping, can be ip address or hostname (default "1.1.1.1")
-n int
    Number of pings to do in each test (default 3)
-t int
    Interval between test execution in seconds (default 60)
```

### As systemd service

- Compile the program with ``go build -o internet-checker main.go``
- Copy the executable to ``/usr/local/bin/``
- Copy the unit located inside the folder ``systemd`` to ``/etc/systemd/system/``
- Execute ``sudo systemctl daemon-reload`` ``sudo systemctl start internet-checker``
  and ``sudo systemctl enable internet-checker``

if you want to change the default parameters used to run the service, edit the file ``internet-checker.service``

### Bash script

If you prefer, you can use the script that do all this for you