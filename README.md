# *mps* | Merkle Proof Service

This Merkle Proof Service (*mps*), composed of two components (a client and a server), uses Merkle tree-based proofs to ensure that files sent by the client to the server and then downloaded back by the client are not corrupted.

(This project is a working proof of concept implemented for educational purposes. Contributions are welcomed!)

## Run

The client and the server can be run natively (check their respective `README`s).

To run them using Docker Compose:

```
$ docker-compose up
```

They can additionally be run individually:

```
$ docker-compose up server
$ docker-compose up client
```

## Test (End-to-End test)

While the client and the server are running, run from the `e2e` subdirectory:

```
$ go run main.go
```

## Components

### Client

The *mps* client (client subdirectory), written in Go, implements a REST API server that provides endpoints through which files can be sent by batch to the *mps* server (`/upload`) and individually fetched from the *mps* server (`/download`). It pushes these requests to the server via a WebSocket connection.

Before sending a set of files to the server, the client constructs the corresponding Merkle tree root hash. Then, the client stores the root hash in the database alongside the receipt ID. The files are never stored on the client side. If the files have already been sent to the server, an error message is returned with the relevant receipt ID.

To download a file from the server, a valid receipt ID and a filename are required (a receipt ID is used for two reasons: asking for a file just based on its filename would lead to collisions, and to make the caller, who is not necessarily a cryptograph enthusiast, deal with a familiar UUID instead of thinking in terms of a Merkle proof). The client receives the file with the proof, reconstructs the root hash based on it, and then compares it with the one stored in its database. If they match, the client returns the file to the caller. An error message is returned if the file does not exist on the server or if the verification fails (in that case, with a `427` status code—invalid digital signature—used as an umbrella term as it is not a signature per se).

### Server

When the client sends files to the server (server subdirectory), the latter, also written in Go, computes their Merkle tree root hash. If this hash does not match the one provided by the client, the server does not store the files and returns an error. Otherwise, the server saves the files and saves the proof in its database.

The database contains the following tables:
* `RECEIPTS`, which stores the receipt IDs and the corresponding Merkle tree root hashes.
* `FILES`, which stores the filenames and the hashes of the files they refer to. 
* `TREES`, which stores a representation of the Merkle trees (node hash, sibling, sibling type—left, right, or none—, and parent).

Generating a proof is then a question of retrieving the hash for a given file and, up to the root, identifying the sibling of the current child and its position in the subtree (left, right). The proof is then Protobuf serialized and sent to the client with the file.

### Library

The library defines and implements the Protobuf messages. It also defines the structure of the parts that composes a proof.
