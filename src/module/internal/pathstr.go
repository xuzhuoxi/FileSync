package internal

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"sort"
)

// pathStrList

type IPathStrList interface {
	sort.Interface
	// 排序
	Sort()
	// 取全部元素
	GetAll() []string
	// 取元素
	Get(index int) string
	// 追加元素
	Append(path string)
	// 追加元素列表
	AppendArray(arr []string)
	// 追加元素列表
	AppendList(list IPathStrList)
	// 移除元素
	Remove(index int) string
}

func NewPathStrList(ln int, cap int) IPathStrList {
	return &pathStrList{ItemArray: make([]string, ln, cap)}
}

type pathStrList struct {
	ItemArray []string
}

func (l *pathStrList) Len() int {
	return len(l.ItemArray)
}

func (l *pathStrList) Less(i, j int) bool {
	lenI := infra.GetRuneCount(l.ItemArray[i], infra.DirSeparator)
	lenJ := infra.GetRuneCount(l.ItemArray[j], infra.DirSeparator)
	if lenI != lenJ {
		return lenI > lenJ
	}
	return l.ItemArray[i] < l.ItemArray[j]
}

func (l *pathStrList) Swap(i, j int) {
	l.ItemArray[i], l.ItemArray[j] = l.ItemArray[j], l.ItemArray[i]
}

func (l *pathStrList) Sort() { sort.Sort(l) }

func (l *pathStrList) Get(index int) string {
	if index < 0 || index >= len(l.ItemArray) {
		return ""
	}
	return l.ItemArray[index]
}

func (l *pathStrList) GetAll() []string {
	return l.ItemArray
}

func (l *pathStrList) Append(path string) {
	if "" == path {
		return
	}
	l.ItemArray = append(l.ItemArray, path)
}

func (l *pathStrList) AppendArray(arr []string) {
	if len(arr) == 0 {
		return
	}
	l.ItemArray = append(l.ItemArray, arr...)
}

func (l *pathStrList) AppendList(list IPathStrList) {
	if nil == list {
		return
	}
	l.AppendArray(list.GetAll())
}

func (l *pathStrList) Remove(index int) string {
	if index < 0 || index >= len(l.ItemArray) {
		return ""
	}
	rs := l.ItemArray[index]
	l.ItemArray = append(l.ItemArray[:index], l.ItemArray[index+1:]...)
	return rs
}
