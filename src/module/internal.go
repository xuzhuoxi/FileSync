package module

import (
	"errors"
	"github.com/xuzhuoxi/infra-go/filex"
	"io/fs"
	"os"
	"sort"
	"time"
)

type iModuleExecutor interface {
	// 初始化Log
	initArgs()
	// 初始化处理列表
	initExecuteList()
	// 处理列表中文件
	execList()
}

// struct

func newPathList(ln int, cap int) pathList {
	return make(pathList, ln, cap)
}

type pathList []string

func (l pathList) Len() int {
	return len(l)
}

func (l pathList) Less(i, j int) bool {
	lenI := getRuneCount(l[i], filex.UnixSeparator)
	lenJ := getRuneCount(l[j], filex.UnixSeparator)
	if lenI != lenJ {
		return lenI > lenJ
	}
	return l[i] < l[j]
}

func (l pathList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l pathList) Sort() { sort.Sort(l) }

type detailPath struct {
	Index        int
	SrcRoot      string // Src根目录
	RelativePath string // 对应以src为根目录的相对路径
	FileInfo     fs.FileInfo

	SrcRelativePath string // 运行时源相对路径
	TarRelativePath string // 运行时目录相对路径

	SrcAbsPath string // 源绝对路径
	TarAbsPath string // 目标绝对路径
}

func newDetailPathList(ln int, cap int) detailPathList {
	return make(detailPathList, ln, cap)
}

type detailPathList []detailPath

func (l detailPathList) Len() int {
	return len(l)
}

func (l detailPathList) Less(i, j int) bool {
	if l[i].Index != l[j].Index {
		return l[i].Index < l[j].Index
	}
	lenI := getRuneCount(l[i].SrcRelativePath, filex.UnixSeparator)
	lenJ := getRuneCount(l[j].SrcRelativePath, filex.UnixSeparator)
	if lenI != lenJ {
		return lenI > lenJ
	}
	return l[i].SrcRelativePath < l[j].SrcRelativePath
}

func (l detailPathList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l detailPathList) Sort() { sort.Sort(l) }

func getRuneCount(str string, r rune) int {
	runes := []rune(str)
	rs := 0
	for _, o := range runes {
		if r == o {
			rs += 1
		}
	}
	return rs
}

func cloneTime(tarPath string, srcInfo os.FileInfo) {
	os.Chtimes(tarPath, srcInfo.ModTime(), srcInfo.ModTime())
}

func compareTime(tarPath string, mt time.Time) bool {
	fileInfo, err := os.Stat(tarPath)
	// 不存在
	if err != nil && !errors.Is(err, os.ErrExist) {
		return true
	}
	return fileInfo.ModTime().UnixNano() < mt.UnixNano()
}
