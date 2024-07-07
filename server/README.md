# mps | Server

This *mps* (Merkle Proof Service) server stores files sent by the *mps* client and, when the client requires it, sends back files alongside a Merkle tree-based proof.

## Run

Assuming that Go is installed on your system (if not: [how to install Go](https://go.dev/doc/install)), run:

```
$ go run main.go
```

The log level can be set with the environment variable `LOG_LEVEL`. Example: 

```
$ LOG_LEVEL=DEBUG go run main.go
```

The port (default: `3000`) can be changed with the environment variable `PORT`. Example: 

```
$ PORT=1234 go run main.go
```

## Usage

You need to use the *mps* client to interact with the server.

