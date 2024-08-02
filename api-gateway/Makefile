CURRENT_DIR=${shell pwd}

gen-proto:
	@protoc \
	--go_out=. \
	--go-grpc_out=. \
	protos/auth.proto \
	protos/deviceControl.proto


SWAGGER := ~/go/bin/swag
SWAGGER_DOCS := docs
SWAGGER_INIT := $(SWAGGER) init -g ./main.go -o $(SWAGGER_DOCS)

swag-gen:
	$(SWAGGER_INIT)
