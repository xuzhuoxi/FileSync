package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/logx"
)

func newDeleteExecutor() IModuleExecutor {
	return &deleteExecutor{}
}

type deleteExecutor struct {
	target *infra.RuntimeTarget
	logger logx.ILogger
	list   infra.PathList
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
	if nil == target {
		return
	}
	if len(target.SrcArr) == 0 {
		return
	}
	e.target = target
	e.initLogger(target.ArgMarks)
	e.initExecuteList()
	e.execList()
}

func (e *deleteExecutor) initLogger(mark infra.ArgMark) {
	e.logger = infra.GenLogger(mark)
}

func (e *deleteExecutor) initExecuteList() {
	panic("implement me")
}

func (e *deleteExecutor) execList() {
	panic("implement me")
}
