package module

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module/internal"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/mathx"
	"os"
	"strings"
)

func newSyncExecutor() IModeExecutor {
	return &syncExecutor{srcList: internal.NewPathSearcher(), tarList: internal.NewPathSearcher()}
}

type syncExecutor struct {
	task   *infra.RuntimeTask
	srcDir string
	tarDir string

	logger     logx.ILogger
	double     bool // 双向同步
	ignore     bool // 忽略空目录，查找文件时使用
	recurse    bool // 保持目录结构，处理文件时使用
	timeUpdate bool // 只处理新时间文件，处理文件时使用
	sizeUpdate bool // 只处理size更大文件，处理文件时使用

	srcList internal.IPathSearcher
	tarList internal.IPathSearcher

	mixedArr  []internal.IPathInfo
	srcNewArr []internal.IPathInfo
	tarNewArr []internal.IPathInfo
}

func (e *syncExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTask{Name: "Sync", Mode: infra.ModeSyncValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTask(config)
}

func (e *syncExecutor) ExecConfigTask(config infra.ConfigTask) {
	runtimeTask, err := infra.NewRuntimeTask(config)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Err : %v", err))
	}
	e.ExecRuntimeTask(runtimeTask)
}

func (e *syncExecutor) ExecRuntimeTask(task *infra.RuntimeTask) {
	if nil == task {
		return
	}
	if len(task.SrcArr) != 1 {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Src Len Err! "))
		return
	}
	e.srcDir = strings.TrimSpace(task.SrcArr[0].FormattedSrc)
	if e.srcDir == "" {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Src Empty Err! "))
		return
	}
	e.tarDir = strings.TrimSpace(task.Tar)
	if e.tarDir == "" {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Tar Empty Err! "))
		return
	}
	e.task = task
	err := e.initArgs()
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[sync] Init args error='%s'", err))
		return
	}
	e.initExecuteList()
	e.execList()
}

func (e *syncExecutor) initArgs() (err error) {
	argsMark := e.task.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.double = argsMark.MatchArg(infra.MarkDouble)
	e.ignore = argsMark.MatchArg(infra.MarkIgnoreEmpty)
	e.recurse = argsMark.MatchArg(infra.MarkRecurse)
	e.timeUpdate = argsMark.MatchArg(infra.MarkTimeUpdate)
	e.sizeUpdate = argsMark.MatchArg(infra.MarkSizeUpdate)

	if e.double && e.timeUpdate && e.sizeUpdate {
		err = errors.New(fmt.Sprintf("[sync] Error with args: '%s'存在时,'%s'与'%s'互斥! ",
			infra.ArgDouble, infra.ArgTimeUpdate, infra.ArgSizeUpdate))
		return
	}
	if !e.timeUpdate && !e.sizeUpdate {
		err = errors.New(fmt.Sprintf("[sync] Error with args: '%s'不存在时,'%s'与'%s'有且至少有一个! ",
			infra.ArgDouble, infra.ArgTimeUpdate, infra.ArgSizeUpdate))
		return
	}

	e.srcList.SetParams(e.recurse, !e.ignore, e.logger)
	e.tarList.SetParams(e.recurse, !e.ignore, e.logger)
	return nil
}

func (e *syncExecutor) initExecuteList() {
	// 查找源
	e.srcList.SetFitting(e.fileFitting, e.dirFitting)
	e.srcList.InitSearcher()
	src := e.task.SrcArr[0].FormattedSrc // 同步时，源目录有且只有一个目录
	e.srcList.Search(src, false)
	e.srcList.SortResults()

	// 查找目标
	e.tarList.SetFitting(e.fileFitting, e.dirFitting)
	e.tarList.InitSearcher()
	e.tarList.Search(e.task.Tar, false)
	e.tarList.SortResults()

	// 源集合与目录集合的交集与差集计算
	e.OperateSets()
}

func (e *syncExecutor) execList() {
	e.logger.Infoln(fmt.Sprintf("[sync] SrcNew(Len=%d),TarNew(Len=%d),Mixed(Len=%d)", len(e.srcNewArr), len(e.tarNewArr), len(e.mixedArr)))
	e.execMixedList()
	e.execSrcNew()
	e.execTarNew()
}

