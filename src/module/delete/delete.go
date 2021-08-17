package delete

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module"
)

type deleteExecutor struct {
}

func (e *deleteExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *deleteExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget := infra.NewRuntimeTarget(config)
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *deleteExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	infra.Logger.Info("Delete", target)
}

func newExecutor() module.IModuleExecutor {
	return &deleteExecutor{}
}

func init() {
	module.RegisterExecutor(infra.ModeDelete, newExecutor)
}
