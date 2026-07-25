package server

import (
	pb "bastiond/generated"
	"context"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
}

func (m *MetricsServer) GetMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.Metrics, error) {
	return nil, nil
}

func (m *MetricsServer) Connect(ctx context.Context, req *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	return nil, nil
}
