package main

import (
	"github.com/google/btree"
	"github.com/thoas/go-funk"
)

type likersList struct {
	likeeId int
	ids     []int
}

func (l likersList) Less(than btree.Item) bool {
	return l.likeeId < than.(likersList).likeeId
}

type LikeeToLikerIndex struct {
	btree *btree.BTree
}

func CreateLikeeToLikerIndex() *LikeeToLikerIndex {
	return &LikeeToLikerIndex{
		btree: btree.New(8),
	}
}

func (m *LikeeToLikerIndex) GetLikers(likeeId int) []int {
	item := m.btree.Get(likersList{likeeId: likeeId})
	if item != nil {
		list := item.(likersList)
		return list.ids
	}

	return nil
}

func (m *LikeeToLikerIndex) AppendLiker(likeeId int, likerId int) {
	item := m.btree.Get(likersList{likeeId: likeeId})
	if item == nil {
		m.btree.ReplaceOrInsert(likersList{likeeId, []int{likerId}})
	} else {
		list := item.(likersList)
		list.ids = append(list.ids, likeeId)
	}
}

func (m *LikeeToLikerIndex) AccountsWithLikesContains(likees []int) (ids []int) {
	var result *[]int
	for _, likee := range likees {
		item := m.btree.Get(likersList{likeeId: likee})
		if item != nil {
			list := item.(likersList)
			if result == nil {
				result = new([]int)
				*result = append(*result, list.ids...)
			} else {
				*result = funk.Intersect(*result, list.ids).([]int)
			}
		}
	}

	if result != nil {
		ids = *result
	}

	return
}

func (m *LikeeToLikerIndex) AccountsWithLikesAny(likees []int) (ids []int) {
	ids = []int{}

	for _, likee := range likees {
		item := m.btree.Get(likersList{likeeId: likee})
		if item != nil {
			list := item.(likersList)
			ids = append(ids, list.ids...)
		}
	}

	ids = funk.Uniq(ids).([]int)

	return
}
