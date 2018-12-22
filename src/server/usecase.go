package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinRepository interface {
	registerProducer(*Producer) error
	getProducer(string) (*Producer, error)
	getIdelsCount() int
	getLatestIdle() *pb.Idle
	getIdles() []*pb.Idle
	addIdle(*pb.Idle) error
	getTransactions() []*pb.Transaction
	getHeadTask() *pb.Transaction
	addTask(*pb.Transaction)
}

type uzucoinUsecase interface {
	registerProducer(*pb.RegisterProducerRequest) (*pb.RegisterProducerResponse, error)
	getHistory(*pb.GetHistoryRequest) (*pb.History, error)
	getBalance(*pb.GetBalanceRequest) (*pb.Balance, error)
	getChain(*pb.GetChainRequest) (*pb.Chain, error)
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

func calcIdleHash(idle *pb.Idle) string {
	payload := fmt.Sprintf(payloadFormat, idle.Transaction.Timestamp, idle.Nonce, idle.PrevHash)
	hash := fmt.Sprintf(hashFormat, sha256.Sum256([]byte(payload)))

	return hash
}

func calcUzucoin(hash string) (bool, float64) {
	succeeded := false
	reward := 0.0
	if strings.Contains(hash, asciiUzuki) {
		succeeded = true
		reward = 24.0
	} else if strings.Contains(hash, asciiUzu) || strings.Contains(hash, asciiZuki) {
		succeeded = true
		reward = 4.0
	} else if strings.Contains(hash, asciiU) || strings.Contains(hash, asciiZu) || strings.Contains(hash, asciiKi) {
		succeeded = true
		reward = 1.0
	}

	return succeeded, reward
}

func (usecase *uzucoinUsecaseImpl) registerProducer(request *pb.RegisterProducerRequest) (*pb.RegisterProducerResponse, error) {
	producer := &Producer{
		uid:  request.Uid,
		name: request.Name,
	}

	if err := usecase.repository.registerProducer(producer); err != nil {
		return &pb.RegisterProducerResponse{Succeeded: false}, err
	}

	return &pb.RegisterProducerResponse{Succeeded: true}, nil
}

func (usecase *uzucoinUsecaseImpl) getHistory(request *pb.GetHistoryRequest) (*pb.History, error) {
	uid := request.Uid
	if _, err := usecase.repository.getProducer(uid); err != nil {
		return &pb.History{Transactions: []*pb.Transaction{}}, err
	}

	transactions := usecase.repository.getTransactions()
	history := make([]*pb.Transaction, 0)
	for _, transaction := range transactions {
		if transaction.FromUID == uid || transaction.ToUID == uid {
			history = append(history, transaction)
		}
	}

	return &pb.History{Transactions: history}, nil
}

func (usecase *uzucoinUsecaseImpl) calcBalance(uid string) (float64, error) {
	if _, err := usecase.repository.getProducer(uid); err != nil {
		return 0.0, err
	}

	idles := usecase.repository.getIdles()
	balance := initialUzucoin
	for _, idle := range idles {
		if idle.Transaction.FromUID == uid {
			balance -= idle.Transaction.Amount
		}

		if idle.Transaction.ToUID == uid {
			balance += idle.Transaction.Amount
		}

		if idle.ResolverUID == uid {
			hash := calcIdleHash(idle)
			succeeded, reward := calcUzucoin(hash)
			if !succeeded {
				return 0.0, errors.New("chain may be broken")
			}

			balance += reward
		}
	}

	return balance, nil

}

func (usecase *uzucoinUsecaseImpl) getBalance(request *pb.GetBalanceRequest) (*pb.Balance, error) {
	balance, err := usecase.calcBalance(request.Uid)
	if err != nil {
		return &pb.Balance{Balance: 0.0}, err
	}

	return &pb.Balance{Balance: balance}, nil
}

func (usecase *uzucoinUsecaseImpl) getChain(request *pb.GetChainRequest) (*pb.Chain, error) {
	idles := usecase.repository.getIdles()

	return &pb.Chain{Idles: idles}, nil
}

func (usecase *uzucoinUsecaseImpl) addTransaction(request *pb.AddTransactionRequest) (*pb.AddTransactionResponse, error) {
	if _, err := usecase.repository.getProducer(request.FromUID); err != nil {
		return &pb.AddTransactionResponse{Timestamp: ""}, err
	}

	if _, err := usecase.repository.getProducer(request.ToUID); err != nil {
		return &pb.AddTransactionResponse{Timestamp: ""}, err
	}

	balance, err := usecase.calcBalance(request.FromUID)
	if err != nil {
		return &pb.AddTransactionResponse{Timestamp: ""}, err
	}

	if request.Amount > balance {
		return &pb.AddTransactionResponse{Timestamp: ""}, errors.New("lack of balance")
	}

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

func (usecase *uzucoinUsecaseImpl) resolveNonce(request *pb.ResolveNonceRequest) (*pb.ResolveNonceResponse, error) {
	if _, err := usecase.repository.getProducer(request.ResolverUID); err != nil {
		return &pb.ResolveNonceResponse{Succeeded: false, Reward: 0.0}, err
	}

	var prevHash string
	if usecase.repository.getIdelsCount() == 0 {
		prevHash = usecase.initialHash
	} else {
		prevHash = calcIdleHash(usecase.repository.getLatestIdle())
	}

	if request.PrevHash != prevHash {
		return &pb.ResolveNonceResponse{Succeeded: false, Reward: 0.0}, nil
	}

	transaction := usecase.repository.getHeadTask()
	idle := &pb.Idle{
		Transaction: transaction,
		Nonce:       request.Nonce,
		PrevHash:    request.PrevHash,
		ResolverUID: request.ResolverUID,
	}
	newHash := calcIdleHash(idle)
	succeeded, reward := calcUzucoin(newHash)
	if !succeeded {
		return &pb.ResolveNonceResponse{Succeeded: false, Reward: 0.0}, nil
	}

	if err := usecase.repository.addIdle(idle); err != nil {
		return &pb.ResolveNonceResponse{Succeeded: false, Reward: 0.0}, err
	}

	return &pb.ResolveNonceResponse{Succeeded: true, Reward: reward}, nil
}

func initUzucoinUsecase(initialHash string, repository uzucoinRepository) (uzucoinUsecase, error) {
	usecase := &uzucoinUsecaseImpl{
		initialHash: initialHash,
		repository:  repository,
	}

	return usecase, nil
}
