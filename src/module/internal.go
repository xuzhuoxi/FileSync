package module

import (
	"errors"
	"github.com/xuzhuoxi/FileSync/src/infra"
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
	index        int
	relativeBase string
	relativePath string
	fileInfo     fs.FileInfo
}

func (dp detailPath) GetFullPath() string {
	return filex.Combine(infra.RunningDir, dp.relativeBase, dp.relativePath)
}

func newDetailPathList(ln int, cap int) detailPathList {
	return make(detailPathList, ln, cap)
}

type detailPathList []detailPath

func (l detailPathList) Len() int {
	return len(l)
}

func (l detailPathList) Less(i, j int) bool {
	if l[i].index != l[j].index {
		return l[i].index < l[j].index
	}
	lenI := getRuneCount(l[i].relativePath, filex.UnixSeparator)
	lenJ := getRuneCount(l[j].relativePath, filex.UnixSeparator)
	if lenI != lenJ {
		return lenI > lenJ
	}
	return l[i].relativePath < l[j].relativePath
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
