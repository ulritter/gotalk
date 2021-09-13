# gotalk

## Simple Multi-user ad-hoc communication program

**Server invocation:**

	gotalk server [<port>] 

Server termination by SIGHUP (for the time being)

**Client invocation:**

	gotalk client [<nickname> [<address>] [<port>]]

Client termination by entering STOP


In all cases \<address\> defaults to `localhost` and port defaults to `8080`

