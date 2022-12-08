run:
	go run server.go -f etc/server.yaml

rpc:
	goctl rpc protoc server.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style=go_zero