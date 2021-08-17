package copy

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module"
)

type copyExecutor struct {
}

func (e *copyExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *copyExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget := infra.NewRuntimeTarget(config)
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *copyExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	infra.Logger.Info("Copy", target)
}

func newExecutor() module.IModuleExecutor {
	return &copyExecutor{}
}

func init() {
	module.RegisterExecutor(infra.ModeCopy, newExecutor)
}