func (e *syncExecutor) execMixedList() {
	e.logger.Infoln(fmt.Sprintf("[sync] Mixed(Len=%d) Src=>Tar", len(e.mixedArr)))
	if e.double {
		e.execMixedDouble()
	} else {
		e.execMixedMirroring()
	}
}

// 镜像处理共同文件
func (e *syncExecutor) execMixedDouble() {
	for _, m := range e.mixedArr {
		srcRelative := m.GetRelativePath()
		srcFull := m.GetFullPath()
		tarFull := m.GenFullPath(e.tarDir)
		srcFileInfo := infra.GetFileInfo(srcFull)
		tarFileInfo := infra.GetFileInfo(tarFull)
		if e.timeUpdate {
			cr := infra.CompareWithTime(srcFileInfo, tarFileInfo)
			if cr == 0 {
				e.logger.Infoln(fmt.Sprintf("[sync] Ignored by '%s':'%s'", infra.ArgTimeUpdate, srcRelative))
				continue
			}
			if cr > 0 {
				e.src2tar(srcFull, tarFull, m.GetRootSubPath())
			} else {
				e.tar2src(srcFull, tarFull, m.GetRootSubPath())
			}
			continue
		}
		if e.sizeUpdate {
			cr := infra.CompareWithSize(srcFileInfo, tarFileInfo)
			if cr == 0 {
				e.logger.Infoln(fmt.Sprintf("[sync] Ignored by '%s':'%s'", infra.ArgSizeUpdate, srcRelative))
				continue
			}
			if cr > 0 {
				e.src2tar(srcFull, tarFull, m.GetRootSubPath())
			} else {
				e.tar2src(srcFull, tarFull, m.GetRootSubPath())
			}
			continue
		}
		e.logger.Infoln(fmt.Sprintf("[sync] Ignored by doesn't meet the conditions:'%s'", srcRelative))
	}
}

// 镜像处理共同文件
func (e *syncExecutor) execMixedMirroring() {
	for _, m := range e.mixedArr {
		srcRelative := m.GetRelativePath()
		srcFull := m.GetFullPath()
		tarFull := m.GenFullPath(e.tarDir)
		srcFileInfo := infra.GetFileInfo(srcFull)
		tarFileInfo := infra.GetFileInfo(tarFull)
		if e.timeUpdate && infra.CompareWithTime(srcFileInfo, tarFileInfo) <= 0 {
			e.logger.Infoln(fmt.Sprintf("[sync] Ignored by '%s':'%s'", infra.ArgTimeUpdate, srcRelative))
			continue
		}
		if e.sizeUpdate && infra.CompareWithSize(srcFileInfo, tarFileInfo) <= 0 {
			e.logger.Infoln(fmt.Sprintf("[sync] Ignored by '%s':'%s'", infra.ArgSizeUpdate, srcRelative))
			continue
		}
		e.src2tar(srcFull, tarFull, m.GetRootSubPath())
	}
}

func (e *syncExecutor) src2tar(srcFull string, tarFull string, subRelative string) {
	e.logger.Infoln(fmt.Sprintf("[sync] => '%s'", subRelative))
	internal.DoCopy(srcFull, tarFull, nil)
}

func (e *syncExecutor) tar2src(srcFull string, tarFull string, subRelative string) {
	e.logger.Infoln(fmt.Sprintf("[sync] <= '%s'", subRelative))
	internal.DoCopy(tarFull, srcFull, nil)
}

func (e *syncExecutor) execSrcNew() {
	e.logger.Infoln(fmt.Sprintf("[sync] SrcNew(Len=%d) Src=>Tar", len(e.srcNewArr)))
	for _, sn := range e.srcNewArr {
		e.src2tar(sn.GetFullPath(), sn.GenFullPath(e.tarDir), sn.GetRootSubPath())
	}
}

