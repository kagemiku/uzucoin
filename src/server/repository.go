package main

import (
	"errors"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinDataStore interface {
	getProducers() []*Producer
	addProducer(*Producer)
	getIdles() []*pb.Idle
	addIdle(*pb.Idle) error
	getTransactions() []*pb.Transaction
	getTasks() []*pb.Transaction
	addTask(*pb.Transaction)
}

type uzucoinRepositoryImpl struct {
	datastore uzucoinDataStore
}

func (repository *uzucoinRepositoryImpl) registerProducer(producer *Producer) error {
	producers := repository.datastore.getProducers()
	for _, p := range producers {
		if *producer == *p {
			return errors.New("Producer is already registered")
		}
	}

	repository.datastore.addProducer(producer)

	return nil
}

func (repository *uzucoinRepositoryImpl) getProducer(uid string) (*Producer, error) {
	producers := repository.datastore.getProducers()
	for _, producer := range producers {
		if producer.uid == uid {
			return producer, nil
		}
	}

	return nil, errors.New("No such Producuer")
}

func (repository *uzucoinRepositoryImpl) getIdelsCount() int {
	idles := repository.datastore.getIdles()

	return len(idles)
}

func (repository *uzucoinRepositoryImpl) getLatestIdle() *pb.Idle {
	idles := repository.datastore.getIdles()
	if len(idles) == 0 {
		return nil
	}

	return idles[len(idles)-1]
}

func (repository *uzucoinRepositoryImpl) getIdles() []*pb.Idle {
	return repository.datastore.getIdles()
}

func (repository *uzucoinRepositoryImpl) addIdle(idle *pb.Idle) error {
	return repository.datastore.addIdle(idle)
}

func (repository *uzucoinRepositoryImpl) getTransactions() []*pb.Transaction {
	return repository.datastore.getTransactions()
}

func (repository *uzucoinRepositoryImpl) getHeadTask() *pb.Transaction {
	tasks := repository.datastore.getTasks()
	if len(tasks) == 0 {
		return nil
	}

	return tasks[0]
}

func (repository *uzucoinRepositoryImpl) addTask(task *pb.Transaction) {
	repository.datastore.addTask(task)
}

func initUzucoinRepository(datastore uzucoinDataStore) (uzucoinRepository, error) {
	repository := &uzucoinRepositoryImpl{datastore: datastore}

	return repository, nil
}
