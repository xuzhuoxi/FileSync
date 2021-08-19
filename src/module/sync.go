package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/logx"
)

func newSyncExecutor() IModuleExecutor {
	return &syncExecutor{}
}

type syncExecutor struct {
	target *infra.RuntimeTarget
	logger logx.ILogger
	list   pathList
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

func (e *syncExecutor) initLogger(mark infra.ArgMark) {
	e.logger = infra.GenLogger(mark)
}

func (e *syncExecutor) initExecuteList() {
	panic("implement me")
}

func (e *syncExecutor) execList() {
	panic("implement me")
}
