package main

// Account ...
type Account struct {
	Id         int       `bson:"id" json:"id"`
	FName      *string   `bson:"fname,omitempty" json:"fname"`
	SName      *string   `bson:"sname,omitempty" json:"sname"`
	Phone      *string   `bson:"phone,omitempty" json:"phone"`
	Email      *string   `bson:"email" json:"email"`
	Sex        *string   `bson:"sex,omitempty" json:"sex"`
	Birth      *int      `bson:"birth,omitempty" json:"birth"`
	Country    *string   `bson:"country,omitempty" json:"country"`
	City       *string   `bson:"city,omitempty" json:"city"`
	Joined     *int      `bson:"joined,omitempty" json:"joined"`
	JoinedYear *int      `bson:"joinedYear,omitempty" json:"-"`
	Interests  *[]string `bson:"-" json:"interests"`
	Status     string    `bson:"status,omitempty" json:"status"`
	Premium    *struct {
		Start  int64 `bson:"start,omitempty" json:"start"`
		Finish int64 `bson:"finish,omitempty" json:"finish"`
	} `bson:"premium,omitempty" json:"premium"`
	Likes *[]struct {
		Id int `bson:"-" json:"id"`
		Ts int `bson:"-" json:"ts"`
	} `bson:"-" json:"likes"`

	BirthYear     *int   `bson:"birthYear,omitempty"  json:"-"` // custom field
	LikeIds       *[]int `bson:"-" json:"-"`                    // custom
	StatusId      int    `bson:"statusId,omitempty" json:"-"`   // custom
	PremiumStatus int    `bson:"premiumStatus" json:"-"`        // custom
	PhoneCode     int    `bson:"phoneCode" json:"-"`            // custom
}

const PremiumNull = 1
const PremiumNotActive = 2
const PremiumActive = 3

const StatusFree = 3
const StatusWtf = 2
const StatusBusy = 1
