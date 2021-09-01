package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module/internal"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/mathx"
	"os"
	"strings"
)

func newSyncExecutor() IModeExecutor {
	return &syncExecutor{srcList: internal.NewPathSearcher(), tarList: internal.NewPathSearcher()}
}

type syncExecutor struct {
	target *infra.RuntimeTarget

	logger  logx.ILogger
	double  bool // 双向同步
	ignore  bool // 忽略空目录，查找文件时使用
	recurse bool // 保持目录结构，处理文件时使用
	update  bool // 只处理新时间文件，处理文件时使用

	srcList internal.IPathSearcher
	tarList internal.IPathSearcher

	mixedArr  []internal.IPathInfo
	srcNewArr []internal.IPathInfo
	tarNewArr []internal.IPathInfo
}

func (e *syncExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTarget{Name: "Sync", Mode: infra.ModeSyncValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTarget(config)
}

func (e *syncExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget, err := infra.NewRuntimeTarget(config)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Err : %v", err))
	}
	e.ExecRuntimeTarget(runtimeTarget)
}

func (e *syncExecutor) ExecRuntimeTarget(target *infra.RuntimeTarget) {
	if nil == target {
		return
	}
	if len(target.SrcArr) != 1 {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Src Len Err! "))
		return
	}
	if target.SrcArr[0].FormattedSrc == "" || strings.TrimSpace(target.SrcArr[0].FormattedSrc) == "" {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Src Empty Err! "))
		return
	}
	if target.Tar == "" || strings.TrimSpace(target.Tar) == "" {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Tar Empty Err! "))
		return
	}
	e.target = target
	e.initArgs()
	e.initExecuteList()
	e.execList()
}

func (e *syncExecutor) initArgs() {
	argsMark := e.target.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.double = argsMark.MatchArg(infra.MarkDouble)
	e.ignore = argsMark.MatchArg(infra.MarkIgnoreEmpty)
	e.recurse = argsMark.MatchArg(infra.MarkRecurse)
	e.update = argsMark.MatchArg(infra.MarkTimeUpdate)

	e.srcList.SetParams(e.recurse, !e.ignore, e.logger)
	e.tarList.SetParams(e.recurse, !e.ignore, e.logger)
}

func (e *syncExecutor) initExecuteList() {
	// 查找源
	e.srcList.SetFitting(e.fileFitting, e.dirFitting)
	e.srcList.InitSearcher()
	src := e.target.SrcArr[0].FormattedSrc
	e.srcList.Search(src, false)
	e.srcList.SortResults()

	// 查找目标
	e.tarList.SetFitting(e.fileFitting, e.dirFitting)
	e.tarList.InitSearcher()
	e.tarList.Search(e.target.Tar, false)
	e.tarList.SortResults()

	// 源集合与目录集合的交集与差集计算
	e.OperateSets()
}

func (e *syncExecutor) execList() {
	e.logger.Infoln(fmt.Sprintf("Mixed Len=%d", len(e.mixedArr)))
	for _, m := range e.mixedArr {
		e.logger.Infoln("\t", m.GetRootSubPath())
	}
	e.logger.Infoln(fmt.Sprintf("SrcNew Len=%d", len(e.srcNewArr)))
	for _, sn := range e.srcNewArr {
		e.logger.Infoln("\t", sn.GetRootSubPath())
	}
	e.logger.Infoln(fmt.Sprintf("TarNew Len=%d", len(e.tarNewArr)))
	for _, tn := range e.tarNewArr {
		e.logger.Infoln("\t", tn.GetRootSubPath())
	}
}

func (e *syncExecutor) fileFitting(fileInfo os.FileInfo) bool {
	if nil == fileInfo {
		return false
	}
	// 名称不匹配
	if !e.target.CheckFileFitting(fileInfo.Name()) {
		return false
	}
	return true
}

func (e *syncExecutor) dirFitting(dirInfo os.FileInfo) bool {
	if nil == dirInfo {
		return false
	}
	if !e.target.CheckDirFitting(dirInfo.Name()) {
		return false
	}
	return true
}

func (e *syncExecutor) OperateSets() {
	srcArr := e.srcList.GetResults()
	tarArr := e.tarList.GetResults()

	idx0, idx1 := 0, 0
	sLen, tLen := len(srcArr), len(tarArr)
	minLen := mathx.MinInt(sLen, tLen)
	srcIdxArr := make([]int, 0, minLen)
	tarIdxArr := make([]int, 0, minLen)
	for idx0 < sLen && idx1 < tLen { // 找相同
		sInfo := srcArr[idx0]
		tInfo := tarArr[idx1]
		//fmt.Println("C:", sInfo.GetRootSubPath(), tInfo.GetRootSubPath())
		if sInfo.GetRootSubPath() == tInfo.GetRootSubPath() {
			srcIdxArr = append(srcIdxArr, idx0)
			tarIdxArr = append(tarIdxArr, idx1)
			idx0, idx1 = idx0+1, idx1+1
			continue
		}
		if sInfo.LessTo(tInfo) {
			idx0 += 1
		} else {
			idx1 += 1
		}
	}
	sameSize := len(srcIdxArr)
	e.mixedArr = make([]internal.IPathInfo, 0, sameSize)
	e.srcNewArr = make([]internal.IPathInfo, 0, sLen-sameSize)
	for idx0, idx1 = 0, 0; idx0 < sLen; {
		if idx1 < sameSize && idx0 == srcIdxArr[idx1] { // 相同
			e.mixedArr = append(e.mixedArr, srcArr[idx0])
			idx0, idx1 = idx0+1, idx1+1
			continue
		}
		e.srcNewArr = append(e.srcNewArr, srcArr[idx0])
		idx0 += 1
	}
	if e.double { // 双向
		e.tarNewArr = make([]internal.IPathInfo, 0, tLen-sameSize)
		for idx0, idx1 = 0, 0; idx0 < tLen; {
			if idx1 < sameSize && idx0 == srcIdxArr[idx1] { // 相同
				idx0, idx1 = idx0+1, idx1+1
				continue
			}
			e.tarNewArr = append(e.tarNewArr, tarArr[idx0])
			idx0 += 1
		}
	}
}
