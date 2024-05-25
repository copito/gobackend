PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := protoc-gen-grpc-gateway

# point to proto files outside service
PROTO_DIR := idl/grpc
OUT_DIR := idl-gen

# find all proto files
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

## IDL
# proto: $(PROTO_FILES)
# 	for file in $(PROTO_FILES); do \
# 		protoc \
# 			--proto_path=$(PROTO_DIR) \
# 			--go_out=$(OUT_DIR) \
# 			--go_opt=paths=source_relative \
# 			--go-grpc_out=$(OUT_DIR) \
# 			--go-grpc_opt=paths=source_relative \
# 			--grpc-gateway_out=$(OUT_DIR) \
# 			--grpc-gateway_opt=paths=source_relative \
# 			$$file; \
# 	done

proto:
	buf generate

clean_proto:
	rm -rf $(OUT_DIR)/*

install_proto:
	go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

install_buf:
	# proto
	sudo apt install -y protobuf-compiler 
	# sudo go install github.com/grpc-ecosystem/grpc-gateway/v2@latest
	# sudo go install google.golang.org/protobuf

	# buf
	sudo GO111MODULE=on GOBIN=/usr/local/bin go install github.com/bufbuild/buf/cmd/buf@v1.32.1 
