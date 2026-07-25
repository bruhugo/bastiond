.PHONY: run build lint proto build 

proto: 
	mkdir -p generated
	protoc -I=proto --go_out=generated --go_opt=paths=source_relative \
    --go-grpc_out=generated --go-grpc_opt=paths=source_relative \
    metrics.proto
	
run:
	go run .

build:
	go build .

lint:
	go lint ./...
