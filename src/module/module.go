package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
)

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
