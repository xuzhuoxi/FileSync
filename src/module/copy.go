package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module/internal"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"os"
	"strings"
)

func newCopyExecutor() IModeExecutor {
	return &copyExecutor{searcher: internal.NewPathSearcher()}
}

type copyExecutor struct {
	target *infra.RuntimeTarget

	logger     logx.ILogger
	ignore     bool // 忽略空目录，查找文件时使用
	recurse    bool // 递归，查找文件时使用
	stable     bool // 保持目录结构，处理文件时使用
	timeUpdate bool // 只处理新时间文件，处理文件时使用
	sizeUpdate bool // 只处理size更大文件，处理文件时使用

	searcher    internal.IPathSearcher
	tempSrcInfo infra.SrcInfo
}

func (e *copyExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTarget{Name: "Copy", Mode: infra.ModeCopyValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTarget(config)
}

func (e *copyExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget, err := infra.NewRuntimeTarget(config)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[copy] Err : %v", err))
	}
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *copyExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	if nil == target {
		return
	}
	if len(target.SrcArr) == 0 || target.Tar == "" || strings.TrimSpace(target.Tar) == "" {
		return
	}
	e.target = target
	err := e.initArgs()
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[copy] Init args error='%s'", err))
		return
	}
	e.initExecuteList()
	e.execList()
}

func (e *copyExecutor) initArgs() error {
	argsMark := e.target.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.ignore = argsMark.MatchArg(infra.MarkIgnoreEmpty)
	e.recurse = argsMark.MatchArg(infra.MarkRecurse)
	e.stable = argsMark.MatchArg(infra.MarkStable)
	e.timeUpdate = argsMark.MatchArg(infra.MarkTimeUpdate)
	e.sizeUpdate = argsMark.MatchArg(infra.MarkSizeUpdate)

	e.searcher.SetParams(e.recurse, !e.ignore, e.logger)
	return nil
}

func (e *copyExecutor) initExecuteList() {
	e.searcher.SetFitting(e.fileFitting, e.dirFitting)
	e.searcher.InitSearcher()
	for _, src := range e.target.SrcArr {
		e.tempSrcInfo = src
		e.searcher.Search(src.FormattedSrc, src.IncludeSelf)
	}
	e.searcher.SortResults()
}

func (e *copyExecutor) execList() {
	e.logger.Infoln(fmt.Sprintf("[copy] Start(RunningPath='%s', Len=%d).", infra.RunningDir, e.searcher.ResultLen()))
	count := 0
	for _, srcPathInfo := range e.searcher.GetResults() {
		srcFileInfo := srcPathInfo.GetFileInfo()
		_, tarFull := internal.GetTarPaths(srcPathInfo, e.stable, e.target.Tar)
		tarFileInfo := infra.GetFileInfo(tarFull)
		if nil != tarFileInfo {
			if e.timeUpdate && infra.CompareWithTime(srcFileInfo, tarFileInfo) <= 0 { // 忽略目标新文件
				e.logger.Infoln(fmt.Sprintf("[copy] Ignored by '%s': '%s'", infra.ArgTimeUpdate, srcPathInfo.GetRelativePath()))
				continue
			}
			if e.sizeUpdate && infra.CompareWithSize(srcFileInfo, tarFileInfo) <= 0 { // 忽略目标大文件
				e.logger.Infoln(fmt.Sprintf("[move] Ignored by '%s':'%s'", infra.ArgSizeUpdate, srcPathInfo.GetRelativePath()))
				continue
			}
		}
		e.doCopy(srcPathInfo)
		count += 1
	}
	ignoreLen := e.searcher.ResultLen() - count
	e.logger.Infoln(fmt.Sprintf("[copy] Finish(CopyLen=%d, IgnoreLen=%d).", count, ignoreLen))
}

func (e *copyExecutor) doCopy(pathInfo internal.IPathInfo) {
	tarRelative, tarFull := internal.GetTarPaths(pathInfo, e.stable, e.target.Tar)
	fileInfo := pathInfo.GetFileInfo()
	e.logger.Infoln(fmt.Sprintf("[copy] Copy '%s' => '%s'", pathInfo.GetRelativePath(), tarRelative))
	if fileInfo.IsDir() {
		os.MkdirAll(tarFull, fileInfo.Mode())
	} else {
		filex.CopyAuto(pathInfo.GetFullPath(), tarFull, fileInfo.Mode())
	}
	infra.SetModTime(tarFull, fileInfo.ModTime())
}

func (e *copyExecutor) fileFitting(fileInfo os.FileInfo) bool {
	if nil == fileInfo {
		return false
	}
	if !e.tempSrcInfo.CheckFitting(fileInfo.Name()) { // 路径匹配不成功
		return false
	}
	// 名称不匹配
	if !e.target.CheckFileFitting(fileInfo.Name()) {
		return false
	}
	return true
}

func (e *copyExecutor) dirFitting(dirInfo os.FileInfo) bool {
	if nil == dirInfo {
		return false
	}
	if !e.target.CheckDirFitting(dirInfo.Name()) {
		return false
	}
	return true
}
