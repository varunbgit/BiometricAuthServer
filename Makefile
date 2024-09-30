BUF_VERSION := v1.4.0
PROTOC_GEN_GO_VERSION := latest

.PHONY: proto-generate-deps
proto-generate-deps:
	@echo "\n + Fetching protobuf dependencies \n"
	@go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	@go install	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install golang.org/x/lint/golint@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	@go install github.com/golang/protobuf/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)


.PHONY: proto-generate ## Compile protobuf
proto-generate:
	@echo "\n + Generating pb language bindings\n"
	@buf mod update
	@buf ls-files $(PROTO_ROOT)
	@buf generate

.PHONY: clean ## Remove previous builds, protobuf files, and proto compiled code
clean:
	@echo " + Removing cloned and generated files\n"
	##- todo: use go clean here
	@rm -rf $(API_OUT) $(MIGRATION_OUT) $(RPC_ROOT)


.PHONY: all  ## build project
all: clean proto-generate-deps proto-generate