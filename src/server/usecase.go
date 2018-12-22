package main

import (
	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinRepository interface {
}

type uzucoinUsecase interface {
	registerUser(*pb.RegisterUserRequest) (*pb.RegisterUserResponse, error)
	getHistory(*pb.GetHistoryRequest) (*pb.History, error)
	getBalance(*pb.GetBalanceRequest) (*pb.Balance, error)
	addTransaction(*pb.TransactionRequest) (*pb.AddTransactionResponse, error)
	getTask(*pb.GetTaskRequest) (*pb.Task, error)
	resolveNonce(*pb.Nonce) (*pb.ResolveNonceResponse, error)
}

type uzucoinUsecaseImpl struct {
	repository uzucoinRepository
}

func (usecase *uzucoinUsecaseImpl) registerUser(request *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return nil, nil
}

func (usecase *uzucoinUsecaseImpl) addTransaction(request *pb.TransactionRequest) (*pb.AddTransactionResponse, error) {
	return nil, nil
}

func (usecase *uzucoinUsecaseImpl) getHistory(request *pb.GetHistoryRequest) (*pb.History, error) {
	return nil, nil
}

func (usecase *uzucoinUsecaseImpl) getBalance(request *pb.GetBalanceRequest) (*pb.Balance, error) {
	return nil, nil
}

func (usecase *uzucoinUsecaseImpl) getTask(request *pb.GetTaskRequest) (*pb.Task, error) {
	return nil, nil
}

func (usecase *uzucoinUsecaseImpl) resolveNonce(nonce *pb.Nonce) (*pb.ResolveNonceResponse, error) {
	return nil, nil
}

func initUzucoinUsecase(repository uzucoinRepository) (uzucoinUsecase, error) {
	usecase := &uzucoinUsecaseImpl{repository: repository}

	return usecase, nil
}
