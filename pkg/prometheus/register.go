package prometheus

import (
	"fmt"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

var (
	UnaryServerInterceptor  = grpcPrometheus.UnaryServerInterceptor
	StreamServerInterceptor = grpcPrometheus.StreamServerInterceptor

	EnableHandlingTimeHistogram = grpcPrometheus.EnableHandlingTimeHistogram
)

func RegisterServer(server *grpc.Server, port string, domain string, job string) {
	fmt.Println("will register")
	grpcPrometheus.Register(server)
	EtcdRegister(domain, job)
	go RpcHandler(domain)
}
