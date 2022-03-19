package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/lodthe/mafia/internal/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DefaultMaxRetries = 5
const DefaultRetryTimeout = 3 * time.Second

func main() {
	log.SetFlags(0)

	var address string
	flag.StringVar(&address, "address", "127.0.0.1:9000", "server address")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var username string
	fmt.Printf("Enter your username: ")
	_, _ = fmt.Scanf("%s", &username)

	conn, err := createConnection(address)
	if err != nil {
		log.Fatalf("failed to connect to %s: %v\n", address, err)
	}
	defer conn.Close()

	cli, err := client.NewClient(ctx, username, conn)
	if err != nil {
		log.Fatalf("failed to init gRPC client: %v\n", err)
	}

	messenger := client.NewMessenger(ctx)
	engine := client.NewEngine(ctx, cli, messenger)

	go messenger.Start()
	go engine.Start()
	go cli.ForwardEvents()

	<-stop
	cancel()
}

func createConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithMax(DefaultMaxRetries),
			grpc_retry.WithPerRetryTimeout(DefaultRetryTimeout),
		)),
	)

	return conn, err
}
