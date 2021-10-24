# Example Dockerfile to create a server container
#
FROM golang:1.16-alpine
# Move to working directory /build
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the code into the container
COPY language.json ./
RUN mkdir ./utils ./secret ./models ./server-tcp
COPY utils/* ./utils
COPY secret/* ./secret
COPY models/* ./models
COPY server-tcp/* ./server-tcp

# Build the application
WORKDIR /app/server-tcp
RUN go build -o /gotalk-server

# Export default port
EXPOSE 8089

# Command to run when starting the container
CMD ["/gotalk-server"]
# Command example below to run when starting the container with german locale
#CMD ["/gotalk-server", "-l", "de"]
