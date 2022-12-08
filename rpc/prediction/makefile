run:
	go run prediction.go -f etc/prediction.yaml

rpc:
	goctl rpc protoc prediction.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style=go_zero