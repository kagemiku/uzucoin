package main

import (
	"context"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinHandler struct {
	usecase uzucoinUsecase
}

func (handler *uzucoinHandler) RegisterProducer(ctx context.Context, in *pb.RegisterProducerRequest) (*pb.RegisterProducerResponse, error) {
	return handler.usecase.registerProducer(in)
}

func (handler *uzucoinHandler) GetHistory(ctx context.Context, in *pb.GetHistoryRequest) (*pb.History, error) {
	return handler.usecase.getHistory(in)
}

func (handler *uzucoinHandler) GetBalance(ctx context.Context, in *pb.GetBalanceRequest) (*pb.Balance, error) {
	return handler.usecase.getBalance(in)
}

func (handler *uzucoinHandler) GetChain(ctx context.Context, in *pb.GetChainRequest) (*pb.Chain, error) {
	return handler.usecase.getChain(in)
}

func (handler *uzucoinHandler) AddTransaction(ctx context.Context, in *pb.AddTransactionRequest) (*pb.AddTransactionResponse, error) {
	return handler.usecase.addTransaction(in)
}

func (handler *uzucoinHandler) GetTask(ctx context.Context, in *pb.GetTaskRequest) (*pb.Task, error) {
	return handler.usecase.getTask(in)
}

func (handler *uzucoinHandler) ResolveNonce(ctx context.Context, in *pb.ResolveNonceRequest) (*pb.ResolveNonceResponse, error) {
	return handler.usecase.resolveNonce(in)
}

func initUzucoinHandler(usecase uzucoinUsecase) (pb.UzucoinServer, error) {
	handler := &uzucoinHandler{usecase: usecase}

	return handler, nil
}
