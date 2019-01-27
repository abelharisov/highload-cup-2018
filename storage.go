package main

// Storage interface
type Storage interface {
	CreateIndexes()
	DropIndexes()
	LoadAccounts(accounts []Account)
	Find(query *AccountsQuery) ([]map[string]interface{}, error)
	Group(query *AccountsGroupQuery) ([]map[string]interface{}, error)
	Recommend(q *AccountsRecommendQuery) ([]map[string]interface{}, error)
	Suggset(q *AccountsRecommendQuery) (result []map[string]interface{}, err error)
}
