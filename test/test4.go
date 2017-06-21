package main

import (
	"log"
	"sort"
)

type CommonInterface interface {
	getCmpKey() interface{}
}

type By func(p1, p2 CommonInterface) bool

func (by By) Sort(inters []CommonInterface) {
	ps := &InterSorter{
		Inters: inters,
		by:     by,
	}
	sort.Sort(ps)
}

type InterSorter struct {
	Inters []CommonInterface
	by     func(p1, p2 CommonInterface) bool
}

func (s *InterSorter) Len() int {
	return len(s.Inters)
}

func (s *InterSorter) Swap(i, j int) {
	s.Inters[i], s.Inters[j] = s.Inters[j], s.Inters[i]
}

func (s *InterSorter) Less(i, j int) bool {
	return s.by(s.Inters[i], s.Inters[j])
}

type A struct {
	num  int
	name string
	CommonInterface
}

func (a A) getCmpKey() interface{} {
	return a.num
}

func main() {

	myNum := func(p1, p2 CommonInterface) bool {
		return p1.getCmpKey().(int) < p2.getCmpKey().(int)
	}

	var aa []CommonInterface

	aa = append(aa, A{
		num:  10,
		name: "a1",
	})

	aa = append(aa, A{
		num:  5,
		name: "a2",
	})

	aa = append(aa, A{
		num:  8,
		name: "a3",
	})

	log.Println(aa)

	By(myNum).Sort(aa)

	log.Println(aa)

}
