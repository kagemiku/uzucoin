package main

import (
	"context"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinHandler struct {
	usecase uzucoinUsecase
}

func (handler *uzucoinHandler) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return nil, nil
}

func (handler *uzucoinHandler) GetHistory(ctx context.Context, in *pb.GetHistoryRequest) (*pb.History, error) {

	return nil, nil
}

func (handler *uzucoinHandler) GetBalance(ctx context.Context, in *pb.GetBalanceRequest) (*pb.Balance, error) {

	return nil, nil
}

func (handler *uzucoinHandler) AddTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.AddTransactionResponse, error) {
	return nil, nil
}

func (handler *uzucoinHandler) GetTask(ctx context.Context, in *pb.GetTaskRequest) (*pb.Task, error) {

	return nil, nil
}

func (handler *uzucoinHandler) ResolveNonce(ctx context.Context, in *pb.Nonce) (*pb.ResolveNonceResponse, error) {

	return nil, nil
}

func initUzucoinHandler(usecase uzucoinUsecase) (pb.UzucoinServer, error) {
	handler := &uzucoinHandler{usecase: usecase}

	return handler, nil
}
