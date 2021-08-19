package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"os"
)

func newDeleteExecutor() IModuleExecutor {
	return &deleteExecutor{list: newPathList(0, 128)}
}

type deleteExecutor struct {
	target *infra.RuntimeTarget
	logger logx.ILogger
	list   pathList
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
	e.initLogger(target.ArgsMark)
	e.initExecuteList()
	e.execList()
}

func (e *deleteExecutor) initLogger(mark infra.ArgMark) {
	e.logger = infra.GenLogger(mark)
}

func (e *deleteExecutor) initExecuteList() {
	recurse := e.target.ArgsMark.MatchArg(infra.ArgMarkRecurse)
	for index, src := range e.target.SrcArr {
		path := filex.Combine(infra.RunningDir, src)
		if !filex.IsExist(path) {
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]: %s", index, src))
			continue
		}
		e.checkPath(path, recurse)
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

func (e *deleteExecutor) checkPath(fullPath string, recurse bool) {
	match := e.checkName(fullPath)

	if match {
		return
	}
	if recurse && filex.IsFolder(fullPath) {
		subPaths, _ := filex.GetPathsInDir(fullPath, nil)
		if len(subPaths) == 0 {
			return
		}
		for _, dir := range subPaths {
			e.checkPath(dir, true)
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
