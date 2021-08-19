package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"io/fs"
	"sort"
)

// interface

type iModuleExecutor interface {
	// 初始化Log
	initLogger(mark infra.ArgMark)
	// 初始化处理列表
	initExecuteList()
	// 处理列表中文件
	execList()
}

type IModuleExecutor interface {
	// 执行任务
	Exec(src, tar, include, exclude, args string, wildcardCase bool)
	// 执行任务
	ExecConfigTarget(config infra.ConfigTarget)
	// 执行任务
	ExecRuntimeTarget(target *infra.RuntimeTarget)
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

// register

type ModuleExecutorGenerator func() IModuleExecutor

var generators = make([]ModuleExecutorGenerator, 64)

func GetExecutor(mode infra.RuntimeMode) IModuleExecutor {
	gen := generators[mode]
	if nil == gen {
		return nil
	}
	return gen()
}

func RegisterExecutor(mode infra.RuntimeMode, generator ModuleExecutorGenerator) {
	generators[mode] = generator
}

func init() {
	RegisterExecutor(infra.ModeClear, newClearExecutor)
	RegisterExecutor(infra.ModeCopy, newCopyExecutor)
	RegisterExecutor(infra.ModeDelete, newDeleteExecutor)
	RegisterExecutor(infra.ModeMove, newMoveExecutor)
	RegisterExecutor(infra.ModeSync, newSyncExecutor)
}
