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

func (e *clearExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTarget(config)
}

func (e *clearExecutor) ExecConfigTarget(cfgTarget infra.ConfigTarget) {
	runtimeTarget, err := infra.NewRuntimeTarget(cfgTarget)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[clear] Err : %v", err))
	}
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
	e.recurse = argsMark.MatchArg(infra.ArgRecurse)
}

func (e *clearExecutor) initExecuteList() {
	for index, src := range e.target.SrcArr {
		path := filex.Combine(infra.RunningDir, src.FormattedSrc)
		if !filex.IsFolder(path) {
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]: %s", index, src.OriginalSrc))
			continue
		}
		if src.IncludeSelf {
			e.checkPath(path, src)
		} else {
			e.CheckSubDir(path, src)
		}
	}
	e.list.Sort()
}

func (e *clearExecutor) execList() {
	if e.list.Len() == 0 {
		return
	}
	for _, dir := range e.list {
		e.logger.Infoln(fmt.Sprintf("[clear] Clear Folder='%s'", dir))
		os.RemoveAll(dir)
	}
}

func (e *clearExecutor) checkPath(fullPath string, srcInfo infra.SrcInfo) {
	isFile := e.checkDir(fullPath, srcInfo)
	if isFile {
		return
	}
	if !e.recurse {
		return
	}
	e.CheckSubDir(fullPath, srcInfo)
}

func (e *clearExecutor) CheckSubDir(fullDir string, srcInfo infra.SrcInfo) {
	dirPaths, _ := filex.GetPathsInDir(fullDir, func(subPath string, info os.FileInfo) bool {
		return info.IsDir()
	})
	if len(dirPaths) == 0 {
		return
	}
	for _, dir := range dirPaths {
		e.checkPath(dir, srcInfo)
	}
}

func (e *clearExecutor) checkDir(fullDir string, srcInfo infra.SrcInfo) (isFile bool) {
	// 非目录
	if !filex.IsFolder(fullDir) {
		return true
	}
	_, filename := filex.Split(fullDir)
	// 路径通配符不匹配
	if !srcInfo.CheckFitting(filename) {
		return false
	}
	// 名称不匹配
	if !e.target.CheckDirFitting(filename) {
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
