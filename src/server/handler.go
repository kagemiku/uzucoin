package main

import (
	"context"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type server struct {
	usecase *uzucoinUsecase
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

func (s *server) AddTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.AddTransactionResponse, error) {
	return nil, nil
}

func (s *server) GetTask(ctx context.Context, in *pb.GetTaskRequest) (*pb.Task, error) {

	return nil, nil
}

func (s *server) ResolveNonce(ctx context.Context, in *pb.Nonce) (*pb.ResolveNonceResponse, error) {

	return nil, nil
}
