package main

import (
	"errors"
	"math/bits"
	"sort"
)

type account struct {
	country    uint
	city       uint
	birth      int
	premium    int
	interestsA uint64
	interestsB uint64
}

const countryI = 0
const cityI = 8

type AccountRecIndex struct {
	InterestDict *Dict
	CountryDict  *Dict
	CityDict     *Dict

	fFree map[int]account
	fBusy map[int]account
	fWtf  map[int]account
	mFree map[int]account
	mBusy map[int]account
	mWtf  map[int]account
}

func (i *AccountRecIndex) Init() {
	i.fFree = make(map[int]account, 0)
	i.fBusy = make(map[int]account, 0)
	i.fWtf = make(map[int]account, 0)

	i.mFree = make(map[int]account, 0)
	i.mBusy = make(map[int]account, 0)
	i.mWtf = make(map[int]account, 0)
}

func (i *AccountRecIndex) getCollection(sex string, status int) *map[int]account {
	if sex == "f" {
		if status == StatusFree {
			return &i.fFree
		} else if status == StatusWtf {
			return &i.fWtf
		} else if status == StatusBusy {
			return &i.fBusy
		}
	}

	if sex == "m" {
		if status == StatusFree {
			return &i.mFree
		} else if status == StatusWtf {
			return &i.mWtf
		} else if status == StatusBusy {
			return &i.mBusy
		}
	}

	return nil
}

func (i *AccountRecIndex) Add(a Account) {
	col := i.getCollection(*a.Sex, a.StatusId)
	if col == nil {
		panic("accounts collection is nil")
	}

	binary := account{}

	if a.City == nil {
		binary.city = 0
	} else {
		binary.city = i.CityDict.GetId(*a.City)
	}

	if a.Country == nil {
		binary.country = 0
	} else {
		binary.country = i.CountryDict.GetId(*a.Country)
	}

	binary.birth = *a.Birth
	binary.premium = a.PremiumStatus

	if a.Interests != nil {
		for _, interest := range *a.Interests {
			id := i.InterestDict.GetId(interest)
			var byte = uint64(1) << id
			if id > 63 {
				binary.interestsB += byte
			} else {
				binary.interestsA += byte
			}
		}
	}

	(*col)[a.Id] = binary
}

type AccountRecArray []AccountRec

type AccountRec struct {
	Id    int
	Score int
}

func (a AccountRecArray) Len() int {
	return len(a)
}

func (a AccountRecArray) Less(i, j int) bool {
	return a[i].Score > a[j].Score
}

func (a AccountRecArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

const PremiumСoefficient__ = 100000000000000
const StatusСoefficient___ = 1000000000000000
const InterestsСoefficient = 10000000000
const MaxBirthScore_______ = 2147483648

func (i *AccountRecIndex) Recommend(a Account, country string, city string, limit int) (result []int) {
	sex := "f"
	if a.Sex == nil {
		panic(errors.New("wtf!!!!"))
	}
	if *a.Sex == "f" {
		sex = "m"
	}

	target := (*i.getCollection(*a.Sex, a.StatusId))[a.Id]
	countryId := i.CountryDict.GetId(country)
	cityId := i.CityDict.GetId(city)

	ids := make([]AccountRec, 0, 100)

	for status := StatusFree; status >= StatusBusy; status-- {
		col := *i.getCollection(sex, status)

		statusScore := status * StatusСoefficient___

		for id, binary := range col {
			if countryId != 0 && binary.country != countryId {
				continue
			}

			if cityId != 0 && binary.city != cityId {
				continue
			}

			if ((target.interestsA & binary.interestsA) | (target.interestsB & binary.interestsB)) != 0 {
				birthScore := MaxBirthScore_______ - int(abs(int64(target.birth)-int64(binary.birth)))
				premiumScore := PremiumСoefficient__
				if binary.premium != PremiumActive {
					premiumScore = 0
				}
				interestsScore := (bits.OnesCount64(target.interestsA) + bits.OnesCount64(target.interestsB)) * InterestsСoefficient
				ids = append(ids, AccountRec{
					id,
					statusScore + interestsScore + birthScore + premiumScore,
				})
			}
		}
	}

	sort.Sort(AccountRecArray(ids))

	// log.Println(limit)
	result = make([]int, 0, limit)

	for i := 0; i < limit && i < len(ids); i++ {
		result = append(result, ids[i].Id)
	}

	return
}

func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
