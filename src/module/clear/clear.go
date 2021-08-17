package clear

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module"
)

type clearExecutor struct {
}

func (e *clearExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *clearExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget := infra.NewRuntimeTarget(config)
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *clearExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	panic("implement me")
}

func newClearExecutor() module.IModuleExecutor {
	return &clearExecutor{}
}

func init() {
	module.RegisterExecutor(infra.ModeClear, newClearExecutor)
}
