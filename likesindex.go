package main

type Likes map[int]([]int)

type LikesIndex struct {
	m map[int]Likes
}

func CreateLikesIndex() *LikesIndex {
	return &LikesIndex{
		make(map[int]Likes),
	}
}

func (i *LikesIndex) HasLike(id int, likeId int) bool {
	_, ok := (*i).m[id][likeId]
	return ok
}

func (i *LikesIndex) AddLikes(a Account) {
	if a.Likes == nil {
		i.m[a.Id] = Likes{}
		return
	}

	for _, like := range *a.Likes {
		if _, ok := i.m[a.Id]; !ok {
			i.m[a.Id] = Likes{}
		}
		if _, ok := i.m[a.Id][like.Id]; !ok {
			i.m[a.Id][like.Id] = []int{}
		}
		i.m[a.Id][like.Id] = append(i.m[a.Id][like.Id], like.Ts)
	}
}

func (i *LikesIndex) GetLikes(id int) *Likes {
	if likes, ok := i.m[id]; ok {
		return &likes
	} else {
		return nil
	}
}
