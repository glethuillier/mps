FROM golang:1.22 AS builder

WORKDIR /app

COPY server/go.mod server/go.sum ./

COPY lib ../lib

RUN go mod download

COPY server .

RUN go build -o /app/server .


FROM debian:12-slim

COPY --from=builder /app/server .

EXPOSE 3000

ENTRYPOINT ["./server"]
