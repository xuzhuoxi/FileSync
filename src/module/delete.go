package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module/internal"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"os"
)

func newDeleteExecutor() IModeExecutor {
	return &deleteExecutor{list: internal.NewPathStrList(0, 128)}
}

type deleteExecutor struct {
	task *infra.RuntimeTask
	list internal.IPathStrList

	logger  logx.ILogger
	recurse bool // 递归，查找文件时使用

	tempSrcInfo infra.SrcInfo // 运行时临时
}

func (e *deleteExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTask{Name: "Delete", Mode: infra.ModeDeleteValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTask(config)
}

func (e *deleteExecutor) ExecConfigTask(config infra.ConfigTask) {
	runtimeTask, err := infra.NewRuntimeTask(config)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[delete] Err : %v", err))
	}
	e.ExecRuntimeTask(runtimeTask)
}

func (e *deleteExecutor) ExecRuntimeTask(task *infra.RuntimeTask) {
	if nil == task {
		return
	}
	if len(task.SrcArr) == 0 {
		return
	}
	e.task = task
	err := e.initArgs()
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[delete] Init args error='%s'", err))
		return
	}
	e.initExecuteList()
	e.execList()
}

func (e *deleteExecutor) initArgs() error {
	argsMark := e.task.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.recurse = argsMark.MatchArg(infra.MarkRecurse)
	return nil
}

func (e *deleteExecutor) initExecuteList() {
	for index, src := range e.task.SrcArr {
		e.tempSrcInfo = src

		path := filex.Combine(infra.RunningDir, src.FormattedSrc)
		fileInfo, err := os.Stat(path)
		if err != nil && !os.IsExist(err) { //不存在
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]:%s", index, src.OriginalSrc))
			continue
		}
		if fileInfo.IsDir() {
			e.checkDir(path, fileInfo)
		} else {
			e.checkFile(path, fileInfo)
		}
	}
	e.list.Sort()
}

func (e *deleteExecutor) execList() {
	if e.list.Len() == 0 {
		return
	}
	for _, dir := range e.list.GetAll() {
		e.logger.Infoln("[delete] Delete Path:", dir)
		os.RemoveAll(dir)
	}
}

func (e *deleteExecutor) fileFitting(fileInfo os.FileInfo) bool {
	filename := fileInfo.Name()
	// 路径通配符不匹配
	if !e.tempSrcInfo.CheckFitting(filename) {
		return false
	}
	// 名称不匹配
	if !e.task.CheckFileFitting(filename) {
		return false
	}
	return true
}

func (e *deleteExecutor) dirFitting(dirInfo os.FileInfo) bool {
	if !e.task.CheckDirFitting(dirInfo.Name()) { // 过滤不匹配目录
		return false
	}
	return true
}

func (e *deleteExecutor) checkDir(fullPath string, fileInfo os.FileInfo) {
	if !e.recurse { // 非递归
		return
	}
	if !e.dirFitting(fileInfo) {
		return
	}
	e.checkSubPath(fullPath)
}

func (e *deleteExecutor) checkSubPath(fullPath string) {
	subPaths, _ := filex.GetPathsInDir(fullPath, nil)
	if len(subPaths) == 0 {
		return
	}
	for _, subPath := range subPaths {
		fileInfo, _ := os.Stat(subPath)
		if fileInfo.IsDir() {
			e.checkDir(subPath, fileInfo)
		} else {
			e.checkFile(subPath, fileInfo)
		}
	}
}

func (e *deleteExecutor) checkFile(fullPath string, fileInfo os.FileInfo) {
	if !e.fileFitting(fileInfo) {
		return
	}
	e.list.Append(fullPath)
}
