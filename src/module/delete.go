package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"os"
)

func newDeleteExecutor() IModeExecutor {
	return &deleteExecutor{list: newPathList(0, 128)}
}

type deleteExecutor struct {
	target *infra.RuntimeTarget
	list   pathList

	logger  logx.ILogger
	recurse bool
}

func (e *deleteExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Delete", Mode: infra.ModeDeleteValue, Src: src,
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
	e.initArgs()
	e.initExecuteList()
	e.execList()
}

func (e *deleteExecutor) initArgs() {
	argsMark := e.target.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.recurse = argsMark.MatchArg(infra.ArgMarkRecurse)
}

func (e *deleteExecutor) initExecuteList() {
	for index, src := range e.target.SrcArr {
		path := filex.Combine(infra.RunningDir, src)
		if !filex.IsExist(path) {
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]: %s", index, src))
			continue
		}
		e.checkPath(path)
	}
	e.list.Sort()
}

func (e *deleteExecutor) execList() {
	if e.list.Len() == 0 {
		return
	}
	for _, dir := range e.list {
		e.logger.Infoln("[delete] Delete Path:", dir)
		os.RemoveAll(dir)
	}
}

func (e *deleteExecutor) checkPath(fullPath string) {
	match := e.checkName(fullPath)

	if match {
		return
	}
	if e.recurse && filex.IsFolder(fullPath) {
		subPaths, _ := filex.GetPathsInDir(fullPath, nil)
		if len(subPaths) == 0 {
			return
		}
		for _, dir := range subPaths {
			e.checkPath(dir)
		}
	}
}

func (e *deleteExecutor) checkName(fullPath string) bool {
	_, filename := filex.Split(fullPath)
	// 名称不匹配
	if !e.target.CheckNameFitting(filename) {
		return false
	}
	e.list = append(e.list, fullPath)
	return true
}
