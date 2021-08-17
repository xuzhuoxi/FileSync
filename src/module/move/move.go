package move

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module"
)

type moveExecutor struct {
}

func (e *moveExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *moveExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget := infra.NewRuntimeTarget(config)
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *moveExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	panic("implement me")
}

func newMoveExecutor() module.IModuleExecutor {
	return &moveExecutor{}
}

func init() {
	module.RegisterExecutor(infra.ModeMove, newMoveExecutor)
}
