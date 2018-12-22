package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

const (
	port = ":50051"
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

	usecase, err := initUzucoinUsecase(repository)
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
