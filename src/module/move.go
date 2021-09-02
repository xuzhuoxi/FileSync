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

func (e *moveExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTarget{Name: "Move", Mode: infra.ModeMoveValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTarget(config)
}

func (e *moveExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget, err := infra.NewRuntimeTarget(config)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[move] Err : %v", err))
	}
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *moveExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	if nil == target {
		return
	}
	if len(target.SrcArr) == 0 || target.Tar == "" || strings.TrimSpace(target.Tar) == "" {
		return
	}
	e.target = target
	err := e.initArgs()
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[move] Init args error='%s'", err))
		return
	}
	e.initExecuteList()
	e.execList()
}

func (e *moveExecutor) initArgs() error {
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

func (e *moveExecutor) initExecuteList() {
	e.searcher.SetFitting(e.fileFitting, e.dirFitting)
	e.searcher.InitSearcher()
	for _, src := range e.target.SrcArr {
		e.tempSrcInfo = src
		e.searcher.Search(src.FormattedSrc, src.IncludeSelf)
	}
	e.searcher.SortResults()
}

func (e *moveExecutor) execList() {
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
		_, tarFull = internal.GetTarPaths(srcPathInfo, e.stable, e.target.Tar)
		tarFileInfo := infra.GetFileInfo(tarFull)
		if nil != tarFileInfo {
			if e.timeUpdate && infra.CompareWithTime(srcFileInfo, tarFileInfo) <= 0 { // 忽略目标新文件
				e.logger.Infoln(fmt.Sprintf("[move] Ignored by '%s':'%s'", infra.ArgTimeUpdate, srcPathInfo.GetRelativePath()))
				continue
			}
			if e.sizeUpdate && infra.CompareWithSize(srcFileInfo, tarFileInfo) <= 0 { // 忽略目标大文件
				e.logger.Infoln(fmt.Sprintf("[move] Ignored by '%s':'%s'", infra.ArgSizeUpdate, srcPathInfo.GetRelativePath()))
				continue
			}
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

func (e *moveExecutor) doMoveFile(pathInfo internal.IPathInfo) {
	tarRelative, tarFull := internal.GetTarPaths(pathInfo, e.stable, e.target.Tar)
	e.logger.Infoln(fmt.Sprintf("[move] Move file '%s' => '%s'", pathInfo.GetRelativePath(), tarRelative))

	filex.CompleteParentPath(tarFull, pathInfo.GetFileInfo().Mode())
	os.Rename(pathInfo.GetFullPath(), tarFull)
}

func (e *moveExecutor) doMoveDir(pathInfo internal.IPathInfo) {
	tarRelative, tarFull := internal.GetTarPaths(pathInfo, e.stable, e.target.Tar)
	e.logger.Infoln(fmt.Sprintf("[move] Move Dir '%s' => '%s'", pathInfo.GetRelativePath(), tarRelative))

	fileInfo := pathInfo.GetFileInfo()
	if filex.IsDir(tarFull) { // 目录存在
		infra.SetModTime(tarFull, fileInfo.ModTime())
		filex.Remove(pathInfo.GetFullPath())
		return
	}
	filex.CompleteParentPath(tarFull, fileInfo.Mode())
	os.Rename(pathInfo.GetFullPath(), tarFull)
}

func (e *moveExecutor) fileFitting(fileInfo os.FileInfo) bool {
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

func (e *moveExecutor) dirFitting(dirInfo os.FileInfo) bool {
	if nil == dirInfo {
		return false
	}
	if !e.target.CheckDirFitting(dirInfo.Name()) {
		return false
	}
	return true
}
