# gotalk

## Simple Multi-user ad-hoc communication program.
## The communication is secured using tls over tcp.
## The program can be started in server mode or in client mode (see below)

Build: install golang, clone / download this repo and

    go build

to build the software

&NewLine;  
&NewLine;  

**Server mode invocation:**

	gotalk server [<port>] 

Server termination by SIGHUP (for the time being)

**Client mode invocation:**

	gotalk client [<nickname> [<address>] [<port>]]

Client commands:
- /exit - terminate connection and exit
- /list - displays active users in room
- /nick <nickname> - change nickname

&NewLine;   
&NewLine;   

In all cases \<address\> defaults to `localhost` and port defaults to `8080`

