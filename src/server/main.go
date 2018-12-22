package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) AddTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.AddTransactionResponse, error) {
	return nil, nil
}

func (s *server) GetTask(ctx context.Context, in *pb.GetTaskRequest) (*pb.Task, error) {

	return nil, nil
}

func (s *server) ResolveNonce(ctx context.Context, in *pb.Nonce) (*pb.ResolveNonceResponse, error) {

	return nil, nil
}

func (s *server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return nil, nil
}

func (s *server) GetHistory(ctx context.Context, in *pb.GetHistoryRequest) (*pb.History, error) {

	return nil, nil
}

func (s *server) GetBalance(ctx context.Context, in *pb.GetBalanceRequest) (*pb.Balance, error) {

	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUzucoinServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
