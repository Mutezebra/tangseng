package prometheus

import (
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// Register grpcServer with prometheus
// and addr with etcd
func Register(server *grpc.Server, addr string, job string) {
	grpcPrometheus.Register(server)
	EtcdRegister(addr, job)
	go RpcHandler(addr)
}
