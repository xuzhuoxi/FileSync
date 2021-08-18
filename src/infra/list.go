package infra

import "sort"

func NewPathList(ln int, cap int) PathList {
	list := make([]string, ln, cap)
	return PathList(list)
}

type PathList []string

func (l PathList) Len() int {
	return len(l)
}

func (l PathList) Less(i, j int) bool {
	lenI := len(l[i])
	lenJ := len(l[j])
	if lenI != lenJ {
		return lenI < lenJ
	}
	return l[i] < l[j]
}

func (l PathList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l PathList) Sort() { sort.Sort(l) }
