package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
)

// interface

type IModeExecutor interface {
	// 执行任务
	Exec(src, tar, include, exclude, args string)
	// 执行任务
	ExecConfigTarget(config infra.ConfigTarget)
	// 执行任务
	ExecRuntimeTarget(target *infra.RuntimeTarget)
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
