package server

import (
	pb "bastiond/generated"
	"context"
	"log/slog"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

func getMetrics(ctx context.Context, request *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	cpu, err := getCPUMetrics(ctx, request.MultipleCores)
	if err != nil {
		return nil, err
	}

	mem, err := getMemoryMetrics(ctx)
	if err != nil {
		return nil, err
	}

	net, err := getNetworkMetrics(ctx, request.MultipleInterfaces)
	if err != nil {
		return nil, err
	}

	disk, err := getDiskMetrics(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.MetricsResponse{
		CpuMetrics:     cpu,
		MemoryMetrics:  mem,
		NetworkMetrics: net,
		DiskMetrics:    disk,
	}, nil
}

func getCPUMetrics(ctx context.Context, multipleCores bool) (*pb.CPUMetrics, error) {
	prev, err := cpu.TimesWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	prevGen, err := cpu.TimesWithContext(ctx, false)
	if err != nil {
		return nil, err
	}

	time.Sleep(500 * time.Millisecond)

	curr, err := cpu.TimesWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	currGen, err := cpu.TimesWithContext(ctx, false)
	if err != nil {
		return nil, err
	}

	if len(prev) == 0 || len(curr) == 0 {
		return &pb.CPUMetrics{IsMultipleCores: false}, nil
	}

	list := make([]*pb.CoreCPUMetrics, 0, 1)
	for i, stat := range curr {
		if i >= len(prev) {
			break
		}

		metrics := getCpuMetricsFromStats(prev[i], stat)
		if metrics == nil {
			continue
		}

		list = append(list, metrics)
	}

	slog.Debug("collected cpu metrics")

	return &pb.CPUMetrics{
		IsMultipleCores: false,
		Cores:           list,
		General:         getCpuMetricsFromStats(prevGen[0], currGen[0]),
	}, nil
}

func getCpuMetricsFromStats(b, a cpu.TimesStat) *pb.CoreCPUMetrics {
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

	if total == 0 {
		return nil
	}

	return &pb.CoreCPUMetrics{
		UserUtilization:   user / total,
		SystemUtilization: system / total,
		IdleUtilization:   idle / total,
		IoWait:            iowait / total,
		Core:              a.CPU,
	}
}

func getMemoryMetrics(ctx context.Context) (*pb.MemoryMetrics, error) {
	stat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		slog.Error("error getting memory status", slog.String("error", err.Error()))
	}

	slog.Debug("collected memory metrics")

	return &pb.MemoryMetrics{
		TotalRam:          stat.Total,
		UsedRam:           stat.Used,
		AvailableRam:      stat.Available,
		MemoryUtilization: stat.UsedPercent,
		TotalSwap:         stat.SwapTotal,
		UsedSwap:          stat.Total - stat.Free,
	}, nil
}

func getNetworkMetrics(ctx context.Context, pernic bool) (*pb.NetworkMetrics, error) {
	stats, err := net.IOCountersWithContext(ctx, pernic)
	if err != nil {
		slog.Error("error getting network settings", slog.String("error", err.Error()))
		return nil, err
	}

	if len(stats) == 0 {
		return &pb.NetworkMetrics{IsMultipleInterfaces: false}, nil
	}

	list := make([]*pb.NetworkInterfaceMetrics, 0)

	for _, stat := range stats {
		list = append(list, &pb.NetworkInterfaceMetrics{
			RxBytes:       stat.BytesRecv,
			TxBytes:       stat.BytesSent,
			RxPackets:     stat.PacketsRecv,
			TxPackets:     stat.PacketsSent,
			InterfaceName: stat.Name,
		})
	}

	slog.Debug("collected network metrics",
		slog.Int("interfaces", len(stats)),
		slog.Bool("pernic", pernic),
	)

	return &pb.NetworkMetrics{
		IsMultipleInterfaces:    false,
		NetworkInterfaceMetrics: list,
	}, nil
}

func getDiskMetrics(ctx context.Context) (*pb.DiskMetrics, error) {
	usage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		slog.Error("error collecting disk metrics", slog.String("error", err.Error()))
		return nil, err
	}

	partitions, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		slog.Error("error collecting disk metrics", slog.String("error", err.Error()))
		return nil, err
	}

	partitionsMetrics := make([]*pb.PartitionMetrics, 0)

	for _, partition := range partitions {
		partitionsMetrics = append(partitionsMetrics, &pb.PartitionMetrics{
			Device:     partition.Device,
			Mountpoint: partition.Mountpoint,
			Fstype:     partition.Fstype,
		})
	}

	return &pb.DiskMetrics{
		TotalDisk:   usage.Total,
		UsedDisk:    usage.Used,
		FreeDisk:    usage.Free,
		DiskUsage:   usage.UsedPercent,
		InodesTotal: usage.InodesTotal,
		InodesUsed:  usage.InodesUsed,
		InodesFree:  usage.InodesFree,
	}, nil
}
