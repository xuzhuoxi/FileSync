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

	logger  logx.ILogger
	ignore  bool // 处理复制列表时使用
	recurse bool // 处理复制列表时使用
	stable  bool // 处理复制列表时使用
	update  bool // 真实复制时使用

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
	e.initArgs()
	e.initExecuteList()
	e.execList()
}

func (e *moveExecutor) initArgs() {
	argsMark := e.target.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.ignore = argsMark.MatchArg(infra.ArgIgnoreEmpty)
	e.recurse = argsMark.MatchArg(infra.ArgRecurse)
	e.stable = argsMark.MatchArg(infra.ArgStable)
	e.update = argsMark.MatchArg(infra.ArgUpdate)

	e.searcher.SetParams(e.recurse, !e.ignore, e.logger)
}

func (e *moveExecutor) initExecuteList() {
	e.searcher.SetFitting(e.fileFitting, e.dirFitting)
	e.searcher.InitSearcher()
	for _, src := range e.target.SrcArr {
		e.tempSrcInfo = src
		e.searcher.Search(src.FormattedSrc, src.IncludeSelf)
	}
	e.searcher.SortResult()
}

func (e *moveExecutor) execList() {
	results := e.searcher.GetResults()
	resultLen := len(results)
	e.logger.Infoln(fmt.Sprintf("[move] Start(RunningPath='%s', Len=%d).", infra.RunningDir, resultLen))
	count := 0
	var fileInfo os.FileInfo
	var tarFull string
	for _, moveInfo := range results {
		fileInfo = moveInfo.GetFileInfo()
		// 忽略目录
		if fileInfo.IsDir() {
			continue
		}
		_, tarFull = internal.GetTarPaths(moveInfo, e.stable, e.target.Tar)
		// 忽略新文件
		if e.update && !infra.CheckPathByTime(tarFull, fileInfo.ModTime()) {
			e.logger.Infoln(fmt.Sprintf("[move] Ignore by '/u': '%s'", moveInfo.GetRelativePath()))
			continue
		}
		e.doMoveFile(moveInfo)
		count += 1
	}
	for _, moveInfo := range results {
		fileInfo = moveInfo.GetFileInfo()
		// 忽略文件
		if !fileInfo.IsDir() {
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
		infra.CloneTime(tarFull, fileInfo)
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
