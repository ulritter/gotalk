# gotalk

## Simple Multi-user ad-hoc communication program.
**The communication is secured using tls over tcp. The program can be started in server mode or in client mode (see below). In client mode the program starts a graphical user interface to accomodate both conversations and status messages. The client GUI is built using `fyne` (https://fyne.io/), a portable graphical toolkit.**

&NewLine; 
**Build the software:**
- install golang
- install fyne (see also https://developer.fyne.io/index.html)
- Linux: on a standard Ubuntu 20.04 distro I had to install:
  `sudo apt-get install libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama1 libxinerama-dev libxi-dev libxxf86vm-dev`
- windows:  for 64 bit gcc (if not already installed) get the MinGW-w64 installer on the website below and chose x86_64 architecture during install:
  `http://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win32/Personal%20Builds/mingw-builds/installer/mingw-w64-install.exe/download`
- install localization package: `go get github.com/moemoe89/go-localization`
- clone / download this repo
- rename `secret.go.example` to `secret.go`
- run `openssl ecparam -genkey -name prime256v1 -out server.key`
- replace `serverKey` constant dummy content with content of `server.key` file
- run `openssl req -new -x509 -key server.key -out server.pem -days 3650`
- replace `rootCert` constant dummy content with content of `server.pem` file
- run `go build`


&NewLine;  
&NewLine;  

**Run the software in server mode:**

	gotalk server [<port>] 

**Examples:**

    ./gotalk server
    ./gotalk server 8089

Server termination by SIGHUP (for the time being)

**Run the software in client mode:**

	gotalk client [<nickname> [<address>] [<port>]]

**Examples:**

    ./gotalk client
    ./gotalk client MyNick
    ./gotalk client MyNick 127.0.0.1
    ./gotalk client MyNick 127.0.0.1 8089

![Client example](https://github.com/ulritter/gotalk/blob/main/example.png)

&NewLine;   

**Client commands:**
- ,`/exit`,`/quit`,`/q` - close connection and exit
- `/list` - displays active users in room
- `/nick <nickname>` - change nickname

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


In all cases \<address\> defaults to `localhost`, \<port\> defaults to `8089`, and \<nickname\> defaults to `J_Doe`

