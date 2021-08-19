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
	logger  logx.ILogger
	srcList detailPathList

	force   bool
	recurse bool
	ignore  bool
	stable  bool
}

func (e *copyExecutor) Exec(src, tar, include, exclude, args string, wildcardCase bool) {
	config := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
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
	e.initLogger(target.ArgsMark)
	e.initExecuteList()
	e.execList()
}

func (e *copyExecutor) initLogger(mark infra.ArgMark) {
	e.logger = infra.GenLogger(mark)
}

func (e *copyExecutor) initExecuteList() {
	argsMark := e.target.ArgsMark
	e.force = argsMark.MatchArg(infra.ArgMarkForce)
	e.ignore = argsMark.MatchArg(infra.ArgMarkIgnore)
	e.stable = argsMark.MatchArg(infra.ArgMarkStable)
	e.recurse = argsMark.MatchArg(infra.ArgMarkRecurse)

	e.initSrcList()
}

func (e *copyExecutor) execList() {
	e.logger.Infoln("execList:", e.srcList.Len())
	for _, path := range e.srcList {
		e.logger.Infoln(path.GetFullPath())
	}
}

func (e *copyExecutor) initSrcList() {
	for index, src := range e.target.SrcArr {
		e.checkSrcRoot(index, src)
	}
	e.srcList.Sort()
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
