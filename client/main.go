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
	certificate, err := tls.LoadX509KeyPair("client/cert.pem", "client/key.pem")
	certPool := x509.NewCertPool()

	caCert, err := ioutil.ReadFile("ca/cacert.pem")
	if err != nil {
		log.Printf("Read CA certs failed %s\n", err)
		return
	}

	if !certPool.AppendCertsFromPEM(caCert) {
		log.Println("Add CA cert to certs pool failed")
		return
	}

	// Make sure that `127.0.0.1 => server.dev`
	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "server.dev",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial("server.dev:10200", grpc.WithTransportCredentials(transportCreds))
	if err != nil {
		log.Printf("failed to dial server: %s\n", err)
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
