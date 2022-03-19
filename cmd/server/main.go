package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lodthe/mafia/internal/server"
	"github.com/lodthe/mafia/pkg/mafiapb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	DefaultTimeout   = 10 * time.Second
	DefaultKeepAlive = 500 * time.Millisecond
)

func main() {
	var address string
	flag.StringVar(&address, "address", "0.0.0.0:9000", "address for listening")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	engine, err := server.NewEngine(ctx, server.DefaultConfig)
	if err != nil {
		log.Fatalf("failed to create engine: %v\nconfig: %v\n", err, server.DefaultConfig)
	}

	srv, lis, err := registerServer(ctx, address, engine)
	if err != nil {
		log.Fatalf("server registration failed on %s: %v\n", address, err)
	}

	go func() {
		err := srv.Serve(lis)
		if err != nil {
			log.Fatalf("server failed: %v\n", err)
		}
	}()

	log.Printf("server started on %s\n", address)

	<-stop
	cancel()
}

func registerServer(ctx context.Context, address string, engine *server.Engine) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, errors.Wrap(err, "listen failed")
	}

	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(DefaultTimeout),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: DefaultKeepAlive,
			Time:              DefaultKeepAlive,
			Timeout:           DefaultKeepAlive,
		}),
		StdUnaryMiddleware(),
		StdStreamMiddleware(),
	)

	mafiapb.RegisterMafiaServer(grpcServer, server.NewServer(ctx, engine))

	reflection.Register(grpcServer)

	return grpcServer, lis, nil
}
