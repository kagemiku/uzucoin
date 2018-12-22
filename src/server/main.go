package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

const (
	port = ":50051"
	seed = "Uzuki Shimamura"
)

var handler pb.UzucoinServer

func init() {
	datastore, err := initUzucoinMemoryDataStore()
	if err != nil {
		log.Fatalf("failed to init datastore: %v\n", err)
	}

	repository, err := initUzucoinRepository(datastore)
	if err != nil {
		log.Fatalf("failed to init repository: %v\n", err)
	}

	initialHash := fmt.Sprintf(hashFormat, sha256.Sum256([]byte(seed)))
	usecase, err := initUzucoinUsecase(initialHash, repository)
	if err != nil {
		log.Fatalf("failed to init repository: %v\n", err)
	}

	handler, err = initUzucoinHandler(usecase)
	if err != nil {
		log.Fatalf("failed to init repository: %v\n", err)
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUzucoinServer(s, handler)
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Printf("server started with port: %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
