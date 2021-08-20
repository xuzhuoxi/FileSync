package module

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/logx"
	"strings"
)

func newMoveExecutor() IModeExecutor {
	return &moveExecutor{}
}

type moveExecutor struct {
	target *infra.RuntimeTarget
	list   pathList

	logger  logx.ILogger
	force   bool
	ignore  bool
	recurse bool
	stable  bool
	update  bool
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
	if nil == target {
		return
	}
	if len(target.SrcArr) == 0 || target.Tar == "" || strings.TrimSpace(target.Tar) == "" {
		return
	}
	e.target = target
	e.initArgs()
	e.initExecuteList()
	e.execList()
}

func (e *moveExecutor) initArgs() {
	argsMark := e.target.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.force = argsMark.MatchArg(infra.ArgMarkForce)
	e.ignore = argsMark.MatchArg(infra.ArgMarkIgnore)
	e.recurse = argsMark.MatchArg(infra.ArgMarkRecurse)
	e.stable = argsMark.MatchArg(infra.ArgMarkStable)
	e.update = argsMark.MatchArg(infra.ArgMarkUpdate)
}

func (e *moveExecutor) initExecuteList() {
	panic("implement me")
}

func (e *moveExecutor) execList() {
	panic("implement me")
}
