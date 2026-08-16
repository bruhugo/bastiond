package server

import (
	pb "bastiond/generated"
	"context"
	"fmt"
	"log/slog"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	config *BastionConfig
}

func NewMetricsServer(config *BastionConfig) *MetricsServer {
	return &MetricsServer{
		config: config,
	}
}

func (m *MetricsServer) GetMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	fmt.Print(m.config)
	metrics, err := getMetrics(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while collecting metrics: %w", err)
	}
	return metrics, nil
}

func (m *MetricsServer) Connect(ctx context.Context, req *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	fmt.Print(m.config)
	if m.config.Password != "" && req.Password != m.config.Password {
		slog.Error("connection failed to authenticate")
		return &pb.ConnectionResponse{
			Status: pb.ConnectionStatus_DENIED,
		}, nil
	}

	return &pb.ConnectionResponse{
		Status: pb.ConnectionStatus_ACCEPTED,
	}, nil
}
