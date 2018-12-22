package main

import pb "github.com/kagemiku/uzucoin/src/server/pb"

type Idle struct {
	transaction *pb.Transaction
	nonce       string
	prevHash    string
}
