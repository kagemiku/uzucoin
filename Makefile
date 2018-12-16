PROTOC = protoc

SERVICE = uzucoin
PROTO_DIR = protobuf
SERVER_DIR = src/server
PROTO_SERVER_DIR = $(SERVER_DIR)/pb

proto:
	@$(PROTOC) -I $(PROTO_DIR)/ $(PROTO_DIR)/$(SERVICE).proto --go_out=plugins=grpc:$(PROTO_SERVER_DIR)
