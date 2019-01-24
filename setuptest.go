package main

import (
	"testing"
)

func SetupTest(t *testing.T) (teardown func(), storage *MongoStorage) {
	storage = &MongoStorage{
		Uri:      MongoUri,
		Database: "hl_test",
	}
	storage.Init()
	Parse(DataFile, OptionsFile, storage, true)

	// Test teardown - return a closure for use by 'defer'
	teardown = func() {
		storage.DropDatabase()
	}

	return
}
