run:
	go run server.go -f etc/server.yaml

api:
	goctl api go -api server.api --style=go_zero -dir .

swagger:
	goctl api plugin -plugin goctl-swagger="swagger -filename swagger.json" -api server.api 