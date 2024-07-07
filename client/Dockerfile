FROM golang:1.22 AS builder

WORKDIR /app

COPY client/go.mod client/go.sum ./

COPY lib ../lib

RUN go mod download

COPY client .

RUN go build -o /app/client .


FROM debian:12-slim

COPY --from=builder /app/client .

EXPOSE 3001

ENTRYPOINT ["./client"]
