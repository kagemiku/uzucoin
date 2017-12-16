package main

import (
	pb "../pb"
	"crypto/sha256"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const seed = "Uzuki Shimamura"

var initialHash = fmt.Sprintf("%x", sha256.Sum256([]byte(seed)))

type Idle struct {
	transaction *pb.FixedTransaction
	nonce       string
	prevHash    string
}

type Service struct {
	idles            []*Idle
	transactionQueue []*pb.FixedTransaction
	m                sync.Mutex
}

func (s *Service) AddTransaction(c context.Context, p *pb.Transaction) (*pb.AddTransactionResponse, error) {
	s.m.Lock()
	defer s.m.Unlock()
	fmt.Println("new transaction", p)

	timestamp := time.Now().String()
	s.transactionQueue = append(s.transactionQueue, &pb.FixedTransaction{p.FromUID, p.ToUID, p.Amount, timestamp})
	fmt.Println("current queue", s.transactionQueue)
	fmt.Println()

	return &pb.AddTransactionResponse{timestamp}, nil
}

func CalcIdleHash(idle *Idle) (string, error) {
	payload := fmt.Sprintf("%s%s%s", idle.transaction.Timestamp, idle.nonce, idle.prevHash)
	fmt.Println("payload:", payload)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	fmt.Println("hash:", hash)
	fmt.Println()

	return hash, nil
}

func (s *Service) GetTask(c context.Context, p *pb.GetTaskRequest) (*pb.Task, error) {
	s.m.Lock()
	defer s.m.Unlock()
	fmt.Println("get task ", p)

	if len(s.transactionQueue) == 0 {
		task := &pb.Task{false, &pb.FixedTransaction{}, ""}
		return task, nil
	} else if len(s.idles) == 0 {
		task := &pb.Task{true, s.transactionQueue[0], initialHash}
		return task, nil
	}

	prevHash, _ := CalcIdleHash(s.idles[len(s.idles)-1])
	task := &pb.Task{true, s.transactionQueue[0], prevHash}
	return task, nil
}

func ShowChain(s *Service) {
	for index, idle := range s.idles {
		fmt.Printf("[%d], %s\n", index, idle)
	}
}

func (s *Service) ResolveNonce(c context.Context, p *pb.Nonce) (*pb.ResolveNonceResponse, error) {
	s.m.Lock()
	defer s.m.Unlock()
	fmt.Println("resolve", p)

	var prevHash string
	if len(s.idles) == 0 {
		prevHash = initialHash
	} else {
		prevHash, _ = CalcIdleHash(s.idles[len(s.idles)-1])
	}

	if p.PrevHash != prevHash {
		return &pb.ResolveNonceResponse{false, 0.0}, nil
	}

	transaction := s.transactionQueue[0]
	amount := 0.0
	succeeded := false
	newHash, _ := CalcIdleHash(&Idle{transaction, p.Nonce, p.PrevHash})
	if strings.Contains(newHash, "757a756b69") {
		amount = 24.0
		succeeded = true
	} else if strings.Contains(newHash, "757a75") || strings.Contains(newHash, "7a756b69") {
		amount = 4.0
		succeeded = true
	} else if strings.Contains(newHash, "75") || strings.Contains(newHash, "7a75") || strings.Contains(newHash, "6b69") {
		amount = 1.0
		succeeded = true
	}

	if succeeded {
		s.transactionQueue = s.transactionQueue[1:]
		s.idles = append(s.idles, &Idle{transaction, p.Nonce, p.PrevHash})
		fmt.Println("current chain")
		ShowChain(s)
		fmt.Println()
	}

	return &pb.ResolveNonceResponse{succeeded, amount}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":11111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	fmt.Println(initialHash)

	pb.RegisterAPIServer(server, new(Service))
	server.Serve(lis)
}
