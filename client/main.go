package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"example.com/mtls/hello"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("client/client.local-cert.pem", "client/client.local-key.pem")
	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile("demoCA/cacert.pem")
	if err != nil {
		log.Fatalf("Read CA certs failed %s\n", err)
	}

	certPool.AppendCertsFromPEM(caCert)
	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "dev.local",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	// add `127.0.0.1 dev.local` to `/etc/hosts`
	conn, err := grpc.Dial("dev.local:10200", grpc.WithTransportCredentials(transportCreds))
	if err != nil {
		log.Fatalf("failed to dial server: %s\n", err)
		return
	}
	defer conn.Close()

	client := hello.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := client.SayHello(ctx, &hello.Request{Name: "world"})
	if err != nil {
		log.Println("Call SayHello failed", err.Error())
		return
	}
	log.Println("Message from server", resp.GetGreet())
}