func (e *syncExecutor) execTarNew() {
	if e.double {
		e.logger.Infoln(fmt.Sprintf("[sync] TarNew(Len=%d) Src<=Tar", len(e.tarNewArr)))
		for _, tn := range e.tarNewArr {
			e.tar2src(tn.GenFullPath(e.srcDir), tn.GetFullPath(), tn.GetRootSubPath())
		}
	} else {
		for _, tn := range e.tarNewArr {
			e.logger.Infoln(fmt.Sprintf("[sync] TarNew(Len=%d) Delete '%s'", len(e.tarNewArr), tn.GetRelativePath()))
			filex.Remove(tn.GetFullPath())
		}
	}
}

func (e *syncExecutor) overrideFilter(srcFileInfo, tarFileInfo os.FileInfo) bool {
	if e.timeUpdate && infra.CompareWithTime(srcFileInfo, tarFileInfo) <= 0 {
		return false
	}
	if e.sizeUpdate && infra.CompareWithSize(srcFileInfo, tarFileInfo) <= 0 {
		return false
	}
	return true
}

func (e *syncExecutor) fileFitting(fileInfo os.FileInfo) bool {
	if nil == fileInfo {
		return false
	}
	// 名称不匹配
	if !e.task.CheckFileFitting(fileInfo.Name()) {
		return false
	}
	return true
}

func (e *syncExecutor) dirFitting(dirInfo os.FileInfo) bool {
	if nil == dirInfo {
		return false
	}
	if !e.task.CheckDirFitting(dirInfo.Name()) {
		return false
	}
	return true
}

func (e *syncExecutor) OperateSets() {
	srcArr := e.srcList.GetResults()
	tarArr := e.tarList.GetResults()

	idx0, idx1 := 0, 0
	sLen, tLen := len(srcArr), len(tarArr)
	//e.logger.Debugln(fmt.Sprintf("SrcLen=%d, TarLen=%d", sLen, tLen))
	minLen := mathx.MinInt(sLen, tLen)
	mixSrcIdxArr := make([]int, 0, minLen)
	mixTarIdxArr := make([]int, 0, minLen)
	for idx0 < sLen && idx1 < tLen { // 找相同
		sInfo := srcArr[idx0]
		tInfo := tarArr[idx1]
		//fmt.Println("SubPath:", idx0, idx1, sInfo.GetRootSubPath(), tInfo.GetRootSubPath())
		if sInfo.GetRootSubPath() == tInfo.GetRootSubPath() {
			mixSrcIdxArr = append(mixSrcIdxArr, idx0)
			mixTarIdxArr = append(mixTarIdxArr, idx1)
			idx0, idx1 = idx0+1, idx1+1
			continue
		}
		if sInfo.LessTo(tInfo) {
			idx0 += 1
		} else {
			idx1 += 1
		}
	}
	sameSize := len(mixSrcIdxArr)
	e.mixedArr = make([]internal.IPathInfo, 0, sameSize)
	e.srcNewArr = make([]internal.IPathInfo, 0, sLen-sameSize)
	for idx0, idx1 = 0, 0; idx0 < sLen; idx0 += 1 { //idx0记录SrcArr下标, idx1记录mixSrcIdxArr下标
		if idx1 >= sameSize {
			break
		}
		if idx0 == mixSrcIdxArr[idx1] { // 相同
			e.mixedArr = append(e.mixedArr, srcArr[idx0])
			//e.logger.Debugln(fmt.Sprintf("MixedPath:%s", srcArr[idx0].GetFullPath()))
			idx1 = idx1 + 1
		} else {
			e.srcNewArr = append(e.srcNewArr, srcArr[idx0])
			//e.logger.Debugln(fmt.Sprintf("SrcPath:%s", srcArr[idx0].GetFullPath()))
		}
	}
	e.tarNewArr = make([]internal.IPathInfo, 0, tLen-sameSize)
	for idx0, idx1 = 0, 0; idx0 < tLen; idx0 += 1 { //idx0记录TarArr下标, idx1记录mixTarIdxArr下标
		if idx1 >= sameSize {
			break
		}
		if idx0 == mixTarIdxArr[idx1] { // 相同
			idx1 = idx1 + 1
		} else {
			e.tarNewArr = append(e.tarNewArr, tarArr[idx0])
			//e.logger.Debugln(fmt.Sprintf("TarPath:%s", tarArr[idx0].GetFullPath()))
		}
	}
}
