package sync

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module"
)

type syncExecutor struct {
}

func (e *syncExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *syncExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget := infra.NewRuntimeTarget(config)
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *syncExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	infra.Logger.Info("Sync", target)
}

func newExecutor() module.IModuleExecutor {
	return &syncExecutor{}
}

func init() {
	module.RegisterExecutor(infra.ModeSync, newExecutor)
}
