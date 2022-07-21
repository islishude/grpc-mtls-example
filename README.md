# GRPC mTLS example

## Install protoc(for MacOS)

```
brew install protobuf clang-format
```

clang-format is used for format proto files


## Install golang code generator

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Build proto files

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/*.proto
```

## Build client and server

```
rm -rf dist && mkdir -p dist
go build -o ./dist ./cmd/...
```

## Run and test

```console
$ ./dist/server 
2022/07/21 22:18:20 listen and serveing...
2022/07/21 22:18:26 request certificate subject: CN=client
```

```console
$ ./dist/client
Hello,world
```
