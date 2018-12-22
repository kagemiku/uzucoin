package main

type uzucoinMemoryDataStore struct{}

func initUzucoinMemoryDataStore() (uzucoinDataStore, error) {
	datastore := &uzucoinMemoryDataStore{}

	return datastore, nil
}
