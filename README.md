# gotalk

## Simple multi-user ad-hoc communication program.
**The communication is secured using tls over tcp. The program can be started in server mode or in client mode (see below). In client mode the program starts a graphical user interface to accomodate both conversations and status messages. The client GUI is built using `fyne` (https://fyne.io/), a portable graphical toolkit.**

&NewLine; 
**Build the software:**
- install golang
- install fyne (see also https://developer.fyne.io/index.html)
- install fyne system dependencies, e.g.:
  - Linux: on a standard Ubuntu 20.04 distro I had to install:
    `sudo apt-get install libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama1 libxinerama-dev libxi-dev libxxf86vm-dev`
  - windows:  for 64 bit gcc (if not already installed) get the MinGW-w64 installer on the website below and chose x86_64 architecture during install:
- install command line parser (`go get github.com/alecthomas/kong`)
- install localization package: (`go get github.com/moemoe89/go-localization`)
- clone / download this repo
- rename `secret.go.example` to `secret.go`
- run `openssl ecparam -genkey -name prime256v1 -out server.key`
- replace `serverKey` constant dummy content with content of `server.key` file
- run `openssl req -new -x509 -key server.key -out server.pem -days 3650`
- replace `rootCert` constant dummy content with content of `server.pem` file
- install `make`if not already present on your system
- run `make all` to build both client and server binaries (target: `./bin` directory)
- run `make client` to build the client binary (target: `./bin` directory)
- run `make server` to build the server binary (target: `./bin` directory)


&NewLine;  
&NewLine;  

**Run the software in server mode:**

    gotalk-server [options] 
    Usage: gotalk-server

    Flags:
        -h, --help           Show context-sensitive help.
        -p, --port="8089"    Port number.
        -l, --locale="en"    Language setting to be used.

**Examples:**

    
    ./gotalk-server 
    ./gotalk-server -p 8089 --locale=de
    ./gotalk-server -l de
    ./gotalk-server --port=8089

Server termination by SIGHUP (for the time being)

**Run the software in client mode:**

	gotalk client [options]

    Usage: gotalk-client

    Flags:
        -h, --help                   Show context-sensitive help.
        -a, --address="localhost"    IP address or domain name.
        -p, --port="8089"            Port number.
        -n, --nick="J_Doe"           Nickname to be used.
        -l, --locale="en"            Language setting to be used.

**Examples:**

    ./gotalk-client
    ./gotalk-client --nick MyNick 
    ./gotalk-client -n MyNick --address=127.0.0.1
    ./gotalk-client --nick=MyNick -a 127.0.0.1 --port 8089 --locale de

![Client example](https://github.com/ulritter/gotalk/blob/main/example.png)

&NewLine;   

**Client commands:**
- `/exit`,`/quit`,`/q` - close connection and exit
- `/list` - displays active users in room
- `/nick <nickname>` - change nickname
- `/help`,`/?` - display help text

&NewLine;   
&NewLine;   

**Color controls:**

  General:
  - a color control followed by space will change the color for the rest of the line
  - a color control attached to a word will change the color for the word
 
 Usage Example:
`$red` this is my `$y`text
 
Color Controls (long form and short form for each control)
`$red $r $cyan $c $yellow $y $green $g`
`$purple $p $white $w $black $b` 

&NewLine;
&NewLine;   


In all cases \<address\> defaults to `localhost`, \<port\> defaults to `8089`, and \<nickname\> defaults to `J_Doe`,
and \<locale\> defaults to the actual system setting. If no translation is available it falls back to english.

TODOS:
- switch to https based communication
- create web client (React)
