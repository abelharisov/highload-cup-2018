package main

// Account ...
type Account struct {
	Id        *int32    `bson:"id"`
	FName     *string   `bson:"fname,omitempty"`
	SName     *string   `bson:"sname,omitempty"`
	Phone     *string   `bson:"phone,omitempty"`
	Email     *string   `bson:"email"`
	Sex       *string   `bson:"sex,omitempty"`
	Birth     *int32    `bson:"birth,omitempty"`
	Year      *int      `bson:"year,omitempty"` // custom field
	Country   *string   `bson:"country,omitempty"`
	City      *string   `bson:"city,omitempty"`
	Joined    *int32    `bson:"joined,omitempty"`
	Interests *[]string `bson:"interests,omitempty"`
	Status    *string   `bson:"status,omitempty"`
	Premium   *struct {
		Start  int32 `bson:"start,omitempty"`
		Finish int32 `bson:"finish,omitempty"`
	} `bson:"premium,omitempty"`
	Likes *[]struct {
		Id int32 `bson:"id,omitempty"`
		Ts int32 `bson:"ts,omitempty"`
	} `bson:"likes,omitempty"`
	LikeIds *[]int32 `bson:"likeIds,omitempty"`
}
