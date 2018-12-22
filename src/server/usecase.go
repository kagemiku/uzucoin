package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinRepository interface {
	getIdelsCount() int
	getLatestIdle() *Idle
	addIdle(*Idle) error
	getHeadTask() *pb.Transaction
	addTask(*pb.Transaction)
}

type uzucoinUsecase interface {
	getHistory(*pb.GetHistoryRequest) (*pb.History, error)
	getBalance(*pb.GetBalanceRequest) (*pb.Balance, error)
	addTransaction(*pb.AddTransactionRequest) (*pb.AddTransactionResponse, error)
	getTask(*pb.GetTaskRequest) (*pb.Task, error)
	resolveNonce(*pb.ResolveNonceRequest) (*pb.ResolveNonceResponse, error)
}

type uzucoinUsecaseImpl struct {
	initialHash string
	repository  uzucoinRepository
}

const (
	payloadFormat = "%s%s%s"
	hashFormat    = "%x"
	asciiUzuki    = "757a756b69"
	asciiUzu      = "757a75"
	asciiZuki     = "7a756b69"
	asciiU        = "75"
	asciiZu       = "7a75"
	asciiKi       = "6b69"
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

func calcUzucoin(hash string) (bool, float64) {
	succeeded := false
	amount := 0.0
	if strings.Contains(hash, asciiUzuki) {
		succeeded = true
		amount = 24.0
	} else if strings.Contains(hash, asciiUzu) || strings.Contains(hash, asciiZuki) {
		succeeded = true
		amount = 4.0
	} else if strings.Contains(hash, asciiU) || strings.Contains(hash, asciiZu) || strings.Contains(hash, asciiKi) {
		succeeded = true
		amount = 1.0
	}

	return succeeded, amount
}

func (usecase *uzucoinUsecaseImpl) resolveNonce(request *pb.ResolveNonceRequest) (*pb.ResolveNonceResponse, error) {
	var prevHash string
	if usecase.repository.getIdelsCount() == 0 {
		prevHash = usecase.initialHash
	} else {
		prevHash = calcIdleHash(usecase.repository.getLatestIdle())
	}

	if request.PrevHash != prevHash {
		return &pb.ResolveNonceResponse{Succeeded: false, Amount: 0.0}, nil
	}

	transaction := usecase.repository.getHeadTask()
	idle := &Idle{
		transaction: transaction,
		nonce:       request.Nonce,
		prevHash:    request.PrevHash,
	}
	newHash := calcIdleHash(idle)
	succeeded, amount := calcUzucoin(newHash)
	if !succeeded {
		return &pb.ResolveNonceResponse{Succeeded: false, Amount: 0.0}, nil
	}

	if err := usecase.repository.addIdle(idle); err != nil {
		return &pb.ResolveNonceResponse{Succeeded: false, Amount: 0.0}, err
	}

	return &pb.ResolveNonceResponse{Succeeded: true, Amount: amount}, nil
}

func initUzucoinUsecase(initialHash string, repository uzucoinRepository) (uzucoinUsecase, error) {
	usecase := &uzucoinUsecaseImpl{
		initialHash: initialHash,
		repository:  repository,
	}

	return usecase, nil
}
