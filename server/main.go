package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	"example.com/mtls/hello"
)

var _ hello.HelloServer = (*helloController)(nil)

type helloController struct{}

func (hc *helloController) SayHello(ctx context.Context, req *hello.Request) (*hello.Response, error) {
	name := req.GetName()
	if name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "")
	}

	return &hello.Response{Greet: fmt.Sprintf("Hello,%s", name)}, nil
}

func main() {
	certificate, err := tls.LoadX509KeyPair("server/cert.pem", "server/key.pem")

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("ca/cacert.pem")
	if err != nil {
		log.Printf("failed to read client ca cert: %s", err)
		return
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Println("failed to append client certs")
		return
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	server := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	hello.RegisterHelloServer(server, new(helloController))

	go func() {
		lis, err := net.Listen("tcp", ":10200")
		if err != nil {
			log.Fatalln("failed to listen", err)
		}
		log.Println("Listening :10200 port")
		if err := server.Serve(lis); err != nil {
			log.Fatalln("grpc.Serve error", err)
		}
	}()

	killSignals := make(chan os.Signal, 1)
	signal.Notify(killSignals, syscall.SIGINT, syscall.SIGTERM)

	<-killSignals

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	graceful := make(chan struct{})

	go func() {
		server.GracefulStop()
		close(graceful)
	}()

	select {
	case <-ctx.Done():
		log.Println("graceful stop timeout")
	case <-graceful:
		log.Println("bye")
	}
}
