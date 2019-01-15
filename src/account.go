package main

// Account ...
type Account struct {
	Id int32
	FName string
	SName string
	Phone string
	Email string
	Sex string
	Birth int32
	Country string
	City string
	Joined int32
	Interests []string
	Status string
	Premium struct {
		Start int32
		Finish int32
	}
	Likes []struct{
		Id int32
		Ts int32
	}
}
