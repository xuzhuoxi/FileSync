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
	task *infra.RuntimeTask

	logger     logx.ILogger
	ignore     bool // 忽略空目录，查找文件时使用
	recurse    bool // 递归，查找文件时使用
	stable     bool // 保持目录结构，处理文件时使用
	file2file  bool // 单文件处理模式
	timeUpdate bool // 只处理新时间文件，处理文件时使用
	sizeUpdate bool // 只处理size更大文件，处理文件时使用
	md5Update  bool // 只处理md5不同的文件，处理文件时使用

	searcher    internal.IPathSearcher
	tempSrcInfo infra.SrcInfo
}

func (e *copyExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTask{Name: "Copy", Mode: infra.ModeCopyValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTask(config)
}

func (e *copyExecutor) ExecConfigTask(config infra.ConfigTask) {
	runtimeTask, err := infra.NewRuntimeTask(config)
	if nil != err {
		e.logger.Errorln(fmt.Sprintf("[copy] Err : %v", err))
	}
	e.ExecRuntimeTask(runtimeTask)
}

func (e *copyExecutor) ExecRuntimeTask(task *infra.RuntimeTask) {
	if nil == task {
		return
	}
	if len(task.SrcArr) == 0 || task.Tar == "" || strings.TrimSpace(task.Tar) == "" {
		return
	}
	e.task = task
	err := e.initArgs()
	if nil != err {
		// 由于logger可能初始化失败
		infra.Logger.Errorln(fmt.Sprintf("[copy] Init args error='%s'", err))
		return
	}
	e.initExecuteList()
	e.execList()
}

func (e *copyExecutor) initArgs() error {
	argsMark := e.task.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.ignore = argsMark.MatchArg(infra.MarkIgnoreEmpty)
	e.recurse = argsMark.MatchArg(infra.MarkRecurse)
	e.stable = argsMark.MatchArg(infra.MarkStable)
	e.file2file = argsMark.MatchArg(infra.MarkFile)
	e.timeUpdate = argsMark.MatchArg(infra.MarkTimeUpdate)
	e.sizeUpdate = argsMark.MatchArg(infra.MarkSizeUpdate)
	e.md5Update = argsMark.MatchArg(infra.MarkMd5Update)

	e.searcher.SetParams(e.recurse, !e.ignore, e.logger)
	return nil
}

func (e *copyExecutor) initExecuteList() {
	if e.file2file {
		e.initExecuteListFile2File()
	} else {
		e.initExecuteListCommon()
	}
}

func (e *copyExecutor) initExecuteListCommon() {
	e.searcher.SetFitting(e.fileFitting, e.dirFitting)
	e.searcher.InitSearcher()
	for _, src := range e.task.SrcArr {
		e.tempSrcInfo = src
		e.searcher.Search(src.FormattedSrc, src.IncludeSelf)
	}
	e.searcher.SortResults()
}

func (e *copyExecutor) initExecuteListFile2File() {
	if 1 != len(e.task.SrcArr) {
		e.logger.Warnln(fmt.Sprintf("[copy] Warn with src, len should be 1. "))
	}
	relateFilePath := e.task.SrcArr[0].FormattedSrc
	fileFull := filex.Combine(infra.RunningDir, relateFilePath)
	fileInfo, err := os.Stat(fileFull)
	if err != nil && !os.IsExist(err) { //不存在
		e.logger.Errorln(fmt.Sprintf("[copy] Error with src. It is not exist! [%s]", relateFilePath))
		return
	}
	e.searcher.InitSearcher()
	e.searcher.AppendResult(relateFilePath, fileInfo)
}

func (e *copyExecutor) execList() {
	e.logger.Infoln(fmt.Sprintf("[copy] Start(RunningPath='%s', Len=%d).", infra.RunningDir, e.searcher.ResultLen()))
	if e.file2file {
		e.execFile2File()
	} else {
		e.execListCommon()
	}
}

func (e *copyExecutor) execListCommon() {
	count := 0
	for _, srcPathInfo := range e.searcher.GetResults() {
		_, tarFull := internal.GetTarPaths(srcPathInfo, e.stable, e.task.Tar)
		if e.checkIgnored(srcPathInfo, tarFull) {
			continue
		}
		e.doCopy(srcPathInfo)
		count += 1
	}
	ignoreLen := e.searcher.ResultLen() - count
	e.logger.Infoln(fmt.Sprintf("[copy] Finish(CopyLen=%d, IgnoreLen=%d).", count, ignoreLen))
}

func (e *copyExecutor) execFile2File() {
	results := e.searcher.GetResults()
	if len(results) != 1 {
		e.logger.Warnln(fmt.Sprintf("[copy] Warn with result, len should be 1. "))
		return
	}
	src := results[0]
	tarFull := filex.Combine(infra.RunningDir, e.task.Tar)
	if e.checkIgnored(src, tarFull) {
		return
	}
	e.logger.Infoln(fmt.Sprintf("[copy] Copy '%s' => '%s'", src.GetRelativePath(), tarFull))
	internal.DoCopy(src.GetFullPath(), tarFull, nil)
}

func (e *copyExecutor) checkIgnored(src internal.IPathInfo, tarFull string) bool {
	tarFileInfo := infra.GetFileInfo(tarFull)
	if nil != tarFileInfo {
		fileInfo := src.GetFileInfo()
		if e.timeUpdate && infra.CompareWithTime(fileInfo, tarFileInfo) <= 0 { // 忽略目标新文件
			e.logger.Infoln(fmt.Sprintf("[copy] Ignored by '%s': '%s'", infra.ArgTimeUpdate, src.GetRelativePath()))
			return true
		}
		if e.sizeUpdate && infra.CompareWithSize(fileInfo, tarFileInfo) <= 0 { // 忽略目标大文件
			e.logger.Infoln(fmt.Sprintf("[copy] Ignored by '%s':'%s'", infra.ArgSizeUpdate, src.GetRelativePath()))
			return true
		}
		if e.md5Update && infra.CompareWithMd5(src.GetFullPath(), tarFull) { // 忽略md5相同文件
			e.logger.Infoln(fmt.Sprintf("[copy] Ignored by '%s':'%s'", infra.ArgMd5Update, src.GetRelativePath()))
			return true
		}
	}
	return false
}

func (e *copyExecutor) doCopy(pathInfo internal.IPathInfo) {
	tarRelative, tarFull := internal.GetTarPaths(pathInfo, e.stable, e.task.Tar)
	e.logger.Infoln(fmt.Sprintf("[copy] Copy '%s' => '%s'", pathInfo.GetRelativePath(), tarRelative))
	internal.DoCopy(pathInfo.GetFullPath(), tarFull, nil)
}

func (e *copyExecutor) fileFitting(fileInfo os.FileInfo) bool {
	if nil == fileInfo {
		return false
	}
	if !e.tempSrcInfo.CheckFitting(fileInfo.Name()) { // 路径匹配不成功
		return false
	}
	// 名称不匹配
	if !e.task.CheckFileFitting(fileInfo.Name()) {
		return false
	}
	return true
}

func (e *copyExecutor) dirFitting(dirInfo os.FileInfo) bool {
	if nil == dirInfo {
		return false
	}
	if !e.task.CheckDirFitting(dirInfo.Name()) {
		return false
	}
	return true
}
