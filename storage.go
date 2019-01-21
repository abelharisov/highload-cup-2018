package main

// Storage interface
type Storage interface {
	LoadAccounts(accounts []Account)
	Find(query *AccountsQuery) ([]map[string]interface{}, error)
}
