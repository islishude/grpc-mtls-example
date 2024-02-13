# gRPC mTLS example

## Install protoc

Download it from the latest release page

https://github.com/protocolbuffers/protobuf/releases/latest

If you're using MacOS, you can install protoc with brew:

```
brew install protobuf
```

by the way, you can also install clang-format for formating the proto files.

```
brew install clang-format
```

and run following command to format the protobuf files

```sh
clang-format -i greet/*.proto
```

## Install protoc generator for golang

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Install cfssl

It's used for generating X509 certificates

```
go install github.com/cloudflare/cfssl/cmd/...@latest
```

you can also use OpenSSL, but it will become more complex :)

## Build proto files

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/*.proto
```

## Create root CA

```
cfssl selfsign -config cfssl.json --profile rootca "Dev Testing CA" csr.json | cfssljson -bare root
```

you will get 3 files:

- root.csr ROOT CA CSR(you may don't need it)
- root-key.pem ROOT CA key
- root.pem ROOT CA certificate

the `cfss.json` here is a configuration file for cfssl, if you would like use your own `cfss.json` and `csr.json`, please check out the [cfssl](https://github.com/cloudflare/cfssl) documentation for the deails.

## Create certificates for server and client

```
cfssl genkey csr.json | cfssljson -bare server
cfssl genkey csr.json | cfssljson -bare client
```

you will get 4 files:

- server.csr Server CSR
- server-key.pem Server key
- client.csr Client CSR
- client-key.pem Client key

the CSR files will be used for signing new certificate

## Sign the certificates

```
cfssl sign -ca root.pem -ca-key root-key.pem -config cfssl.json -profile server server.csr | cfssljson -bare server
cfssl sign -ca root.pem -ca-key root-key.pem -config cfssl.json -profile client client.csr | cfssljson -bare client
```

you will get your server and client certificates

- server.pem
- client.pem

## Build client and server

```
mkdir -p dist
go build -o ./dist ./cmd/...
```

## Run and test

```console
$ ./dist/server
2022/07/21 22:18:20 listening on 6443
2022/07/21 22:18:26 request certificate subject: CN=client
```

```console
$ ./dist/client
Hello,world
```
