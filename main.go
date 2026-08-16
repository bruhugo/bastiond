package main

import (
	pb "bastiond/generated"
	server "bastiond/internal"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func getLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		panic("invalid log level provided")
	}
}

func setLogger(level string) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(level),
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}

func main() {
	port := flag.Int("port", 8080, "Port to listen for requests")
	log := flag.String("log", "INFO", "The log level")
	password := flag.String("pass", "", "The password for the client to use. Omit if no password is required.")

	tlsPath := flag.String("cert", "", "The path poiting to your certificate file")
	keyPath := flag.String("key", "", "The path for the private key")

	flag.Parse()

	setLogger(*log)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	config := &server.BastionConfig{
		TLSCertificatePath: *tlsPath,
		Password:           *password,
	}

	var opts []grpc.ServerOption

	// add certificate
	if *tlsPath != "" {
		credential, err := credentials.NewServerTLSFromFile(*tlsPath, *keyPath)
		if err != nil {
			panic(err)
		}

		opts = append(opts, grpc.Creds(credential))
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterMetricsServiceServer(grpcServer, server.NewMetricsServer(config))

	slog.Info("metrics server implementation was registered")
	slog.Info("server is ready", slog.Int("port", *port))

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
