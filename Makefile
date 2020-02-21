build:
	mkdir -p dist
	go build -o ./dist/client client/main.go
	go build -o ./dist/server server/main.go
