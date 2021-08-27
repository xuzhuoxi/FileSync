package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"io/fs"
	"sort"
)

// pathList

func newPathStrList(ln int, cap int) *pathStrList {
	return &pathStrList{ItemArray: make([]string, ln, cap)}
}

type pathStrList struct {
	ItemArray []string
}

func (l *pathStrList) Len() int {
	return len(l.ItemArray)
}

func (l *pathStrList) Less(i, j int) bool {
	lenI := getRuneCount(l.ItemArray[i], infra.DirSeparator)
	lenJ := getRuneCount(l.ItemArray[j], infra.DirSeparator)
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

func (l *pathStrList) Append(path string) {
	if "" == path {
		return
	}
	l.ItemArray = append(l.ItemArray, path)
}

func (l *pathStrList) AppendList(list *pathStrList) {
	if nil == list {
		return
	}
	l.ItemArray = append(l.ItemArray, list.ItemArray...)
}

func (l *pathStrList) Remove(index int) string {
	if index < 0 || index >= len(l.ItemArray) {
		return ""
	}
	rs := l.ItemArray[index]
	l.ItemArray = append(l.ItemArray[:index], l.ItemArray[index+1:]...)
	return rs
}

// detailPath & detailPathList

func newDetailPathList(ln int, cap int) detailPathList {
	return make(detailPathList, ln, cap)
}

type detailPath struct {
	Index    int
	SrcInfo  infra.SrcInfo // Src信息
	FileInfo fs.FileInfo

	SrcRelativePath string // 运行时 源 相对路径
	TarRelativePath string // 运行时 目录 相对路径

	SrcAbsPath string // 源 绝对路径
	TarAbsPath string // 目标 绝对路径
}

type detailPathList []detailPath

func (l detailPathList) Len() int {
	return len(l)
}

func (l detailPathList) Less(i, j int) bool {
	if l[i].Index != l[j].Index {
		return l[i].Index < l[j].Index
	}
	lenI := getRuneCount(l[i].SrcRelativePath, infra.DirSeparator)
	lenJ := getRuneCount(l[j].SrcRelativePath, infra.DirSeparator)
	if lenI != lenJ {
		return lenI > lenJ
	}
	return l[i].SrcRelativePath < l[j].SrcRelativePath
}

func (l detailPathList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l detailPathList) Sort() { sort.Sort(l) }

// pathInfo

type pathInfo struct {
	RelativeRoot string      // 基于 运行时 或 配置文件目录 的相对路径
	SubPath      string      // 基于 RelativeRoot 的相对路径
	FileInfo     fs.FileInfo // 文件或目录信息

	relativePath string
	fullPath     string
}

func (i *pathInfo) GetRelativePath() string {
	if "" == i.relativePath {
		i.relativePath = filex.Combine(i.RelativeRoot, i.SubPath)
	}
	return i.relativePath
}

func (i *pathInfo) GetFullPath() string {
	if "" == i.fullPath {
		i.fullPath = filex.Combine(infra.RunningDir, i.RelativeRoot, i.SubPath)
	}
	return i.fullPath
}

func newPathInfoList(ln int, cap int) *pathInfoList {
	return &pathInfoList{ItemArray: make([]*pathInfo, ln, cap)}
}

type pathInfoList struct {
	ItemArray []*pathInfo
}

func (l *pathInfoList) Len() int {
	return len(l.ItemArray)
}

func (l *pathInfoList) Less(i, j int) bool {
	rpi := l.ItemArray[i].GetRelativePath()
	rpj := l.ItemArray[j].GetRelativePath()
	lenI := getRuneCount(rpi, infra.DirSeparator)
	lenJ := getRuneCount(rpj, infra.DirSeparator)
	if lenI != lenJ {
		return lenI > lenJ
	}
	return rpi < rpj
}

func (l *pathInfoList) Swap(i, j int) {
	l.ItemArray[i], l.ItemArray[j] = l.ItemArray[j], l.ItemArray[i]
}

func (l *pathInfoList) Sort() { sort.Sort(l) }

func (l *pathInfoList) Get(index int) *pathInfo {
	if index < 0 || index >= len(l.ItemArray) {
		return nil
	}
	return l.ItemArray[index]
}

func (l *pathInfoList) Append(path *pathInfo) {
	if nil == path {
		return
	}
	l.ItemArray = append(l.ItemArray, path)
}

func (l *pathInfoList) AppendList(list *pathInfoList) {
	if nil == list {
		return
	}
	l.ItemArray = append(l.ItemArray, list.ItemArray...)
}

func (l *pathInfoList) Remove(index int) *pathInfo {
	if index < 0 || index >= len(l.ItemArray) {
		return nil
	}
	rs := l.ItemArray[index]
	l.ItemArray = append(l.ItemArray[:index], l.ItemArray[index+1:]...)
	return rs
}
