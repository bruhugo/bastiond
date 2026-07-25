package main

import (
	pb "bastiond/generated"
	server "bastiond/internal"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMetricsServiceServer(grpcServer, &server.MetricsServer{})
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
