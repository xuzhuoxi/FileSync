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

	logger  logx.ILogger
	double  bool
	ignore  bool
	recurse bool
	update  bool
}

func (e *syncExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTarget{Name: "Sync", Mode: infra.ModeSyncValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTarget(config)
}

func (e *syncExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget, err := infra.NewRuntimeTarget(config)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Err : %v", err))
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
	e.double = argsMark.MatchArg(infra.ArgDouble)
	e.ignore = argsMark.MatchArg(infra.ArgIgnoreEmpty)
	e.recurse = argsMark.MatchArg(infra.ArgRecurse)
	e.update = argsMark.MatchArg(infra.ArgUpdate)
}

func (e *syncExecutor) initExecuteList() {
	panic("implement me")
}

func (e *syncExecutor) execList() {
	panic("implement me")
}
