package main

// Storage interface
type Storage interface {
	LoadAccounts(accounts []Account)
	Find(filter *AccountsFilter) []map[string]interface{}
}
