package main

import (
	pb "../pb"
	"crypto/sha256"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const rs2Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString2(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rs2Letters[rand.Intn(len(rs2Letters))]
	}
	return string(b)
}

type Idle struct {
	transaction *pb.FixedTransaction
	nonce       string
	prevHash    string
}

func CalcIdleHash(idle *Idle) (string, error) {
	payload := fmt.Sprintf("%s%s%s", idle.transaction.Timestamp, idle.nonce, idle.prevHash)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))

	return hash, nil
}

func main() {
	if len(os.Args) <= 1 {
		os.Exit(1)
	}

	conn, err := grpc.Dial("localhost:11111", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Dial:", err)
	}
	defer conn.Close()
	c := pb.NewAPIClient(conn)

	mode := os.Args[1]
	if mode == "add" {
		amount, _ := strconv.ParseFloat(os.Args[2], 64)
		msg := &pb.Transaction{
			FromUID: "from uid",
			ToUID:   "to uid",
			Amount:  amount,
		}
		if res, err := c.AddTransaction(context.Background(), msg); err != nil {
			panic(err)
		} else {
			fmt.Println("add:", res)
		}
	} else if mode == "get" {
		msg := &pb.GetTaskRequest{}
		res, err := c.GetTask(context.Background(), msg)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("task:", res)
		}
	} else if mode == "res" {
		msg := &pb.GetTaskRequest{}
		res, err := c.GetTask(context.Background(), msg)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("task:", res)
		}

		rand.Seed(time.Now().UnixNano())
		var nonce string
		for {
			nonce_cand := RandString2(10)
			hash, _ := CalcIdleHash(&Idle{res.Transaction, nonce_cand, res.PrevHash})
			fmt.Println(hash)
			if strings.Contains(hash, "757a756b69") ||
				strings.Contains(hash, "757a75") ||
				strings.Contains(hash, "7a756b69") ||
				strings.Contains(hash, "75") ||
				strings.Contains(hash, "7a75") ||
				strings.Contains(hash, "6b69") {
				fmt.Println("found", hash)
				nonce = nonce_cand
				break
			}
		}

		msg2 := &pb.Nonce{
			PrevHash: res.PrevHash,
			Nonce:    nonce,
		}
		if res, err := c.ResolveNonce(context.Background(), msg2); err != nil {
			panic(err)
		} else {
			fmt.Println("resolve:", res)
		}
	}
}
