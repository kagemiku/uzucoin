package main

import pb "github.com/kagemiku/uzucoin/src/server/pb"

type uzucoinDataStore interface {
	getIdles() []*pb.Idle
	addIdle(*pb.Idle) error
	getTasks() []*pb.Transaction
	addTask(*pb.Transaction)
}

type uzucoinRepositoryImpl struct {
	datastore uzucoinDataStore
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
