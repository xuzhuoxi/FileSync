package module

import (
	"errors"
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

func (e *moveExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Move", Mode: infra.ModeMoveValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *moveExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget := infra.NewRuntimeTarget(config)
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
	e.ignore = argsMark.MatchArg(infra.ArgMarkIgnore)
	e.recurse = argsMark.MatchArg(infra.ArgMarkRecurse)
	e.stable = argsMark.MatchArg(infra.ArgMarkStable)
	e.update = argsMark.MatchArg(infra.ArgMarkUpdate)
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
		if moveInfo.FileInfo.IsDir() { // 忽略目录
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
		if !moveInfo.FileInfo.IsDir() { //忽略文件
			continue
		}
		// 当前目录不为空
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

func (e *moveExecutor) checkSrcRoot(rootIndex int, srcRoot string) {
	dir, fileName := filex.Split(srcRoot)
	// 名称不匹配
	if !e.target.CheckNameFitting(fileName) {
		return
	}
	fullSrcRoot := filex.Combine(infra.RunningDir, srcRoot)
	fileInfo, err := os.Stat(fullSrcRoot)
	// 不存在
	if err != nil && !errors.Is(err, os.ErrExist) {
		e.logger.Warnln(fmt.Sprintf("[copy] Ignore src[%d]: %s", rootIndex, srcRoot))
		return
	}
	// 文件
	if !fileInfo.IsDir() {
		e.appendPath(rootIndex, dir, fileName, fileInfo)
		return
	}
	// 目录
	subPaths, _ := ioutil.ReadDir(fullSrcRoot)
	// 真空目录
	if len(subPaths) == 0 {
		return
	}
	// 遍历
	for _, info := range subPaths {
		if info.IsDir() {
			e.checkDir(rootIndex, srcRoot, info.Name(), info)
		} else {
			e.checkFile(rootIndex, srcRoot, info.Name(), info)
		}
	}
}

func (e *moveExecutor) checkFile(rootIndex int, relativeBase string, relativePath string, fileInfo os.FileInfo) {
	// 名称不匹配
	if !e.target.CheckNameFitting(fileInfo.Name()) {
		return
	}
	e.appendPath(rootIndex, relativeBase, relativePath, fileInfo)
}

func (e *moveExecutor) checkDir(rootIndex int, relativeBase string, relativePath string, fileInfo os.FileInfo) {
	// 名称不匹配
	if !e.target.CheckNameFitting(fileInfo.Name()) {
		return
	}
	// 不忽略空目录
	if !e.ignore {
		e.appendPath(rootIndex, relativeBase, relativePath, fileInfo)
	}
	fullPath := filex.Combine(infra.RunningDir, relativeBase, relativePath)
	subPaths, _ := ioutil.ReadDir(fullPath)
	// 真空目录
	if len(subPaths) == 0 {
		return
	}
	// 不递归
	if !e.recurse {
		return
	}
	// 遍历
	for _, info := range subPaths {
		rp := filex.Combine(relativePath, info.Name())
		if info.IsDir() {
			e.checkDir(rootIndex, relativeBase, rp, info)
		} else {
			e.checkFile(rootIndex, relativeBase, rp, info)
		}
	}
}

func (e *moveExecutor) appendPath(rootIndex int, srcRoot string, relativePath string, fileInfo os.FileInfo) {
	srcRelativePath := filex.Combine(srcRoot, relativePath)
	srcAbsPath := filex.Combine(infra.RunningDir, srcRelativePath)
	var tarRelativePath string
	if e.stable {
		tarRelativePath = filex.Combine(e.target.Tar, relativePath)
	} else {
		tarRelativePath = filex.Combine(e.target.Tar, fileInfo.Name())
	}
	tarAbsPath := filex.Combine(infra.RunningDir, tarRelativePath)

	detail := detailPath{
		Index: rootIndex, SrcRoot: srcRoot, RelativePath: relativePath, FileInfo: fileInfo,
		SrcRelativePath: srcRelativePath, SrcAbsPath: srcAbsPath,
		TarRelativePath: tarRelativePath, TarAbsPath: tarAbsPath}
	e.moveList = append(e.moveList, detail)
}
