GO = go
DEP = dep
PROTOC = protoc

SERVICE = uzucoin
PROTO_DIR = protobuf
SERVER_DIR = src/server
PROTO_SERVER_DIR = $(SERVER_DIR)/pb

.DEFAULT_GOAL = run

proto:
	@$(PROTOC) -I $(PROTO_DIR)/ $(PROTO_DIR)/$(SERVICE).proto --go_out=plugins=grpc:$(PROTO_SERVER_DIR)

dep:
	@cd $(SERVER_DIR) && $(DEP) ensure

build: proto
	@cd $(SERVER_DIR) && $(GO) build -o $(SERVICE)

run: build
	./$(SERVER_DIR)/$(SERVICE)

clean:
	@rm $(SERVER_DIR)/$(SERVICE)
