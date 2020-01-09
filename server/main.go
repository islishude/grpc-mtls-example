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
	"test/greet"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func main() {
	server := grpc.NewServer(
		grpc.Creds(NewTLS()),
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

func NewTLS() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("server/cert.pem", "server/cert-key.pem")
	if err != nil {
		panic("Load server certification failed: " + err.Error())
	}

	data, err := ioutil.ReadFile("rootca/rootca.pem")
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(data) {
		panic("can't add ca cert")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}
	return credentials.NewTLS(tlsConfig)
}
