protogen:
	protoc -I hello --go_out=plugins=grpc:hello hello/*.proto
	goimports -w hello/*.go

build:
	go build -o server.out server/main.go
	go build -o client.out client/main.go
