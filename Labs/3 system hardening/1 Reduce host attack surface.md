# Reduce attack surface OS

Install a service on your control node.

```bash
$ sudo apt install ftpd
```

## Using ss to find open ports
Find open ports that are listening.
```bash
$ ss -tulpan | grep LISTEN
```
tcp   LISTEN    0      128               0.0.0.0:21                    0.0.0.0:*

## Using lsof to find PID and user of app
```bash
$ sudo lsof -i :21
```
COMMAND   PID USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
inetd   36386 root    7u  IPv4 190629      0t0  TCP *:ftp (LISTEN)

## Find pull path to binary
```bash
$ ls -l /proc/PID#/exe  #found from output of command above
```
Kill and delete
```bash
$ sudo ls -l /proc/36386/exe
```
lrwxrwxrwx 1 root root 0 Jul 16 05:05 /proc/36386/exe -> /usr/sbin/inetd

```bash
$ sudo kill PID#
$ rm /usr/sbin/inetd
```
# Remove unneeded services
```bash
$ systemctl --type=service --state=running
```
```
  UNIT                        LOAD   ACTIVE SUB     DESCRIPTION
  crio.service                loaded active running Container Runtime Interface for OCI (CRI-O)
  cron.service                loaded active running Regular background program processing daemon
  dbus.service                loaded active running D-Bus System Message Bus
  getty@tty1.service          loaded active running Getty on tty1
  irqbalance.service          loaded active running irqbalance daemon
  kubelet.service             loaded active running kubelet: The Kubernetes Node Agent
  ModemManager.service        loaded active running Modem Manager
  multipathd.service          loaded active running Device-Mapper Multipath Device Controller
  networkd-dispatcher.service loaded active running Dispatcher daemon for systemd-networkd
  polkit.service              loaded active running Authorization Manager
  rsyslog.service             loaded active running System Logging Service
  snapd.service               loaded active running Snap Daemon
  ssh.service                 loaded active running OpenBSD Secure Shell server
  systemd-journald.service    loaded active running Journal Service
  systemd-logind.service      loaded active running User Login Management
  systemd-networkd.service    loaded active running Network Configuration
  systemd-resolved.service    loaded active running Network Name Resolution
  systemd-timedated.service   loaded active running Time & Date Service
  systemd-udevd.service       loaded active running Rule-based Manager for Device Events and Files
  udisks2.service             loaded active running Disk Manager
  user@1000.service           loaded active running User Manager for UID 1000
  vboxadd-service.service     loaded active running vboxadd-service.service
```
$ systemctl list-units --type=service --all

$ systemctl list-unit-files

$ sudo systemctl stop snapd

$ sudo systemctl disable snapd

$ systemctl --type=service --state=running
```
  UNIT                        LOAD   ACTIVE SUB     DESCRIPTION
  crio.service                loaded active running Container Runtime Interface for OCI (CRI-O)
  cron.service                loaded active running Regular background program processing daemon
  dbus.service                loaded active running D-Bus System Message Bus
  getty@tty1.service          loaded active running Getty on tty1
  irqbalance.service          loaded active running irqbalance daemon
  kubelet.service             loaded active running kubelet: The Kubernetes Node Agent
  ModemManager.service        loaded active running Modem Manager
  multipathd.service          loaded active running Device-Mapper Multipath Device Controller
  networkd-dispatcher.service loaded active running Dispatcher daemon for systemd-networkd
  polkit.service              loaded active running Authorization Manager
  rsyslog.service             loaded active running System Logging Service
  ssh.service                 loaded active running OpenBSD Secure Shell server
  systemd-journald.service    loaded active running Journal Service
  systemd-logind.service      loaded active running User Login Management
  systemd-networkd.service    loaded active running Network Configuration
  systemd-resolved.service    loaded active running Network Name Resolution
  systemd-udevd.service       loaded active running Rule-based Manager for Device Events and Files
  udisks2.service             loaded active running Disk Manager
  user@1000.service           loaded active running User Manager for UID 1000
  vboxadd-service.service     loaded active running vboxadd-service.service
```
```bash
$ apt list --installed | grep snapd
```
snapd/now 2.58+22.04 amd64 [installed,upgradable to: 2.58+22.04.1]
```bash
$ sudo apt remove --auto-remove snapd
```
```
Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
The following packages will be REMOVED:
  snapd squashfs-tools
0 upgraded, 0 newly installed, 2 to remove and 100 not upgraded.
After this operation, 103 MB disk space will be freed.
Do you want to continue? [Y/n] y
(Reading database ... 44393 files and directories currently installed.)
Removing snapd (2.58+22.04) ...
Stopping snap.lxd.activate.service
Stopping unit snap.lxd.activate.service
Waiting until unit snap.lxd.activate.service is stopped [attempt 1]
snap.lxd.activate.service is stopped.
Stopping snap.lxd.daemon.service
Stopping unit snap.lxd.daemon.service
Waiting until unit snap.lxd.daemon.service is stopped [attempt 1]
snap.lxd.daemon.service is stopped.
Stopping snap.lxd.user-daemon.service
Stopping unit snap.lxd.user-daemon.service
Waiting until unit snap.lxd.user-daemon.service is stopped [attempt 1]
snap.lxd.user-daemon.service is stopped.
Stopping snap.lxd.daemon.unix.socket
Stopping unit snap.lxd.daemon.unix.socket
Waiting until unit snap.lxd.daemon.unix.socket is stopped [attempt 1]
snap.lxd.daemon.unix.socket is stopped.
Stopping snap.lxd.user-daemon.unix.socket
Stopping unit snap.lxd.user-daemon.unix.socket
Waiting until unit snap.lxd.user-daemon.unix.socket is stopped [attempt 1]
snap.lxd.user-daemon.unix.socket is stopped.
Warning: Stopping snapd.service, but it can still be activated by:
  snapd.socket
Removing squashfs-tools (1:4.5-3build1) ...
Processing triggers for dbus (1.12.20-2ubuntu4.1) ...
Processing triggers for man-db (2.10.2-1) ...
```
