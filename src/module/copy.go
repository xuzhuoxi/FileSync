package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"io/ioutil"
	"os"
	"strings"
)

func newCopyExecutor() IModeExecutor {
	return &copyExecutor{copyList: newDetailPathList(0, 128)}
}

type copyExecutor struct {
	target *infra.RuntimeTarget

	logger  logx.ILogger
	ignore  bool // 处理复制列表时使用
	recurse bool // 处理复制列表时使用
	stable  bool // 处理复制列表时使用
	update  bool // 真实复制时使用

	copyList detailPathList
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
}

func (e *copyExecutor) initExecuteList() {
	for index, src := range e.target.SrcArr {
		e.checkSrcRoot(index, src)
	}
	e.copyList.Sort()
}

func (e *copyExecutor) execList() {
	e.logger.Infoln(fmt.Sprintf("[copy] Start(RunningPath='%s', Len=%d).", infra.RunningDir, e.copyList.Len()))
	count := 0
	for _, copyInfo := range e.copyList {
		// 忽略新文件
		if e.update && !copyInfo.FileInfo.IsDir() && !compareTime(copyInfo.TarAbsPath, copyInfo.FileInfo.ModTime()) {
			e.logger.Infoln(fmt.Sprintf("[copy] Ignore by '/u': '%s'", copyInfo.SrcAbsPath))
			continue
		}
		e.doCopy(copyInfo)
		count += 1
	}
	ignoreLen := e.copyList.Len() - count
	e.logger.Infoln(fmt.Sprintf("[copy] Finish(CopyLen=%d, IgnoreLen=%d).", count, ignoreLen))
}

func (e *copyExecutor) doCopy(copyInfo detailPath) {
	e.logger.Infoln(fmt.Sprintf("[copy] Copy '%s' => '%s'", copyInfo.SrcRelativePath, copyInfo.TarRelativePath))
	if copyInfo.FileInfo.IsDir() {
		os.MkdirAll(copyInfo.TarAbsPath, copyInfo.FileInfo.Mode())
	} else {
		filex.CopyAuto(copyInfo.SrcAbsPath, copyInfo.TarAbsPath, copyInfo.FileInfo.Mode())
	}
	cloneTime(copyInfo.TarAbsPath, copyInfo.FileInfo)
}

func (e *copyExecutor) checkSrcRoot(rootIndex int, srcInfo infra.SrcInfo) {
	fullSrcRoot := filex.Combine(infra.RunningDir, srcInfo.FormattedSrc)
	fileInfo, err := os.Stat(fullSrcRoot)
	if err != nil && !os.IsExist(err) { //不存在
		e.logger.Warnln(fmt.Sprintf("[copy] Ignore src[%d]: %s", rootIndex, srcInfo.OriginalSrc))
		return
	}

	if !fileInfo.IsDir() { // 文件
		e.checkFile(rootIndex, srcInfo, "", fileInfo)
		return
	}

	// 目录
	if srcInfo.IncludeSelf {
		e.checkDir(rootIndex, srcInfo, "", fileInfo)
	} else {
		e.checkSubDir(rootIndex, srcInfo, "")
	}
}

func (e *copyExecutor) checkDir(rootIndex int, srcInfo infra.SrcInfo, srcRelativePath string, fileInfo os.FileInfo) {
	// 名称不匹配
	if !e.target.CheckDirFitting(fileInfo.Name()) {
		return
	}
	// 不忽略空目录，把目录都加入到列表中
	if !e.ignore {
		e.appendPath(rootIndex, srcInfo, srcRelativePath, fileInfo)
	}
	// 不递归
	if !e.recurse {
		return
	}
	e.checkSubDir(rootIndex, srcInfo, srcRelativePath)
}
func (e *copyExecutor) checkSubDir(rootIndex int, srcInfo infra.SrcInfo, srcRelativePath string) {
	fullPath := filex.Combine(infra.RunningDir, srcInfo.FormattedSrc, srcRelativePath)
	subPaths, _ := ioutil.ReadDir(fullPath)
	// 真空目录
	if len(subPaths) == 0 {
		return
	}
	// 遍历
	for _, info := range subPaths {
		rp := filex.Combine(srcRelativePath, info.Name())
		if info.IsDir() {
			e.checkDir(rootIndex, srcInfo, rp, info)
		} else {
			e.checkFile(rootIndex, srcInfo, rp, info)
		}
	}
}

func (e *copyExecutor) checkFile(rootIndex int, srcInfo infra.SrcInfo, srcRelativePath string, fileInfo os.FileInfo) {
	if !srcInfo.CheckFitting(fileInfo.Name()) { // 路径匹配不成功
		return
	}
	// 名称不匹配
	if !e.target.CheckFileFitting(fileInfo.Name()) {
		return
	}
	e.appendPath(rootIndex, srcInfo, srcRelativePath, fileInfo)
}

func (e *copyExecutor) appendPath(rootIndex int, srcInfo infra.SrcInfo, relativePath string, fileInfo os.FileInfo) {
	srcRelativePath := filex.Combine(srcInfo.FormattedSrc, relativePath)
	srcAbsPath := filex.Combine(infra.RunningDir, srcRelativePath)
	var tarRelativePath string
	if e.stable { // 保持目录
		if srcInfo.IncludeSelf { // 包含源目录
			_, selfName := filex.Split(srcInfo.FormattedSrc)
			tarRelativePath = filex.Combine(e.target.Tar, selfName, relativePath)
		} else { // 不包含源目录
			tarRelativePath = filex.Combine(e.target.Tar, relativePath)
		}
	} else { // 不保持目录
		tarRelativePath = filex.Combine(e.target.Tar, fileInfo.Name())
	}
	tarAbsPath := filex.Combine(infra.RunningDir, tarRelativePath)

	detail := detailPath{
		Index: rootIndex, SrcInfo: srcInfo, FileInfo: fileInfo,
		SrcRelativePath: srcRelativePath, SrcAbsPath: srcAbsPath,
		TarRelativePath: tarRelativePath, TarAbsPath: tarAbsPath}
	e.copyList = append(e.copyList, detail)
}
