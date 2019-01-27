package main

import funk "github.com/thoas/go-funk"

type LikeeMap map[int][]int

func (m *LikeeMap) AppendLiker(likeeId int, likerId int) {
	if _, ok := (*m)[likeeId]; !ok {
		(*m)[likeeId] = make([]int, 0)
	}
	(*m)[likeeId] = append((*m)[likeeId], likerId)
}

func (m *LikeeMap) AccountsWithLikesContains(likees []int) (ids []int) {
	var result *[]int
	for _, likee := range likees {
		if news, ok := (*m)[likee]; ok {
			if result == nil {
				result = &[]int{}
				*result = append(*result, news...)
			} else {
				*result = funk.Intersect(*result, news).([]int)
			}
		}
	}

	if result != nil {
		ids = *result
	}

	return
}

func (m *LikeeMap) AccountsWithLikesAny(likees []int) (ids []int) {
	ids = []int{}

	for _, likee := range likees {
		if likers, ok := (*m)[likee]; ok {
			ids = append(ids, likers...)
		}
	}

	ids = funk.Uniq(ids).([]int)

	return
}
