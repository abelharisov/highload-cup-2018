package main

import (
	"github.com/google/btree"
)

type LikesIndex struct {
	m *btree.BTree
}

type idsMap map[int]float32

type Likes struct {
	id  int
	ids idsMap
}

func (l Likes) Less(o btree.Item) bool {
	return l.id < o.(Likes).id
}

func (l Likes) Len() int {
	return len(l.ids)
}

func CreateLikesIndex() *LikesIndex {
	return &LikesIndex{
		btree.New(8),
	}
}

func (i *LikesIndex) HasLike(id int, likeId int) bool {
	item := i.m.Get(Likes{id: id})
	if item != nil {
		return false
	}

	_, ok := item.(Likes).ids[likeId]
	return ok
}

func (i *LikesIndex) AddLikes(a Account) {
	if a.Likes == nil {
		i.m.ReplaceOrInsert(Likes{a.Id, idsMap{}})
		return
	}

	for _, like := range *a.Likes {
		if !i.m.Has(Likes{id: a.Id}) {
			i.m.ReplaceOrInsert(Likes{a.Id, idsMap{like.Id: float32(like.Ts)}})
		} else {
			item := i.m.Get(Likes{id: a.Id}).(Likes)
			if ts, ok := item.ids[like.Id]; ok {
				item.ids[like.Id] = (ts + float32(like.Ts)) / 2.0
			} else {
				item.ids[like.Id] = float32(like.Ts)
			}
		}
	}
}

func (i *LikesIndex) GetLikes(id int) *Likes {
	item := i.m.Get(Likes{id: id})
	if item == nil {
		return nil
	}

	likes := item.(Likes)
	return &likes
}
