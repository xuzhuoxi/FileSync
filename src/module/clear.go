package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"os"
)

func newClearExecutor() IModuleExecutor {
	return &clearExecutor{list: newPathList(0, 128)}
}

type clearExecutor struct {
	target *infra.RuntimeTarget
	logger logx.ILogger
	list   pathList
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
	e.initLogger(target.ArgsMark)
	e.initExecuteList()
	e.execList()
}

func (e *clearExecutor) initLogger(mark infra.ArgMark) {
	e.logger = infra.GenLogger(mark)
}

func (e *clearExecutor) initExecuteList() {
	recurse := e.target.ArgsMark.MatchArg(infra.ArgMarkRecurse)
	for index, src := range e.target.SrcArr {
		path := filex.Combine(infra.RunningDir, src)
		if !filex.IsFolder(path) {
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]: %s", index, src))
			continue
		}
		e.checkPath(path, recurse)
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

func (e *clearExecutor) checkPath(fullPath string, recurse bool) {
	isFile := e.checkDir(fullPath, recurse)
	if isFile {
		return
	}
	if recurse {
		dirPaths, _ := filex.GetPathsInDir(fullPath, func(subPath string, info os.FileInfo) bool {
			return info.IsDir()
		})
		if len(dirPaths) == 0 {
			return
		}
		for _, dir := range dirPaths {
			e.checkPath(dir, true)
		}
	}
}

func (e *clearExecutor) checkDir(fullDir string, recurse bool) (isFile bool) {
	// 非目录
	if !filex.IsFolder(fullDir) {
		return true
	}
	_, filename := filex.Split(fullDir)
	// 名称不匹配
	if !e.target.CheckNameFitting(filename) {
		return false
	}
	if recurse {
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
