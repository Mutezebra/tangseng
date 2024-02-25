package myprometheus

import (
	"fmt"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// Register grpcServer with prometheus
// and addr with etcd
func RegisterServer(server *grpc.Server, addr string, job string) {
	fmt.Println("will register")
	grpcPrometheus.Register(server)
	EtcdRegister(addr, job)
	go RpcHandler(addr)
}
