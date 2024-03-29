package module

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module/internal"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/slicex"
	"os"
)

func newClearExecutor() IModeExecutor {
	return &clearExecutor{list: internal.NewPathStrList(0, 128)}
}

type clearExecutor struct {
	task *infra.RuntimeTask
	list internal.IPathStrList

	logger  logx.ILogger
	recurse bool // 递归，查找文件时使用

	tempSrcInfo infra.SrcInfo
}

func (e *clearExecutor) Exec(src, tar, include, exclude, args string) {
	config := infra.ConfigTask{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Args: args}
	e.ExecConfigTask(config)
}

func (e *clearExecutor) ExecConfigTask(cfgTask infra.ConfigTask) {
	runtimeTask, err := infra.NewRuntimeTask(cfgTask)
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[clear] Err : %v", err))
	}
	e.ExecRuntimeTask(runtimeTask)
}

func (e *clearExecutor) ExecRuntimeTask(task *infra.RuntimeTask) {
	if nil == task {
		return
	}
	if len(task.SrcArr) == 0 {
		return
	}
	e.task = task
	err := e.initArgs()
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[clear] Init args error='%s'", err))
		return
	}
	e.initExecuteList()
	e.execList()
}

func (e *clearExecutor) initArgs() error {
	argsMark := e.task.ArgsMark
	e.logger = infra.GenLogger(argsMark)
	e.recurse = argsMark.MatchArg(infra.MarkRecurse)
	//fmt.Println("initArgs", e.recurse)
	return nil
}

func (e *clearExecutor) initExecuteList() {
	for index, src := range e.task.SrcArr {
		e.tempSrcInfo = src
		path := filex.Combine(infra.RunningDir, src.FormattedSrc)
		if !filex.IsFolder(path) {
			e.logger.Warnln(fmt.Sprintf("[clear] Ignore src[%d]: %s", index, src.OriginalSrc))
			continue
		}
		if src.IncludeSelf {
			e.checkDir(path)
		} else {
			e.checkSubDir(path)
		}
	}
}

func (e *clearExecutor) execList() {
	if e.list.Len() == 0 {
		return
	}
	all := e.list.GetAll()
	slicex.ReverseString(all)
	for _, dir := range all {
		if e.isDirEmpty(dir) {
			e.logger.Infoln(fmt.Sprintf("[clear] Clear Folder='%s'", dir))
			os.Remove(dir)
		}
	}
}

func (e *clearExecutor) fileFitting(fileInfo os.FileInfo) bool {
	return false
}

func (e *clearExecutor) dirFitting(dirInfo os.FileInfo) bool {
	// 路径通配符不匹配
	filename := dirInfo.Name()
	if !e.tempSrcInfo.CheckFitting(filename) {
		return false
	}
	// 名称不匹配
	if !e.task.CheckDirFitting(filename) {
		return false
	}
	return true
}

func (e *clearExecutor) checkSubDir(fullDir string) {
	dirPaths, _ := filex.GetPathsInDir(fullDir, func(subPath string, info os.FileInfo) bool {
		return info.IsDir()
	})
	if len(dirPaths) == 0 {
		return
	}
	for _, dir := range dirPaths {
		e.checkDir(dir)
	}
}

func (e *clearExecutor) checkDir(fullDir string) {
	fileInfo, err := os.Stat(fullDir)
	if err != nil && !os.IsExist(err) { //不存在
		return
	}
	if !fileInfo.IsDir() { //文件
		return
	}
	fit := e.dirFitting(fileInfo)
	if fit {
		e.appendDir(fullDir)
	}
	e.checkSubDir(fullDir)
}

func (e *clearExecutor) isDirEmpty(fullDir string) bool {
	dirPaths, _ := filex.GetPathsInDir(fullDir, nil)
	return 0 == len(dirPaths)
}

func (e *clearExecutor) isDirDepthEmpty(fullDir string) bool {
	size, _ := filex.GetFolderSize(fullDir)
	return size <= 0
}

func (e *clearExecutor) appendDir(fullDir string) {
	e.list.Append(fullDir)
}
