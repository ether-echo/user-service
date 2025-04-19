package main

import (
	"context"
	up "github.com/ether-echo/protos/userProcessor"
	"github.com/ether-echo/user-service/internal/adapter"
	"github.com/ether-echo/user-service/internal/kafka"
	"github.com/ether-echo/user-service/internal/repository"
	"github.com/ether-echo/user-service/internal/rpc"
	"github.com/ether-echo/user-service/internal/service"
	"github.com/ether-echo/user-service/pkg/config"
	"github.com/ether-echo/user-service/pkg/debug"
	"github.com/ether-echo/user-service/pkg/logger"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	logger.BuildLogger(conf.LogLevel)
	log := logger.Logger().Named("main").Sugar()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Infof("Starting debug server on %s", conf.DebugPort)
		debug.Run(":" + conf.DebugPort)
	}()

	grpcServer := rpc.NewGrpcServer("message-processor:50052")
	defer grpcServer.Close()

	pDB := repository.NewPostgresDB(conf)
	defer pDB.Close()

	serviceTg := service.NewService(pDB, grpcServer)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen grps_users_service: %v", err)
	}

	grpcUserServer := grpc.NewServer()
	up.RegisterUserServiceServer(grpcUserServer, &adapter.UserService{
		Repository: serviceTg,
	})

	go func() {
		if err := grpcUserServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	handlerCG := kafka.NewConsumerGroupHandler(serviceTg)

	consumer, err := kafka.NewConsumer(conf.KafkaBrokers, []string{"start"}, conf.KafkaGroup, handlerCG)
	if err != nil {
		log.Fatalf("Consumer stopped with error: %v", err)
	}

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("Consumer stopped with error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info("Gracefully shutting down...")

	cancel()

	if err := consumer.Stop(); err != nil {
		log.Fatalf("Consumer stopped with error: %v", err)
	}

	log.Info("Service stopped")
}
