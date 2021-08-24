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

func newCopyExecutor() IModeExecutor {
	return &copyExecutor{srcList: newDetailPathList(0, 128)}
}

type copyExecutor struct {
	target  *infra.RuntimeTarget
	srcList detailPathList

	logger  logx.ILogger
	ignore  bool // 处理复制列表时使用
	recurse bool // 处理复制列表时使用
	stable  bool // 真实复制时使用
	update  bool // 真实复制时使用
}

func (e *copyExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Copy", Mode: infra.ModeCopyValue, Src: src,
		Include: include, Exclude: exclude, Args: args, Case: wildcardCase}
	e.ExecConfigTarget(config)
}

func (e *copyExecutor) ExecConfigTarget(config infra.ConfigTarget) {
	runtimeTarget := infra.NewRuntimeTarget(config)
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
	e.ignore = argsMark.MatchArg(infra.ArgMarkIgnore)
	e.recurse = argsMark.MatchArg(infra.ArgMarkRecurse)
	e.stable = argsMark.MatchArg(infra.ArgMarkStable)
	e.update = argsMark.MatchArg(infra.ArgMarkUpdate)
}

func (e *copyExecutor) initExecuteList() {
	for index, src := range e.target.SrcArr {
		e.checkSrcRoot(index, src)
	}
	e.srcList.Sort()
}

func (e *copyExecutor) execList() {
	e.logger.Infoln(fmt.Sprintf("[copy] Start: Copy(len=%d).", e.srcList.Len()))
	count := 0
	for _, srcInfo := range e.srcList {
		srcFullPath := srcInfo.GetFullPath()
		var tarFullPath string
		if e.stable {
			tarFullPath = filex.Combine(infra.RunningDir, e.target.Tar, srcInfo.relativePath)
		} else {
			tarFullPath = filex.Combine(infra.RunningDir, e.target.Tar, srcInfo.fileInfo.Name())
		}
		// 忽略新文件
		if e.update && !srcInfo.fileInfo.IsDir() && !compareTime(tarFullPath, srcInfo.fileInfo.ModTime()) {
			continue
		}
		e.doCopy(srcFullPath, tarFullPath, srcInfo.fileInfo)
		count += 1
	}
	e.logger.Infoln(fmt.Sprintf("[copy] Finish: Copy(len=%d), Ignore(len=%d).", count, e.srcList.Len()-count))
}

func (e *copyExecutor) doCopy(srcFullPath, tarFullPath string, srcFileInfo os.FileInfo) {
	e.logger.Infoln(fmt.Sprintf("[copy] Copy file '%s' \n\t\t=> '%s'", srcFullPath, tarFullPath))
	if srcFileInfo.IsDir() {
		os.MkdirAll(tarFullPath, srcFileInfo.Mode())
	} else {
		filex.CopyAuto(srcFullPath, tarFullPath, srcFileInfo.Mode())
	}
	cloneTime(tarFullPath, srcFileInfo)
}

func (e *copyExecutor) checkSrcRoot(rootIndex int, srcRoot string) {
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

func (e *copyExecutor) checkFile(rootIndex int, relativeBase string, relativePath string, fileInfo os.FileInfo) {
	// 名称不匹配
	if !e.target.CheckNameFitting(fileInfo.Name()) {
		return
	}
	detail := detailPath{index: rootIndex, relativeBase: relativeBase, relativePath: relativePath, fileInfo: fileInfo}
	e.srcList = append(e.srcList, detail)
}

func (e *copyExecutor) checkDir(rootIndex int, relativeBase string, relativePath string, fileInfo os.FileInfo) {
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

func (e *copyExecutor) appendPath(rootIndex int, relativeBase string, relativePath string, fileInfo os.FileInfo) {
	detail := detailPath{index: rootIndex, relativeBase: relativeBase, relativePath: relativePath, fileInfo: fileInfo}
	e.srcList = append(e.srcList, detail)
}
