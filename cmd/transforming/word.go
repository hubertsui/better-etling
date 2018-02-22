package main

type Word struct {
	name  string
	count int
}

type Words []Word

func (c Words) Len() int {
	return len(c)
}
func (c Words) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Words) Less(i, j int) bool {
	return c[i].count > c[j].count
}
