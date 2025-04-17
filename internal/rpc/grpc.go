package rpc

import (
	"context"
	"fmt"
	"github.com/ether-echo/user-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"

	pb "github.com/ether-echo/protos/userProcessor"
)

var log = logger.Logger().Named("GRPC").Sugar()

type GrpcServer struct {
	conn              *grpc.ClientConn
	userServiceClient pb.UserServiceClient
}

func NewGrpcServer(addr string) *GrpcServer {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}

	userServiceClient := pb.NewUserServiceClient(conn)
	return &GrpcServer{
		conn:              conn,
		userServiceClient: userServiceClient,
	}
}

func (g *GrpcServer) StartMessage(chatID int64, firstName string, exist bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := g.userServiceClient.StartMessage(ctx, &pb.StartRequest{
		ChatId:    chatID,
		FirstName: firstName,
		Exist:     exist,
	})

	if err != nil {
		return fmt.Errorf("could not send notification: %v", err)
	}

	log.Infof("Notification sent successfully: %v", resp.Success)

	return nil
}

func (g *GrpcServer) Close() {
	err := g.conn.Close()
	if err != nil {
		log.Errorf("could not close gRPC connection: %v", err)
	}
}
