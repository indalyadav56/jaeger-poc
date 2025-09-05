PROTO_DIR=api/proto
PB_DIR=pb
THIRD_PARTY_DIR=third_party
SWAGGER_DIR=api/swagger

AUTH_PROTO=$(PROTO_DIR)/auth.proto

# Dynamically resolve grpc-gateway v2 module path from Go modules
PROTOC_GEN_OPENAPIV2=$(shell go list -m -f '{{.Dir}}' github.com/grpc-ecosystem/grpc-gateway/v2)

proto: $(SWAGGER_DIR)
	protoc -I. \
		-I$(THIRD_PARTY_DIR) \
		-I$(PROTOC_GEN_OPENAPIV2) \
		--go_out=$(PB_DIR) --go-grpc_out=$(PB_DIR) \
		--grpc-gateway_out=$(PB_DIR) \
		--openapiv2_out=$(SWAGGER_DIR) \
		--openapiv2_opt=allow_merge=true,merge_file_name=auth \
		$(AUTH_PROTO)

$(SWAGGER_DIR):
	mkdir -p $(SWAGGER_DIR)

.PHONY: proto

run:
	go run cmd/server.go