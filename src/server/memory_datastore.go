package main

import (
	"fmt"
	"sync"

	pb "github.com/kagemiku/uzucoin/src/server/pb"
)

type uzucoinMemoryDataStore struct {
	idles     []*Idle
	taskQueue []*pb.Transaction
	m         sync.RWMutex
}

func (datastore *uzucoinMemoryDataStore) getIdles() []*Idle {
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

func initUzucoinMemoryDataStore() (uzucoinDataStore, error) {
	datastore := &uzucoinMemoryDataStore{
		idles:     make([]*Idle, 0),
		taskQueue: make([]*pb.Transaction, 0),
		m:         sync.RWMutex{},
	}

	return datastore, nil
}
