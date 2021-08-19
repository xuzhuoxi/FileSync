package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/logx"
)

func newMoveExecutor() IModuleExecutor {
	return &moveExecutor{}
}

type moveExecutor struct {
	target *infra.RuntimeTarget
	logger logx.ILogger
	list   pathList
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
	infra.Logger.Info("Move", target)
}

func (e *moveExecutor) initLogger(mark infra.ArgMark) {
	e.logger = infra.GenLogger(mark)
}

func (e *moveExecutor) initExecuteList() {
	panic("implement me")
}

func (e *moveExecutor) execList() {
	panic("implement me")
}
