build:
	mkdir -p dist
	go build -o ./dist/client client/main.go
	go build -o ./dist/server server/main.go

protogen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/*.proto
