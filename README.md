gRPC mLTS example
---

## generate eckey file

```
openssl ecparam -genkey -name secp384r1 -noout -out key.pem
```

## self sign ca

```
openssl req -x509 -key key.pem -new -out myca.local
```

## generate certification request file

```
openssl req -key key.pem  -new -config openssl.conf -out csr.pem
```

## sign a certification signing request file

```
openssl ca -config openssl.conf -in csr.pem -out cert.pem
```
