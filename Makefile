PROTO_PATH = proto
GO_OUT_DIR = pkg/messages
PROTO_FILE = messages.proto

all: proto build_docker

proto:
	cd lib && \
	protoc --proto_path=$(PROTO_PATH) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative $(PROTO_FILE)

build_docker:
	docker build -f server/Dockerfile -t server .
	docker build -f client/Dockerfile -t client .

.PHONY: all proto build_docker
