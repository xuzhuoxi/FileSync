package internal

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"io/fs"
	"sort"
)

type IPathInfo interface {
	// 基于(运行目录/配置文件目录) 的相对路径
	GetRelativeRoot() string
	// 基于RelativeRoot的相对路径
	GetRootSubPath() string
	// 文件或目录信息
	GetFileInfo() fs.FileInfo

	// 完整相对路径
	GetRelativePath() string
	// 完整绝对路径
	GetFullPath() string

	// 基于新的RelativeRoot生成完整相对路径
	GenRelativePath(relativeRoot string) string
	// 基于新的RelativeRoot生成完整绝对路径
	GenFullPath(relativeRoot string) string

	// 比较
	LessTo(target IPathInfo) bool
}

type pathInfo struct {
	RelativeRoot string      // 基于(运行目录/配置文件目录) 的相对路径
	RootSubPath  string      // 基于RelativeRoot的相对路径
	FileInfo     fs.FileInfo // 文件或目录信息

	relativePath string // 临时补全的完整相对路径,基于(运行目录/配置文件目录)
	fullPath     string // 临时补全的完整绝对路径
}

func (i *pathInfo) GetRelativeRoot() string {
	return i.RelativeRoot
}

func (i *pathInfo) GetRootSubPath() string {
	return i.RootSubPath
}

func (i *pathInfo) GetFileInfo() fs.FileInfo {
	return i.FileInfo
}

func (i *pathInfo) GetRelativePath() string {
	if "" == i.relativePath {
		i.relativePath = filex.Combine(i.RelativeRoot, i.RootSubPath)
	}
	return i.relativePath
}

func (i *pathInfo) GetFullPath() string {
	if "" == i.fullPath {
		i.fullPath = filex.Combine(infra.RunningDir, i.RelativeRoot, i.RootSubPath)
	}
	return i.fullPath
}

func (i *pathInfo) GenRelativePath(relativeRoot string) string {
	return filex.Combine(relativeRoot, i.RootSubPath)
}

func (i *pathInfo) GenFullPath(relativeRoot string) string {
	return filex.Combine(infra.RunningDir, relativeRoot, i.RootSubPath)
}

func (i *pathInfo) LessTo(target IPathInfo) bool {
	dirI := i.FileInfo.IsDir()
	dirJ := target.GetFileInfo().IsDir()
	rpi := i.GetRelativePath()
	rpj := target.GetRelativePath()
	if dirI == dirJ {
		return rpi < rpj
	}
	return dirJ
}

type IPathInfoList interface {
	sort.Interface
	// 排序
	Sort()
	// 取元素
	Get(index int) IPathInfo
	// 取全部元素
	GetAll() []IPathInfo
	// 追加元素
	Append(path IPathInfo)
	// 追加元素列表
	AppendArray(arr []IPathInfo)
	// 追加元素列表
	AppendList(list IPathInfoList)
	// 移除元素
	Remove(index int) IPathInfo
}

func NewPathInfoList(ln int, cap int) IPathInfoList {
	return &pathInfoList{ItemArray: make([]IPathInfo, ln, cap)}
}

type pathInfoList struct {
	ItemArray []IPathInfo
}

func (l *pathInfoList) Len() int {
	return len(l.ItemArray)
}

func (l *pathInfoList) Less(i, j int) bool {
	return l.ItemArray[i].LessTo(l.ItemArray[j])
}

func (l *pathInfoList) Swap(i, j int) {
	l.ItemArray[i], l.ItemArray[j] = l.ItemArray[j], l.ItemArray[i]
}

func (l *pathInfoList) Sort() { sort.Sort(l) }

func (l *pathInfoList) Get(index int) IPathInfo {
	if index < 0 || index >= len(l.ItemArray) {
		return nil
	}
	return l.ItemArray[index]
}

func (l *pathInfoList) GetAll() []IPathInfo {
	return l.ItemArray
}

func (l *pathInfoList) Append(path IPathInfo) {
	if nil == path {
		return
	}
	l.ItemArray = append(l.ItemArray, path)
}

func (l *pathInfoList) AppendArray(arr []IPathInfo) {
	if len(arr) == 0 {
		return
	}
	l.ItemArray = append(l.ItemArray, arr...)
}

func (l *pathInfoList) AppendList(list IPathInfoList) {
	if nil == list {
		return
	}
	l.AppendArray(list.GetAll())
}

func (l *pathInfoList) Remove(index int) IPathInfo {
	if index < 0 || index >= len(l.ItemArray) {
		return nil
	}
	rs := l.ItemArray[index]
	l.ItemArray = append(l.ItemArray[:index], l.ItemArray[index+1:]...)
	return rs
}
