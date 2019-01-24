package main

type Dict struct {
	m map[string]uint
	// r map[int]string
}

func (d *Dict) Init() {
	d.m = make(map[string]uint)
	d.m[""] = 0
	// d.r = make(map[int]string)
}

func (d *Dict) GetId(v string) uint {
	if id, ok := d.m[v]; ok {
		return id
	}

	id := uint(len(d.m))
	d.m[v] = id
	// d.r[id] = v
	return id
}
