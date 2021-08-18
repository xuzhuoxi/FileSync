package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/logx"
)

func newCopyExecutor() IModuleExecutor {
	return &copyExecutor{}
}

type copyExecutor struct {
	target *infra.RuntimeTarget
	logger logx.ILogger
	list   infra.PathList
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
	if nil == target {
		return
	}
	if len(target.SrcArr) == 0 {
		return
	}
	if target.Tar == "" {
		return
	}
	e.target = target
	e.initLogger(target.ArgMarks)
	e.initExecuteList()
	e.execList()
}

func (e *copyExecutor) initLogger(mark infra.ArgMark) {
	e.logger = infra.GenLogger(mark)
}

func (e *copyExecutor) initExecuteList() {
	panic("implement me")
}

func (e *copyExecutor) execList() {
	panic("implement me")
}
