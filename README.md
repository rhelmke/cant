# CAN't

**CAN't** is a multithreaded privacy proxy server for [ISO-11898](https://www.iso.org/standard/63648.html) Controller Area Networks that is mainly targeted at [ISOBUS](https://www.iso.org/standard/57556.html) applications.
Its purpose is to selectively block, perturbate, or round privacy sensitive data transmitted via CAN. The software can easily be configured using an HTML5 web interface.
Implementation details and the motivation behind this project can be found in:

*J. Bauer, R. Helmke, A. Bothe, and N. Aschenbruck, “CAN’t track us: Adaptable Privacy for ISOBUS Controller Area Networks”, Elsevier Computer Standards and Interfaces, 2019.*

Concept, implementation, and evaluation were topic of a bachelor's thesis at the [Distributed Systems Group](https://sys.cs.uos.de/), University of Osnabrück:

*R. Helmke, “Konzeptionierung und Implementierung eines Proxys für erhöhten Datenschutz in Controller Area Networks“,
2018, supervised by: Prof. Dr. Nils Aschenbruck, Prof. Dr. Michael Brinkmeier.*

## Dependencies

* `GOOS=linux GOARCH=amd64`
* Golang (`>=1.13`)
* go-bindata (`go get -u github.com/jteeuwen/go-bindata/...`)
* MariaDB (`apt install mysql-server`)
* npm/yarn
* [SocketCAN](https://www.mjmwired.net/kernel/Documentation/networking/can.txt)
* Packages: See `go.mod`

## Build

Make sure your `CXX` and `CC` variables are properly set:

```bash
export CXX=clang++ # or whatever you fancy
export CC=clang # or whatever you fancy
```

Then build the project:

```bash
git clone --recursive https://github.com/rhelmke/cant.git
cd cant
go generate # execute yarn and go-bindata to generate the webinterface
go build cant
```

## General Synopsis

```plain
Usage:
  cant [command]

Available Commands:
  help        Help about any command
  run         run the main components
  seed        Seed the database
  setup       Interactive cant setup
  version     Print version

Flags:
  -h, --help   help for cant

Use "cant [command] --help" for more information about a command.
```

```plain
Usage:
  cant run [command]

Available Commands:
  proxy       run the proxy component

Flags:
  -h, --help   help for run

Use "cant run [command] --help" for more information about a command.
```

## Manual Setup

### Database

```bash
apt install mysql-server # install mariadb on debian or ubuntu
mysql -uroot -p # login to mariadb and create empty database + user
mysql> create database cant;
mysql> create user 'cant'@'localhost' identified by '<PASSWORD>';
mysql> grant all privileges on cant.* to 'cant'@'localhost';
mysql> exit
```

### CAN Interface

**CAN't** needs two interfaces in order to work as man-in-the-middle between ECUs.
For testing purposes, you might want to create two virtual interfaces. To do so, add following lines `/etc/network/interfaces`:

```bash
auto vcan0
iface vcan0 inet manual
    bitrate 250000 # ISOBUS uses a nominal bitrate of 250kbit/s
    pre-up /sbin/ip link add dev $IFACE type vcan
    post-up /sbin/ip link set $IFACE txqueuelen 1000

auto vcan1
iface vcan1 inet manual
    bitrate 250000 # ISOBUS uses a nominal bitrate of 250kbit/s
    pre-up /sbin/ip link add dev $IFACE type vcan
    post-up /sbin/ip link set $IFACE txqueuelen 1000
```

Also, be sure that you loaded all needed kernel modules:

```bash
modprobe can
modprobe vcan
```

You can then execute `ifup vcan0` or `ifup vcan1` to bring the interfaces up.
You can read and write data using the `can-utils` package from your OS's repositories.

### Seeding

1. Make sure you did everything explained in Section "Database" and "CAN Interfaces".
2. `./cant setup`, this will guide you through a cli-based installation routine. You need a properly functioning MariaDB containing an empty database for **CAN't**.
3. `./cant seed -f <pgn and spn data>.csv -t spnpgn`, this will seed the database with all known PGN's and SPN's.
4. `./cant seed -t filter`, this will inject all implemented and compiled filters into the database.

## Where to get SPN and PGN Data

The [VDMA](https://www.isobus.net/isobus/) maintains a growing list of known PGN's and SPN's. Please use the [csv-formatted database dump](https://www.isobus.net/isobus/attachments/isoExport_csv.zip) and extract `SPNs and PGNs.csv`. This file can be used to seed the database during setup.

## Run the proxy

```bash
./cant run proxy
```

The proxy will expose its webinterface to a port (default: `8080`) configured during setup.

## [IMPORTANT] Security

**This software is considered as proof of concept.**

*As this project originated from a bachelor's thesis with CAN privacy as its main topic, neither encryption nor secure API-Endpoints have been developed and planned as future work. At this point, it is highly discouraged to use the proxy in production or expose it to the internet.*
