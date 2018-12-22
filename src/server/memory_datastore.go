package main

import (
	"errors"
	"sync"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinMemoryDataStore struct {
	idles     []*pb.Idle
	taskQueue []*pb.Transaction
	m         sync.RWMutex
}

func (datastore *uzucoinMemoryDataStore) getIdles() []*pb.Idle {
	datastore.m.RLock()
	defer datastore.m.RUnlock()

	return datastore.idles
}

func (datastore *uzucoinMemoryDataStore) getTasks() []*pb.Transaction {
	datastore.m.RLock()
	defer datastore.m.RUnlock()

	return datastore.taskQueue
}

func (datastore *uzucoinMemoryDataStore) addTask(task *pb.Transaction) {
	datastore.m.Lock()
	defer datastore.m.Unlock()

	datastore.taskQueue = append(datastore.taskQueue, task)
}

func isEqualTransaction(lhs *pb.Transaction, rhs *pb.Transaction) bool {
	return lhs.FromUID == rhs.FromUID &&
		lhs.ToUID == rhs.ToUID &&
		lhs.Amount == rhs.Amount &&
		lhs.Timestamp == rhs.Timestamp
}

func (datastore *uzucoinMemoryDataStore) addIdle(idle *pb.Idle) error {
	datastore.m.Lock()
	defer datastore.m.Unlock()

	if len(datastore.taskQueue) == 0 {
		return errors.New("task queue is empty")
	}

	task := datastore.taskQueue[0]
	if !isEqualTransaction(task, idle.Transaction) {
		return errors.New("the head of task and given idle are different")
	}

	datastore.idles = append(datastore.idles, idle)
	datastore.taskQueue = datastore.taskQueue[1:]

	return nil
}

func initUzucoinMemoryDataStore() (uzucoinDataStore, error) {
	datastore := &uzucoinMemoryDataStore{
		idles:     make([]*pb.Idle, 0),
		taskQueue: make([]*pb.Transaction, 0),
		m:         sync.RWMutex{},
	}

	return datastore, nil
}
