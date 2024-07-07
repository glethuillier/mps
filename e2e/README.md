# mps | End-to-End Test

This end-to-end test generates random files and sends them to the mps client. It then asks the client to verify each file individually.

## Run

Assuming that Go is installed on your system (if not: [how to install Go](https://go.dev/doc/install)), run:

```
$ go run main.go
```

By default, the test generates and sends 100 files. To adjust the number of files, specify it in the command:

```
$ go run main.go {{n}}
```

Example (generate and send 1,000 files):

```
$ go run main.go 1000
```

## Usage

This tool only performs positive tests.

Negative tests can be done manually: run a test then, when the execution of the test pauses (after having sent the files), corrupt one of the files (in `./server/downloads/{{root_hash}}/`), and resume the test. The client should automatically detect the discrepancy.
