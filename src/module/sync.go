package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/logx"
	"strings"
)

func newSyncExecutor() IModeExecutor {
	return &syncExecutor{}
}

type syncExecutor struct {
	target *infra.RuntimeTarget
	list   pathList

	logger  logx.ILogger
	double  bool
	ignore  bool
	recurse bool
	update  bool
}

func (e *syncExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Sync", Mode: infra.ModeSyncValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *syncExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget, err := infra.NewRuntimeTarget(config)
	if nil != err {
		e.logger.Errorln(fmt.Sprintf("[sync] Err : %v", err))
	}
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *syncExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	if nil == target {
		return
	}
	if len(target.SrcArr) != 1 || target.Tar == "" || strings.TrimSpace(target.Tar) == "" {
		return
	}
	e.target = target
	e.initArgs()
	e.initExecuteList()
	e.execList()
}

func (e *syncExecutor) initArgs() {
	argsMark := e.target.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.double = argsMark.MatchArg(infra.ArgMarkDouble)
	e.ignore = argsMark.MatchArg(infra.ArgMarkIgnore)
	e.recurse = argsMark.MatchArg(infra.ArgMarkRecurse)
	e.update = argsMark.MatchArg(infra.ArgMarkUpdate)
}

func (e *syncExecutor) initExecuteList() {
	panic("implement me")
}

func (e *syncExecutor) execList() {
	panic("implement me")
}
