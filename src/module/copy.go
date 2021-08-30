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

	logger  logx.ILogger
	ignore  bool // 处理复制列表时使用
	recurse bool // 处理复制列表时使用
	stable  bool // 处理复制列表时使用
	update  bool // 真实复制时使用

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
	e.initArgs()
	e.initExecuteList()
	e.execList()
}

func (e *copyExecutor) initArgs() {
	argsMark := e.target.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.ignore = argsMark.MatchArg(infra.ArgIgnoreEmpty)
	e.recurse = argsMark.MatchArg(infra.ArgRecurse)
	e.stable = argsMark.MatchArg(infra.ArgStable)
	e.update = argsMark.MatchArg(infra.ArgUpdate)

	e.searcher.SetParams(e.recurse, !e.ignore, e.logger)
}

func (e *copyExecutor) initExecuteList() {
	e.searcher.SetFitting(e.fileFitting, e.dirFitting)
	e.searcher.InitSearcher()
	for _, src := range e.target.SrcArr {
		e.tempSrcInfo = src
		e.searcher.Search(src.FormattedSrc, src.IncludeSelf)
	}
	e.searcher.SortResult()
}

func (e *copyExecutor) execList() {
	e.logger.Infoln(fmt.Sprintf("[copy] Start(RunningPath='%s', Len=%d).", infra.RunningDir, e.searcher.ResultLen()))
	count := 0
	for _, copyInfo := range e.searcher.GetResults() {
		fileInfo := copyInfo.GetFileInfo()
		// 忽略新文件
		if e.update && !fileInfo.IsDir() {
			_, tarFullPath := internal.GetTarPaths(copyInfo, e.stable, e.target.Tar)
			if !infra.CheckPathByTime(tarFullPath, fileInfo.ModTime()) {
				e.logger.Infoln(fmt.Sprintf("[copy] Ignore by '/u': '%s'", copyInfo.GetFullPath()))
				continue
			}
		}
		e.doCopy(copyInfo)
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
	infra.CloneTime(tarFull, fileInfo)
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
