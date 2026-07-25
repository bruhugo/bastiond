package server

import (
	pb "bastiond/generated"
	"context"
	"fmt"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
}

func (m *MetricsServer) GetMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	metrics, err := getMetrics(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while collecting metrics: %w", err)
	}
	return metrics, nil
}

func (m *MetricsServer) Connect(ctx context.Context, req *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	return &pb.ConnectionResponse{
		Status: pb.ConnectionStatus_ACCEPTED,
	}, nil
}
