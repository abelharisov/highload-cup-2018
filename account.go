package main

import (
	"github.com/thoas/go-funk"
)

// Account ...
type Account struct {
	Id         *int      `bson:"id" json:"id"`
	FName      *string   `bson:"fname,omitempty" json:"fname"`
	SName      *string   `bson:"sname,omitempty" json:"sname"`
	Phone      *string   `bson:"phone,omitempty" json:"phone"`
	Email      *string   `bson:"email" json:"email"`
	Sex        *string   `bson:"sex,omitempty" json:"sex"`
	Birth      *int      `bson:"birth,omitempty" json:"birth"`
	BirthYear  *int      `bson:"birthYear,omitempty"  json:"-"` // custom field
	Country    *string   `bson:"country,omitempty" json:"country"`
	City       *string   `bson:"city,omitempty" json:"city"`
	Joined     *int      `bson:"joined,omitempty" json:"joined"`
	JoinedYear *int      `bson:"joinedYear,omitempty" json:"-"`
	Interests  *[]string `bson:"interests,omitempty" json:"interests"`
	Status     *string   `bson:"status,omitempty" json:"status"`
	Premium    *struct {
		Start  int `bson:"start,omitempty" json:"start"`
		Finish int `bson:"finish,omitempty" json:"finish"`
	} `bson:"premium,omitempty" json:"premium"`
	Likes *[]struct {
		Id int `bson:"id,omitempty" json:"id"`
		Ts int `bson:"ts,omitempty" json:"ts"`
	} `bson:"likes,omitempty" json:"likes"`
	LikeIds *[]int `bson:"likeIds,omitempty" json:"-"` // custom
}

// type Likers []int
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

// type Ids []int
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
