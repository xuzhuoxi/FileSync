package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module/internal"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"os"
)

func newClearExecutor() IModeExecutor {
	return &clearExecutor{list: internal.NewPathStrList(0, 128)}
}

type clearExecutor struct {
	target *infra.RuntimeTarget
	list   internal.IPathStrList

	logger  logx.ILogger
	recurse bool // 递归，查找文件时使用

	tempSrcInfo infra.SrcInfo
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
	e.recurse = argsMark.MatchArg(infra.MarkRecurse)
}

func (e *clearExecutor) initExecuteList() {
	for index, src := range e.target.SrcArr {
		e.tempSrcInfo = src
		path := filex.Combine(infra.RunningDir, src.FormattedSrc)
		if !filex.IsFolder(path) {
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]: %s", index, src.OriginalSrc))
			continue
		}
		if src.IncludeSelf {
			e.checkDir(path)
		} else {
			e.checkSubDir(path)
		}
	}
	e.list.Sort()
}

func (e *clearExecutor) execList() {
	if e.list.Len() == 0 {
		return
	}
	for _, dir := range e.list.GetAll() {
		e.logger.Infoln(fmt.Sprintf("[clear] Clear Folder='%s'", dir))
		os.RemoveAll(dir)
	}
}

func (e *clearExecutor) fileFitting(fileInfo os.FileInfo) bool {
	return false
}

func (e *clearExecutor) dirFitting(dirInfo os.FileInfo) bool {
	// 路径通配符不匹配
	filename := dirInfo.Name()
	if !e.tempSrcInfo.CheckFitting(filename) {
		return false
	}
	// 名称不匹配
	if !e.target.CheckDirFitting(filename) {
		return false
	}
	return true
}

func (e *clearExecutor) checkPath(fullPath string) {
	isFile := e.checkDir(fullPath)
	if isFile {
		return
	}
	if !e.recurse {
		return
	}
	e.checkSubDir(fullPath)
}

func (e *clearExecutor) checkSubDir(fullDir string) {
	dirPaths, _ := filex.GetPathsInDir(fullDir, func(subPath string, info os.FileInfo) bool {
		return info.IsDir()
	})
	if len(dirPaths) == 0 {
		return
	}
	for _, dir := range dirPaths {
		e.checkPath(dir)
	}
}

func (e *clearExecutor) checkDir(fullDir string) (Interrupt bool) {
	fileInfo, err := os.Stat(fullDir)
	if err != nil && !os.IsExist(err) { //不存在
		return true
	}
	if !fileInfo.IsDir() {
		return true
	}
	if !e.dirFitting(fileInfo) {
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
	e.list.Append(fullDir)
	return true
}
