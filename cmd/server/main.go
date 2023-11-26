package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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
	tlsConfig, err := LoadTlSConfig("server.pem", "server-key.pem", "root.pem")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(grpc.Creds(tlsConfig), grpc.UnaryInterceptor(MiddlewareHandler))

	greet.RegisterGreetingServer(server, new(GreetServer))

	basectx, casncel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer casncel()

	listener, err := net.Listen("tcp", ":10200")
	if err != nil {
		panic(err)
	}

	go func() {
		<-basectx.Done()
		server.GracefulStop()
		log.Println("bye")
	}()

	log.Println("listen and serving...")
	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}

func MiddlewareHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// you can write your own code here to check client tls certificate
	if p, ok := peer.FromContext(ctx); ok {
		if mtls, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			for _, item := range mtls.State.PeerCertificates {
				log.Println("client certificate subject:", item.Subject)
			}
		}
	}
	return handler(ctx, req)
}

type GreetServer struct {
	greet.UnimplementedGreetingServer
}

func (g *GreetServer) SayHello(ctx context.Context, req *greet.SayHelloRequest) (*greet.SayHelloResponse, error) {
	respdata := "Hello," + req.GetName()
	return &greet.SayHelloResponse{Greet: respdata}, nil
}

func LoadTlSConfig(certFile, keyFile, caFile string) (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certification: %w", err)
	}

	data, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("faild to read CA certificate: %w", err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("unable to append the CA certificate to CA pool")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    capool,
	}
	return credentials.NewTLS(tlsConfig), nil
}
