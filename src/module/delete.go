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

func (e *deleteExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTarget{Name: "Delete", Mode: infra.ModeDeleteValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTarget(config)
}

func (e *deleteExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget, err := infra.NewRuntimeTarget(config)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[delete] Err : %v", err))
	}
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
	e.recurse = argsMark.MatchArg(infra.ArgRecurse)
}

func (e *deleteExecutor) initExecuteList() {
	for index, src := range e.target.SrcArr {
		path := filex.Combine(infra.RunningDir, src.FormattedSrc)
		if !filex.IsExist(path) {
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]: %s", index, src.OriginalSrc))
			continue
		}
		if src.IncludeSelf {
			e.checkPath(path, src)
		} else {
			e.checkSubPath(path, src)
		}
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

func (e *deleteExecutor) checkPath(fullPath string, srcInfo infra.SrcInfo) {
	if filex.IsFile(fullPath) { // 处理文件
		e.checkFileName(fullPath, srcInfo)
		return
	}
	if !e.recurse { // 非递归
		return
	}
	_, filename := filex.Split(fullPath)
	if !e.target.CheckDirFitting(filename) { // 过滤不匹配目录
		return
	}
	e.checkSubPath(fullPath, srcInfo)
}

func (e *deleteExecutor) checkSubPath(fullPath string, srcInfo infra.SrcInfo) {
	subPaths, _ := filex.GetPathsInDir(fullPath, nil)
	if len(subPaths) == 0 {
		return
	}
	for _, dir := range subPaths {
		e.checkPath(dir, srcInfo)
	}
}

func (e *deleteExecutor) checkFileName(fullPath string, srcInfo infra.SrcInfo) {
	_, filename := filex.Split(fullPath)
	// 路径通配符不匹配
	if !srcInfo.CheckFitting(filename) {
		return
	}
	// 名称不匹配
	if !e.target.CheckFileFitting(filename) {
		return
	}
	e.list = append(e.list, fullPath)
	return
}
