package main

import pb "github.com/kagemiku/uzucoin/src/server/pb"

type uzucoinDataStore interface {
	getIdles() []*Idle
	getTasks() []*pb.Transaction
}

type uzucoinRepositoryImpl struct {
	datastore uzucoinDataStore
}

func (repository *uzucoinRepositoryImpl) getIdelsCount() int {
	idles := repository.datastore.getIdles()

	return len(idles)
}

func (repository *uzucoinRepositoryImpl) getLatestIdle() *Idle {
	idles := repository.datastore.getIdles()
	if len(idles) == 0 {
		return nil
	}

	return idles[len(idles)-1]
}

func (repository *uzucoinRepositoryImpl) getHeadTask() *pb.Transaction {
	tasks := repository.datastore.getTasks()
	if len(tasks) == 0 {
		return nil
	}

	return tasks[0]
}

func initUzucoinRepository(datastore uzucoinDataStore) (uzucoinRepository, error) {
	repository := &uzucoinRepositoryImpl{datastore: datastore}

	return repository, nil
}
