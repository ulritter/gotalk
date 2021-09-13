# gotalk

## Simple Multi-user ad-hoc communication program

Build installing golang, by cloning this repo and by

    go build

to build the software

&NewLine;  
&NewLine;  

**Server invocation:**

	gotalk server [<port>] 

Server termination by SIGHUP (for the time being)

**Client invocation:**

	gotalk client [<nickname> [<address>] [<port>]]

Client termination by entering STOP

&NewLine;   
&NewLine;   

In all cases \<address\> defaults to `localhost` and port defaults to `8080`

