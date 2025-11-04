install_utils:
	go install github.com/yoheimuta/protolint/cmd/protolint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

lint:
	@protolint .

lint_fix:
	@protolint --fix .

generate:
	protoc -I=./ --go_out=./ --go-grpc_out=./ ./file_storage.proto