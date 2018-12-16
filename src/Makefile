SERVICE = uzucoin
PROTOC = protoc
PROTO_DIR = protobuf

proto:
	@$(PROTOC) -I $(PROTO_DIR)/ $(PROTO_DIR)/$(SERVICE).proto --go_out=plugins=grpc:$(PROTO_DIR)
