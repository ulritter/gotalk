all: client server

client: # build client
	cd client-tcp; go build -o ../bin/gotalk-client


server: # build server
	cd server-tcp; go build -o ../bin/gotalk-server

clean:
	rm bin/*
