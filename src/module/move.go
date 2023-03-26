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

func newMoveExecutor() IModeExecutor {
	return &moveExecutor{searcher: internal.NewPathSearcher()}
}

type moveExecutor struct {
	task *infra.RuntimeTask

	logger     logx.ILogger
	ignore     bool // 忽略空目录，查找文件时使用
	recurse    bool // 递归，查找文件时使用
	stable     bool // 保持目录结构，处理文件时使用
	file2file  bool // 单文件处理模式
	timeUpdate bool // 只处理新时间文件，处理文件时使用
	sizeUpdate bool // 只处理size更大文件，处理文件时使用
	md5Update  bool // 只处理md5不同文件，处理文件时使用

	searcher    internal.IPathSearcher
	tempSrcInfo infra.SrcInfo
}

func (e *moveExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTask{Name: "Move", Mode: infra.ModeMoveValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTask(config)
}

func (e *moveExecutor) ExecConfigTask(config infra.ConfigTask) {
	runtimeTask, err := infra.NewRuntimeTask(config)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[move] Err : %v", err))
	}
	e.ExecRuntimeTask(runtimeTask)
}

func (e *moveExecutor) ExecRuntimeTask(task *infra.RuntimeTask) {
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
		infra.Logger.Errorln(fmt.Sprintf("[move] Init args error='%s'", err))
		return
	}
	e.initExecuteList()
	e.execList()
}

func (e *moveExecutor) initArgs() error {
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

func (e *moveExecutor) initExecuteList() {
	if e.file2file {
		e.initExecuteListFile2File()
	} else {
		e.initExecuteListCommon()
	}
}

func (e *moveExecutor) initExecuteListCommon() {
	e.searcher.SetFitting(e.fileFitting, e.dirFitting)
	e.searcher.InitSearcher()
	for _, src := range e.task.SrcArr {
		e.tempSrcInfo = src
		e.searcher.Search(src.FormattedSrc, src.IncludeSelf)
	}
	e.searcher.SortResults()
}

func (e *moveExecutor) initExecuteListFile2File() {
	if 1 != len(e.task.SrcArr) {
		e.logger.Warnln(fmt.Sprintf("[move] Warn with src, len should be 1. "))
	}
	relateFilePath := e.task.SrcArr[0].FormattedSrc
	fileFull := filex.Combine(infra.RunningDir, relateFilePath)
	fileInfo, err := os.Stat(fileFull)
	if err != nil && !os.IsExist(err) { //不存在
		e.logger.Errorln(fmt.Sprintf("[move] Error with src. It is not exist! [%s]", relateFilePath))
		return
	}
	e.searcher.InitSearcher()
	e.searcher.AppendResult(relateFilePath, fileInfo)
}

func (e *moveExecutor) execList() {
	if e.file2file {
		e.execFile2File()
	} else {
		e.execListCommon()
	}
}
func (e *moveExecutor) execListCommon() {
	results := e.searcher.GetResults()
	resultLen := len(results)
	e.logger.Infoln(fmt.Sprintf("[move] Start(RunningPath='%s', Len=%d).", infra.RunningDir, resultLen))
	count := 0
	var srcFileInfo os.FileInfo
	var tarFull string
	for _, srcPathInfo := range results {
		srcFileInfo = srcPathInfo.GetFileInfo()
		// 忽略目录
		if srcFileInfo.IsDir() {
			continue
		}
		_, tarFull = internal.GetTarPaths(srcPathInfo, e.stable, e.task.Tar)
		if e.checkIgnored(srcPathInfo, tarFull) {
			continue
		}
		e.doMoveFile(srcPathInfo)
		count += 1
	}
	for _, moveInfo := range results {
		srcFileInfo = moveInfo.GetFileInfo()
		// 忽略文件
		if !srcFileInfo.IsDir() {
			continue
		}
		srcFull := moveInfo.GetFullPath()
		// 忽略非空目录
		if !filex.IsEmptyDir(srcFull) {
			continue
		}
		e.doMoveDir(moveInfo)
		count += 1
	}
	ignoreLen := resultLen - count
	e.logger.Infoln(fmt.Sprintf("[move] Finish(CopyLen=%d, IgnoreLen=%d).", count, ignoreLen))
}

func (e *moveExecutor) execFile2File() {
	results := e.searcher.GetResults()
	if len(results) != 1 {
		e.logger.Warnln(fmt.Sprintf("[move] Warn with result, len should be 1. "))
		return
	}
	src := results[0]
	tarFull := filex.Combine(infra.RunningDir, e.task.Tar)
	if e.checkIgnored(src, tarFull) {
		return
	}
	e.logger.Infoln(fmt.Sprintf("[move] Copy '%s' => '%s'", src.GetRelativePath(), tarFull))
	internal.DoMove(src.GetFullPath(), tarFull, nil)
}

func (e *moveExecutor) checkIgnored(src internal.IPathInfo, tarFull string) bool {
	tarFileInfo := infra.GetFileInfo(tarFull)
	if nil != tarFileInfo {
		fileInfo := src.GetFileInfo()
		if e.timeUpdate && infra.CompareWithTime(fileInfo, tarFileInfo) <= 0 { // 忽略目标新文件
			e.logger.Infoln(fmt.Sprintf("[move] Ignored by '%s': '%s'", infra.ArgTimeUpdate, src.GetRelativePath()))
			return true
		}
		if e.sizeUpdate && infra.CompareWithSize(fileInfo, tarFileInfo) <= 0 { // 忽略目标大文件
			e.logger.Infoln(fmt.Sprintf("[move] Ignored by '%s':'%s'", infra.ArgSizeUpdate, src.GetRelativePath()))
			return true
		}
		if e.md5Update && infra.CompareWithMd5(src.GetFullPath(), tarFull) { // 忽略md5相同文件
			e.logger.Infoln(fmt.Sprintf("[move] Ignored by '%s':'%s'", infra.ArgMd5Update, src.GetRelativePath()))
			return true
		}
	}
	return false
}

func (e *moveExecutor) doMoveFile(pathInfo internal.IPathInfo) {
	tarRelative, tarFull := internal.GetTarPaths(pathInfo, e.stable, e.task.Tar)
	e.logger.Infoln(fmt.Sprintf("[move] Move file2file '%s' => '%s'", pathInfo.GetRelativePath(), tarRelative))

	filex.CompleteParentPath(tarFull, pathInfo.GetFileInfo().Mode())
	internal.DoMove(pathInfo.GetFullPath(), tarFull, nil)
}

func (e *moveExecutor) doMoveDir(pathInfo internal.IPathInfo) {
	tarRelative, tarFull := internal.GetTarPaths(pathInfo, e.stable, e.task.Tar)
	e.logger.Infoln(fmt.Sprintf("[move] Move Dir '%s' => '%s'", pathInfo.GetRelativePath(), tarRelative))

	fileInfo := pathInfo.GetFileInfo()
	if filex.IsDir(tarFull) { // 目录存在
		infra.SetModTime(tarFull, fileInfo.ModTime())
		filex.Remove(pathInfo.GetFullPath())
		return
	}
	filex.CompleteParentPath(tarFull, fileInfo.Mode())
	internal.DoMove(pathInfo.GetFullPath(), tarFull, nil)
}

func (e *moveExecutor) fileFitting(fileInfo os.FileInfo) bool {
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

func (e *moveExecutor) dirFitting(dirInfo os.FileInfo) bool {
	if nil == dirInfo {
		return false
	}
	if !e.task.CheckDirFitting(dirInfo.Name()) {
		return false
	}
	return true
}
