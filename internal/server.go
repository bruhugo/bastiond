package server

import (
	pb "bastiond/generated"
	"context"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
}

func (m *MetricsServer) GetMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.Metrics, error) {
	metrics, err := getMetrics(ctx)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (m *MetricsServer) Connect(ctx context.Context, req *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	// TODO: implement auth
	return &pb.ConnectionResponse{
		Status: pb.ConnectionStatus_ACCEPTED,
	}, nil
}
