package main

import funk "github.com/thoas/go-funk"

type InterestsMap map[string][]int

func (m *InterestsMap) Append(accountId int, interests []string) {
	for _, interest := range interests {
		if _, ok := (*m)[interest]; !ok {
			(*m)[interest] = make([]int, 0)
		}
		(*m)[interest] = append((*m)[interest], accountId)
	}
}

func (m *InterestsMap) AccountsWithInterestsAny(interests []string) (ids []int) {
	for _, interest := range interests {
		if _, ok := (*m)[interest]; ok {
			ids = append(ids, (*m)[interest]...)
		}
	}

	ids = funk.UniqInt(ids)

	return
}

func (m *InterestsMap) AccountsWithInterestsContains(interests []string) (ids []int) {
	var result *[]int
	for _, interest := range interests {
		if news, ok := (*m)[interest]; ok {
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
