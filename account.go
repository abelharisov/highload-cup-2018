package main

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

type Likers []int
type LikeeMap map[int]Likers

func (likeeMap *LikeeMap) AppendLiker(likeeId int, likerId int) {
	if _, ok := (*likeeMap)[likeeId]; !ok {
		(*likeeMap)[likeeId] = make(Likers, 0)
	}
	(*likeeMap)[likeeId] = append((*likeeMap)[likeeId], likerId)
}
