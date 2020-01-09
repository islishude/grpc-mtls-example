package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"test/greet"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewTLS() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("client/cert.pem", "client/cert-key.pem")
	if err != nil {
		panic("Load client certification failed: " + err.Error())
	}

	data, err := ioutil.ReadFile("rootca/rootca.pem")
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(data) {
		panic("can't add CA cert")
	}

	tlsConfig := &tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(tlsConfig)
}

func main() {
	conn, err := grpc.Dial("localhost:10200", grpc.WithTransportCredentials(NewTLS()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := greet.NewGreetingClient(conn)
	resp, err := client.SayHello(context.Background(), &greet.SayHelloRequest{Name: "world"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.GetGreet())
}
