package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Server interface {
	AddTransaction(*Transaction) (*AddTransactionResponse, error)
	GetTask(*GetTaskRequest) (*Task, error)
	ResolveNonce(*Nonce) (*ResolveNonceResponse, error)
}

type ServerImpl struct {
	Idles            []*Idle
	TransactionQueue []*FixedTransaction
	M                sync.Mutex
}

type Idle struct {
	Transaction *FixedTransaction
	Nonce       string
	PrevHash    string
}

const seed = "Uzuki Shimamura"

var initialHash = fmt.Sprintf("%x", sha256.Sum256([]byte(seed)))

func (this *ServerImpl) AddTransaction(t *Transaction) (*AddTransactionResponse, error) {
	this.M.Lock()
	defer this.M.Unlock()
	fmt.Printf("%f うづコイン送られたよ！\n", t.Amount)

	timestamp := time.Now().String()
	this.TransactionQueue = append(this.TransactionQueue, &FixedTransaction{t.FromUID, t.ToUID, t.Amount, timestamp})
	//fmt.Println("current queue", s.transactionQueue)
	//fmt.Println()

	return &AddTransactionResponse{timestamp}, nil

}

func (this *ServerImpl) CalcIdleHash(idle *Idle) (string, error) {
	payload := fmt.Sprintf("%s%s%s", idle.Transaction.Timestamp, idle.Nonce, idle.PrevHash)
	//fmt.Println("payload:", payload)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	//fmt.Println("hash:", hash)
	//fmt.Println()

	return hash, nil
}

func (this *ServerImpl) GetTask(p *GetTaskRequest) (*Task, error) {
	this.M.Lock()
	defer this.M.Unlock()
	fmt.Println("S(min)ING!してね♪")

	if len(this.TransactionQueue) == 0 {
		task := &Task{false, &FixedTransaction{}, ""}
		return task, nil
	} else if len(this.Idles) == 0 {
		task := &Task{true, this.TransactionQueue[0], initialHash}
		return task, nil
	}

	prevHash, _ := this.CalcIdleHash(this.Idles[len(this.Idles)-1])
	task := &Task{true, this.TransactionQueue[0], prevHash}
	return task, nil
}

func (this *ServerImpl) ShowChain() {
	fmt.Println("いまのうづコインのーとです")
	for index, idle := range this.Idles {
		fmt.Printf("[%d], %s\n", index, idle)
	}
}

func (this *ServerImpl) ResolveNonce(p *Nonce) (*ResolveNonceResponse, error) {
	this.M.Lock()
	defer this.M.Unlock()
	fmt.Printf("ありがとー！")

	var prevHash string
	if len(this.Idles) == 0 {
		prevHash = initialHash
	} else {
		prevHash, _ = this.CalcIdleHash(this.Idles[len(this.Idles)-1])
	}

	if p.PrevHash != prevHash {
		return &ResolveNonceResponse{false, 0.0}, nil
	}

	transaction := this.TransactionQueue[0]
	amount := 0.0
	succeeded := false
	newHash, _ := this.CalcIdleHash(&Idle{transaction, p.Nonce, p.PrevHash})
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
		fmt.Printf(" %f うづコインあげるね♪\n", amount)
		this.TransactionQueue = this.TransactionQueue[1:]
		this.Idles = append(this.Idles, &Idle{transaction, p.Nonce, p.PrevHash})
		this.ShowChain()
		fmt.Println()
	}

	return &ResolveNonceResponse{succeeded, amount}, nil
}
