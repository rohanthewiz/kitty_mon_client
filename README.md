# Kitty Monitor Client
This is the client portion of the Kitty Monitor system

#### Project Status: Beta

## Synopsis
Kitty Monitor is a system for monitoring the temperature of a remote location.
* For simplicity and cost-effectiveness, the internal board temperature of an Odroid C1+/C2 (buy at hardkernel.com or ameridroid.com) is the sensor used
* The client is designed to run on a Ubuntu based operating system
* The Odroid locally polls it's internal system temperature and saves it to a SQLite database (typically every 2mins)
* The Odroid sends any unsent readings to the server (typically every 4mins)
* The server collects and displays these readings in a webserver
* The temperature provided is not an exact reading of ambient temperature, but an indication in case of say failure of an airconditioning unit

## Getting Setup

### Download
Get Go for your operating system: http://golang.org/dl/ and install.

## Getting and Building KittyMonitor
```
# cd to the directory of your choice then clone the repo:
git clone https://github.com/rohanthewiz/kitty_mon_client.git
cd kitty_mon_client
go build
```

### COMPILE with Docker
- Build the compiler docker image: `docker build -f Dockerfile.compile -t go:compiler .`
- Use a container to cross-compile, for example, to arm64 linux:
 `docker run --rm -v "$HOME/<path/to/kitty_mon/project>:root" -w /root -e GOOS=linux -e GOARCH=arm64 -e CGO_ENABLED=1 go:compiler go build -v -ldflags '-w -extldflags "-static"' -o app .`

### Standard compile
```
# You need a C compiler that can build sqlite.
CGO_ENABLED=1 go build # this will produce the executable 'kitty_mon' in the current directory
```

Copy the `app` executable where needed. Example: `scp app user@myserver.net:bin/kitty_mon`

## Using KittyMonitor
Note that the single compiled binary can operate as server or client.
KittyMonitor launches a goroutine that polls system board temperature of the Odroid.

The main thread continues in an infinite loop to:
 1. poll temperature every 2 mins
 2. connect to the server every 4 mins and send any unsent readings (up to a limit).
 
Usage info: `kitty_mon -h` 

### Get the server's secret token (you'll need to copy the token unto the client)

```
# At the server run
./kitty_mon -get_server_secret # => 9A7blahblah123...
```

### Starting the Client in Development mode (Readings every 8s. The server_secret is from above)

```
$ ./kitty_mon -synch_client ServerNameOrIP -server_secret 9A7blahblah123...
```

### Starting the Client in Production mode (Readings every 2 minutes)

```
$ ./kitty_mon -synch_client ServerNameOrIP -server_secret 9A7blahblah123... -env prod
```
(The secret code is only needed the first time the client talks to the server)

### Monitor Raw data
In a browser go to http://ServerNameOrIP:9080
See my example at http://gonotes.net:9080

### Other Options
    
    -h -- List available options with defaults
    -db "" -- Sqlite DB path. KittyMon will try to create the database 'kitty_mon.sqlite' in your home directory by default
    -admin="" -- Privileged actions like 'delete_tables' (drops the readings table)

### Local Testing
For testing you can run a client and server on the same machine, just use different databases
Example:

```
# Start the server with the default database
./kitty_mon
# Start the client with test.sqlite
./kitty_mon_client -db test.sqlite -synch_client localhost -node_name TestClient
# View at localhost:9080
```
If no readings are listed then your workstation may not support internal temperature reporting.
That's okay, the Odroid C1/C2 does!

### TODO
- A lot since we are in Beta
- Token based auth webserver mode
- Allow instructions for the client (example: reboot, delete n readings)

### TIPS
- gcc is required to compile SQLite. On Windows you can get 64bit MinGW here http://mingw-w64.sourceforge.net/download.php. Install it and make sure to add the bin directory to your path
  (Could not get this to work with Windows 8.1, Windows 7 did work though)
- For less typing, you might want to do 'go build -o km' (km.exe on Windows) to produce the executable 'km'
- This is a sweet way to learn a modern, highly performant language - The Go Programming Language using a database with an ORM (object relational manager) while building a useful tool!
I recommend using git to checkout an early version of kitty_mon_client so you can start out simple
- Firefox has a great addon called SQLite Manager which you can use to peek into the database file
- Feel free to create a pull request if you'd like to pitch in.

### Credits
- Go -- http://golang.org/  Thanks Google!!
- GORM -- https://github.com/jinzhu/gorm  - Who needs sluggish ActiveRecord, or other interpreted code interfacing to your database.
- SQLite -- http://www.sqlite.org/ - A great place to start. Actually GORM includes all the things needed for SQLite so SQLite gets compiled into KittyMon!
