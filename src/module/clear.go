package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"os"
)

func newClearExecutor() IModeExecutor {
	return &clearExecutor{list: newPathList(0, 128)}
}

type clearExecutor struct {
	target *infra.RuntimeTarget
	list   pathList

	logger  logx.ILogger
	recurse bool
}

func (e *clearExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *clearExecutor) ExecConfigTarget(cfgTarget infra.ConfigTarget) {
	runtimeTarget := infra.NewRuntimeTarget(cfgTarget)
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *clearExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
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

func (e *clearExecutor) initArgs() {
	argsMark := e.target.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.recurse = argsMark.MatchArg(infra.ArgMarkRecurse)
}

func (e *clearExecutor) initExecuteList() {
	for index, src := range e.target.SrcArr {
		path := filex.Combine(infra.RunningDir, src)
		if !filex.IsFolder(path) {
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]: %s", index, src))
			continue
		}
		e.checkPath(path)
	}
	e.list.Sort()
}

func (e *clearExecutor) execList() {
	if e.list.Len() == 0 {
		return
	}
	for _, dir := range e.list {
		e.logger.Infoln("[clear] Clear Folder:", dir)
		os.RemoveAll(dir)
	}
}

func (e *clearExecutor) checkPath(fullPath string) {
	isFile := e.checkDir(fullPath)
	if isFile {
		return
	}
	if e.recurse {
		dirPaths, _ := filex.GetPathsInDir(fullPath, func(subPath string, info os.FileInfo) bool {
			return info.IsDir()
		})
		if len(dirPaths) == 0 {
			return
		}
		for _, dir := range dirPaths {
			e.checkPath(dir)
		}
	}
}

func (e *clearExecutor) checkDir(fullDir string) (isFile bool) {
	// 非目录
	if !filex.IsFolder(fullDir) {
		return true
	}
	_, filename := filex.Split(fullDir)
	// 名称不匹配
	if !e.target.CheckNameFitting(filename) {
		return false
	}
	if e.recurse {
		size, _ := filex.GetFolderSize(fullDir)
		// 非空
		if 0 != size {
			return false
		}
	} else {
		dirPaths, _ := filex.GetPathsInDir(fullDir, nil)
		// 非空
		if len(dirPaths) != 0 {
			return false
		}
	}
	e.list = append(e.list, fullDir)
	return true
}
