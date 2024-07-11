# mps | Client

This *mps* (Merkle Proof Service) client uploads files to a *mps* server and, using a Merkle tree, ensures that they have not been corrupted in any way.

## Run

Assuming that Go is installed on your system (if not: [how to install Go](https://go.dev/doc/install)), run:

```
$ go run main.go
```

The log level can be set with the environment variable `LOG_LEVEL`. Example: 

```
$ LOG_LEVEL=DEBUG go run main.go
```

The server host and port (default:`localhost` and `3000`) can be changed with the environment variables `SERVER_HOST` and `SERVER_PORT`. Example: 

```
$ SERVER_HOST=10.0.0.1 SERVER_PORT=1234 go run main.go
```

## Usage

### Upload files

```
curl --request POST \
  --url 'http://localhost:3001/upload' \
  --header 'content-type: multipart/form-data' \
  --form file={{file1}} \
  --form file={{file2}} \
  . . .
```

Example:

```
curl --request POST \
  --url 'http://localhost:3001/upload?=&=' \
  --header 'content-type: multipart/form-data' \
  --form file=@/Users/you/Documents/file1.txt \
  --form file=@/Users/you/Documents/file2.jpg \
  --form file=@/Users/you/Documents/file3.docx
```

If the request succeeds, the client returns a receipt ID **hat you should keep to download your files subsequently**.

### Download files

```
curl --request POST \
  -v \
  --url 'http://localhost:3001/download' \
  --header 'Content-Type: application/json' \
  --data '{
	"receipt_id": "{{receipt ID}}",
	"filename": "{{filename}}"
}'
  --output {{filename}}
```

Example:

```
curl --request POST \
  -v \
  --url 'http://localhost:8080/download' \
  --header 'Content-Type: application/json' \
  --data '{
	"receipt_id": "0aa9a4bc-7554-4d6b-bebb-77b23dfc321b",
	"filename": "test1.txt"
}'
  --output text1.txt
```

If the verification succeeds, the file is downloaded. Otherwise, an error message is returned.

Note: the proof is returned in the headers (`Proof-*`).
