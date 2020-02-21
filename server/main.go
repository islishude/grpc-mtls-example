package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/islishude/grpc-mtls-example/greet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func main() {
	server := grpc.NewServer(
		grpc.Creds(LoadKeyPair()),
		grpc.UnaryInterceptor(middlefunc),
	)

	greet.RegisterGreetingServer(server, new(GreetServer))

	go func() {
		l, err := net.Listen("tcp", ":10200")
		if err != nil {
			panic(err)
		}
		log.Println("listen and server...")
		if err := server.Serve(l); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	server.GracefulStop()
	log.Println("bye")
}

func middlefunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// get client tls info
	if p, ok := peer.FromContext(ctx); ok {
		if mtls, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			for _, item := range mtls.State.PeerCertificates {
				log.Println(item.Subject)
			}
		}
	}
	return handler(ctx, req)
}

type GreetServer struct{}

func (g *GreetServer) SayHello(ctx context.Context, req *greet.SayHelloRequest) (*greet.SayHelloResponse, error) {
	respdata := "Hello," + req.GetName()
	return &greet.SayHelloResponse{Greet: respdata}, nil
}

func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		panic("Load server certification failed: " + err.Error())
	}

	data, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		panic("can't read ca file")
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(data) {
		panic("can't add ca cert")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    capool,
	}
	return credentials.NewTLS(tlsConfig)
}
