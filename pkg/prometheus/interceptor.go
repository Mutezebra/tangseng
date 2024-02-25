package prometheus

import grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

var (
	UnaryServerInterceptor  = grpcPrometheus.UnaryServerInterceptor
	StreamServerInterceptor = grpcPrometheus.StreamServerInterceptor

	EnableHandlingTimeHistogram = grpcPrometheus.EnableHandlingTimeHistogram
)
