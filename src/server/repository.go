package main

type uzucoinDataStore interface {
}

type uzucoinRepositoryImpl struct {
	datastore uzucoinDataStore
}

func initUzucoinRepository(datastore uzucoinDataStore) (uzucoinRepository, error) {
	repository := &uzucoinRepositoryImpl{datastore: datastore}

	return repository, nil
}
