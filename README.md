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

## How to generate x509 certificates

1. Download [cfssl](https://github.com/cloudflare/cfssl)
2. Generate your self-signed root CA

```
cfssl selfsign -config cfssl.json --profile rootca "My Root CA" csr.json | cfssljson -bare root
```

you will get 3 files:

- root.csr ROOT CA CSR(you may don't need it)
- root-key.pem ROOT CA key
- root.pem ROOT CA certificate

3. Create your server and client certificate

```
cfssl genkey csr.json | cfssljson -bare server
cfssl genkey csr.json | cfssljson -bare client
```

you will get 4 files:

- server.csr Server side CSR
- server-key.pem Server key
- client.csr Client side CSR
- server-key.pem Client key

4. Sign new certificates by your self-signed root CA

```
cfssl sign -ca root.pem -ca-key root-key.pem -config cfssl.json -profile server server.csr | cfssljson -bare server
cfssl sign -ca root.pem -ca-key root-key.pem -config cfssl.json -profile client client.csr | cfssljson -bare client
```

you will get your server and client certificates

- server.pem
- client.pem

For more detail about `cfss.json` and `csr.json`, check out cfsll documentation.
