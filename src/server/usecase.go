package main

import (
	"crypto/sha256"
	"fmt"
	"time"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinRepository interface {
	getIdelsCount() int
	getLatestIdle() *Idle
	getHeadTask() *pb.Transaction
	addTask(*pb.Transaction)
}

type uzucoinUsecase interface {
	getHistory(*pb.GetHistoryRequest) (*pb.History, error)
	getBalance(*pb.GetBalanceRequest) (*pb.Balance, error)
	addTransaction(*pb.AddTransactionRequest) (*pb.AddTransactionResponse, error)
	getTask(*pb.GetTaskRequest) (*pb.Task, error)
	resolveNonce(*pb.Nonce) (*pb.ResolveNonceResponse, error)
}

type uzucoinUsecaseImpl struct {
	initialHash string
	repository  uzucoinRepository
}

const (
	payloadFormat = "%s%s%s"
	hashFormat    = "%x"
)

func calcIdleHash(idle *Idle) string {
	payload := fmt.Sprintf(payloadFormat, idle.transaction.Timestamp, idle.nonce, idle.prevHash)
	hash := fmt.Sprintf(hashFormat, sha256.Sum256([]byte(payload)))

	return hash
}

func (usecase *uzucoinUsecaseImpl) getHistory(request *pb.GetHistoryRequest) (*pb.History, error) {
	return nil, nil
}

func (usecase *uzucoinUsecaseImpl) getBalance(request *pb.GetBalanceRequest) (*pb.Balance, error) {
	return nil, nil
}

func (usecase *uzucoinUsecaseImpl) addTransaction(request *pb.AddTransactionRequest) (*pb.AddTransactionResponse, error) {
	timestamp := time.Now().String()
	task := &pb.Transaction{
		FromUID:   request.FromUID,
		ToUID:     request.ToUID,
		Amount:    request.Amount,
		Timestamp: timestamp,
	}
	usecase.repository.addTask(task)

	return &pb.AddTransactionResponse{Timestamp: timestamp}, nil
}

func (usecase *uzucoinUsecaseImpl) getTask(request *pb.GetTaskRequest) (*pb.Task, error) {
	var task *pb.Task
	transaction := usecase.repository.getHeadTask()
	if transaction == nil {
		task = &pb.Task{
			Exists:      false,
			Transaction: nil,
			PrevHash:    "",
		}
	} else if idlesCount := usecase.repository.getIdelsCount(); idlesCount == 0 {
		task = &pb.Task{
			Exists:      true,
			Transaction: transaction,
			PrevHash:    usecase.initialHash,
		}
	} else {
		prevHash := calcIdleHash(usecase.repository.getLatestIdle())
		task = &pb.Task{
			Exists:      true,
			Transaction: transaction,
			PrevHash:    prevHash,
		}
	}

	return task, nil
}

func (usecase *uzucoinUsecaseImpl) resolveNonce(nonce *pb.Nonce) (*pb.ResolveNonceResponse, error) {
	return nil, nil
}

func initUzucoinUsecase(initialHash string, repository uzucoinRepository) (uzucoinUsecase, error) {
	usecase := &uzucoinUsecaseImpl{
		initialHash: initialHash,
		repository:  repository,
	}

	return usecase, nil
}
