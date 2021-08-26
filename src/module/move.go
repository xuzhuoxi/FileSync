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

func newMoveExecutor() IModeExecutor {
	return &moveExecutor{}
}

type moveExecutor struct {
	target   *infra.RuntimeTarget
	moveList detailPathList

	logger  logx.ILogger
	ignore  bool // 处理复制列表时使用
	recurse bool // 处理复制列表时使用
	stable  bool // 处理复制列表时使用
	update  bool // 真实复制时使用
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
}

func (e *moveExecutor) initExecuteList() {
	for index, src := range e.target.SrcArr {
		e.checkSrcRoot(index, src)
	}
	e.moveList.Sort()
}

func (e *moveExecutor) execList() {
	e.logger.Infoln(fmt.Sprintf("[move] Start(RunningPath='%s', Len=%d).", infra.RunningDir, e.moveList.Len()))
	count := 0
	for _, moveInfo := range e.moveList {
		// 忽略目录
		if moveInfo.FileInfo.IsDir() {
			continue
		}
		// 忽略新文件
		if e.update && !compareTime(moveInfo.TarAbsPath, moveInfo.FileInfo.ModTime()) {
			e.logger.Infoln(fmt.Sprintf("[move] Ignore by '/u': '%s'", moveInfo.SrcAbsPath))
			continue
		}
		e.doMoveFile(moveInfo)
		count += 1
	}
	for _, moveInfo := range e.moveList {
		// 忽略文件
		if !moveInfo.FileInfo.IsDir() {
			continue
		}
		// 忽略非空目录
		if !filex.IsEmptyDir(moveInfo.SrcAbsPath) {
			continue
		}
		// 忽略新目录
		if e.update && !compareTime(moveInfo.TarAbsPath, moveInfo.FileInfo.ModTime()) {
			continue
		}
		e.doMoveDir(moveInfo)
		count += 1
	}
	ignoreLen := e.moveList.Len() - count
	e.logger.Infoln(fmt.Sprintf("[move] Finish(CopyLen=%d, IgnoreLen=%d).", count, ignoreLen))
}

func (e *moveExecutor) doMoveFile(moveInfo detailPath) {
	e.logger.Infoln(fmt.Sprintf("[move] Move file '%s' => '%s'", moveInfo.SrcRelativePath, moveInfo.TarRelativePath))
	os.Rename(moveInfo.SrcAbsPath, moveInfo.TarAbsPath)
	if moveInfo.FileInfo.IsDir() {
		os.MkdirAll(moveInfo.TarAbsPath, moveInfo.FileInfo.Mode())
		cloneTime(moveInfo.TarAbsPath, moveInfo.FileInfo)
	} else {
		filex.CompletePath(moveInfo.TarAbsPath, moveInfo.FileInfo.Mode())
		os.Rename(moveInfo.SrcAbsPath, moveInfo.TarAbsPath)
	}
}

func (e *moveExecutor) doMoveDir(moveInfo detailPath) {
	e.logger.Infoln(fmt.Sprintf("[move] Move Dir '%s' => '%s'", moveInfo.SrcRelativePath, moveInfo.TarRelativePath))
	os.Rename(moveInfo.SrcAbsPath, moveInfo.TarAbsPath)
	if moveInfo.FileInfo.IsDir() {
		os.MkdirAll(moveInfo.TarAbsPath, moveInfo.FileInfo.Mode())
		cloneTime(moveInfo.TarAbsPath, moveInfo.FileInfo)
	} else {
		os.Rename(moveInfo.SrcAbsPath, moveInfo.TarAbsPath)
	}
}

func (e *moveExecutor) checkSrcRoot(rootIndex int, srcInfo infra.SrcInfo) {
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

func (e *moveExecutor) checkDir(rootIndex int, srcInfo infra.SrcInfo, srcRelativePath string, fileInfo os.FileInfo) {
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
func (e *moveExecutor) checkSubDir(rootIndex int, srcInfo infra.SrcInfo, srcRelativePath string) {
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

func (e *moveExecutor) checkFile(rootIndex int, srcInfo infra.SrcInfo, srcRelativePath string, fileInfo os.FileInfo) {
	if !srcInfo.CheckFitting(fileInfo.Name()) { // 路径匹配不成功
		return
	}
	// 名称不匹配
	if !e.target.CheckFileFitting(fileInfo.Name()) {
		return
	}
	e.appendPath(rootIndex, srcInfo, srcRelativePath, fileInfo)
}

func (e *moveExecutor) appendPath(rootIndex int, srcInfo infra.SrcInfo, relativePath string, fileInfo os.FileInfo) {
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
	e.moveList = append(e.moveList, detail)
}
