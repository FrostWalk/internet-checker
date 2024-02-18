# Network Checker

Is a simple Go program that monitors the availability of the connection at specified intervals,
in order to make the check it pings a provided address (e.g., 1.1.1.1) to determine whether the network
is up or down. If the ping fails for a specified number checks, the program will execute the command
provided as argument.

## Features

- Monitors network connectivity using ICMP ping.
- Configurable threshold for consecutive failures before taking action.
- Configurable time interval between tests.
- Configurable time interval after a ping goes in timeout.
- Configurable number of pings to do in each test.
- Configurable command to execute when the network is down.
- Systemd service integration for automatic startup and management.

## Why?

The Network Checker was created to address a specific scenario where I use an old PC as a home server
connected to a UPS which does not provide interface for querying its power status. In this setup, the
router, is situated far from the server and cannot be connected to the UPS.

By continuously pinging a defined host, the program monitors the network connection. If the ping fails,
and the connection doesn't recover within a specified time frame, it indicates a power outage
ant the program takes proactive action by shutting down the PC just a few moments before the battery run out.
This approach helps ensure that the PC is powered off in the event of an extended power outage.
When the power comes back, there is a raspberry pi (which by design boot as soon as it's connected to power)
that when turned on checks that the server is up and if not, uses wake on lan to turn it on.
<br>**I know it is not the best possible solution, but it works**

## Usage

```
./network-checker
  -S int
        Seconds after a ping goes in timeout (default 2)
        
  -c list of strings (separated by ,)
        Command with args to execute when the network is down
        Example: /sbin/shutdown,-h,now
         
  -f int
        Number of fails, after which executing the command (default 5)
        
  -h string
        Host to ping, can be ip address or hostname (default "1.1.1.1")
        
  -p int
        Number of pings to do in each test (default 4)
        
  -t int
        Interval between test execution in seconds (default 60)

```

### As systemd service

- Compile the program with ``go build -ldflags="-w -s -buildid=" -trimpath -o network-checker main.go``
- Copy the executable to ``/usr/local/bin/``
- Copy the unit located inside the folder ``systemd`` to ``/etc/systemd/system/``
- Execute ``sudo systemctl daemon-reload``, ``sudo systemctl start network-checker``
  and ``sudo systemctl enable network-checker``

#### Note about unprivileged ping

The library used to ping the host (pro-bing) attempts to send an "unprivileged" ping via UDP. On Linux,
this must be enabled with the following sysctl command: ``sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"``
because this has to be done on every system boot, the best solution is to create a custom sysctl file
``/etc/sysctl.d/99-allow-unprivileged-ping.conf`` with the following content: ``net.ipv4.ping_group_range=0 2147483647``.
Then you need to load the new configuration with ``sudo sysctl -p /etc/sysctl.d/99-allow-unprivileged-ping.conf``.

in one command ``echo 'net.ipv4.ping_group_range=0 2147483647' | sudo tee /etc/sysctl.d/99-allow-unprivileged-ping.conf >/dev/null 
&& sudo sysctl -p /etc/sysctl.d/99-allow-unprivileged-ping.conf``

### Bash script

If you prefer, you can use the script that do all this for you
