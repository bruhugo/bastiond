package server

import (
	pb "bastiond/generated"
	"context"
	"log/slog"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

func getMetrics(ctx context.Context) (*pb.Metrics, error) {
	cpu, err := getCPUMetrics(ctx)
	if err != nil {
		return nil, err
	}

	mem, err := getMemoryMetrics(ctx)
	if err != nil {
		return nil, err
	}

	net, err := getNetworkMetrics(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.Metrics{
		CpuMetrics:     cpu,
		MemoryMetrics:  mem,
		NetworkMetrics: net,
	}, nil
}

func getCPUMetrics(ctx context.Context) (*pb.CPUMetrics, error) {
	prev, _ := cpu.Times(false)

	time.Sleep(time.Second)

	curr, _ := cpu.Times(false)

	a := prev[0]
	b := curr[0]

	user := b.User - a.User
	system := b.System - a.System
	idle := b.Idle - a.Idle
	iowait := b.Iowait - a.Iowait
	irq := b.Irq - a.Irq
	softirq := b.Softirq - a.Softirq
	nice := b.Nice - a.Nice
	steal := b.Steal - a.Steal

	total :=
		user + system + idle +
			iowait + irq + softirq +
			nice + steal

	return &pb.CPUMetrics{
		IdleUtilization:   idle / total,
		SystemUtilization: system / total,
		IOWait:            iowait / total,
		UserUtilization:   user / total,
	}, nil
}

func getMemoryMetrics(ctx context.Context) (*pb.MemoryMetrics, error) {
	stat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		slog.Error("error getting memory status", slog.String("error", err.Error()))
	}

	return &pb.MemoryMetrics{
		TotalRam:          int64(stat.Total),
		UsedRam:           int64(stat.Used),
		AvailableRam:      int64(stat.Available),
		MemoryUtilization: stat.UsedPercent,
		TotalSwap:         int64(stat.SwapTotal),
		UsedSwap:          int64(stat.Total - stat.Free),
	}, nil
}

func getNetworkMetrics(ctx context.Context) (*pb.NetworkMetrics, error) {
	stats, err := net.IOCountersWithContext(ctx, false)
	if err != nil {
		slog.Error("error getting network settings", slog.String("error", err.Error()))
		return nil, err
	}
	stat := stats[0]

	return &pb.NetworkMetrics{
		RXBytes:    int64(stat.BytesRecv),
		RXPackets:  int64(stat.PacketsRecv),
		TXBytes:    int64(stat.BytesSent),
		TXPackets:  int64(stat.PacketsSent),
		PacketLoss: int64(stat.Dropout + stat.Dropin),
	}, nil
}
