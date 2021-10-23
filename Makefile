all: client server

client: # build client
	cd client-go; go build -o ../bin/gotalk-client


server: # build server
	cd server-go; go build -o ../bin/gotalk-server

clean:
	rm bin/*
