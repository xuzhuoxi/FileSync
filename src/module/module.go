package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"os"
)

// interface

type iModuleExecutor interface {
	// 初始化Log
	initArgs() error
	// 初始化处理列表
	initExecuteList()
	// 处理列表中文件
	execList()

	// 文件过滤
	fileFitting(fileInfo os.FileInfo) bool
	// 目录过滤
	dirFitting(dirInfo os.FileInfo) bool
}

type IModeExecutor interface {
	// 执行任务
	Exec(src, tar, include, exclude, args string)
	// 执行任务
	ExecConfigTask(config infra.ConfigTask)
	// 执行任务
	ExecRuntimeTask(task *infra.RuntimeTask)
}

// register

type ModeExecutorGenerator func() IModeExecutor

var (
	generators = make([]ModeExecutorGenerator, 64)
)

func GetExecutor(mode infra.RuntimeMode) IModeExecutor {
	gen := generators[mode]
	if nil == gen {
		return nil
	}
	return gen()
}

func RegisterExecutor(mode infra.RuntimeMode, generator ModeExecutorGenerator) {
	generators[mode] = generator
}

func init() {
	RegisterExecutor(infra.ModeClear, newClearExecutor)
	RegisterExecutor(infra.ModeCopy, newCopyExecutor)
	RegisterExecutor(infra.ModeDelete, newDeleteExecutor)
	RegisterExecutor(infra.ModeMove, newMoveExecutor)
	RegisterExecutor(infra.ModeSync, newSyncExecutor)
}
