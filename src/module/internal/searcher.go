package internal

import (
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"io/ioutil"
	"os"
)

func NewPathSearcher() IPathSearcher {
	return &pathSearcher{}
}

type CheckFitting func(fileInfo os.FileInfo) bool

type IPathSearcher interface {
	// 设置
	SetParams(recurse, dirInclude bool, logger logx.ILogger)
	// 设置
	SetFitting(fileFitting, dirFitting CheckFitting)

	// 初始化
	InitSearcher()
	// 查找
	Search(RelativeRoot string, checkRoot bool)
	// 结果排序
	SortResult()
	// 取结果
	GetResults() []IPathInfo
	// 结果数量
	ResultLen() int
}

type pathSearcher struct {
	recurse    bool // 是否递归
	dirInclude bool // 是否记录目录
	logger     logx.ILogger

	fileFitting CheckFitting
	dirFitting  CheckFitting
	resultList  IPathInfoList
}

func (e *pathSearcher) SetParams(recurse, dirInclude bool, logger logx.ILogger) {
	e.recurse, e.dirInclude, e.logger = recurse, dirInclude, logger
}

func (e *pathSearcher) SetFitting(fileFitting, dirFitting CheckFitting) {
	e.fileFitting, e.dirFitting = fileFitting, dirFitting
}

func (e *pathSearcher) InitSearcher() {
	e.resultList = NewPathInfoList(0, 128)
}

func (e *pathSearcher) Search(relativeRoot string, checkRoot bool) {
	fullSrcRoot := filex.Combine(infra.RunningDir, relativeRoot)
	fileInfo, err := os.Stat(fullSrcRoot)
	if err != nil && !os.IsExist(err) { //不存在
		e.logger.Warnln(fmt.Sprintf("[search] Ignore src[%s]", relativeRoot))
		return
	}
	if !fileInfo.IsDir() { // 文件
		e.checkFile(relativeRoot, "", fileInfo)
		return
	}
	// 目录
	if checkRoot {
		parentDir, selfName := filex.Split(relativeRoot)
		e.checkDir(parentDir, selfName, fileInfo)
	} else {
		e.checkSubDir(relativeRoot, "")
	}
}

func (e *pathSearcher) SortResult() {
	e.resultList.Sort()
	//for _, info := range e.resultList.GetAll() {
	//	fmt.Println("列表：", info.GetRelativePath())
	//}
}

func (e *pathSearcher) GetResults() []IPathInfo {
	if e.resultList == nil || e.resultList.Len() == 0 {
		return nil
	}
	return e.resultList.GetAll()
}

func (e *pathSearcher) ResultLen() int {
	if e.resultList == nil {
		return 0
	}
	return e.resultList.Len()
}

func (e *pathSearcher) checkDir(relativeRoot string, rootSubPath string, fileInfo os.FileInfo) {
	if nil != e.dirFitting && !e.dirFitting(fileInfo) {
		return
	}
	if e.dirInclude {
		e.appendPath(relativeRoot, rootSubPath, fileInfo)
	}
	// 不递归
	if !e.recurse {
		return
	}
	e.checkSubDir(relativeRoot, rootSubPath)
}
func (e *pathSearcher) checkSubDir(relativeRoot string, rootSubPath string) {
	fullPath := filex.Combine(infra.RunningDir, relativeRoot, rootSubPath)
	subPaths, _ := ioutil.ReadDir(fullPath)
	// 真空目录
	if len(subPaths) == 0 {
		return
	}
	// 遍历
	for _, info := range subPaths {
		rp := filex.Combine(rootSubPath, info.Name())
		if info.IsDir() {
			e.checkDir(relativeRoot, rp, info)
		} else {
			e.checkFile(relativeRoot, rp, info)
		}
	}
}

func (e *pathSearcher) checkFile(relativeRoot string, rootSubPath string, fileInfo os.FileInfo) {
	if nil != e.fileFitting && !e.fileFitting(fileInfo) {
		return
	}
	e.appendPath(relativeRoot, rootSubPath, fileInfo)
}

func (e *pathSearcher) appendPath(relativeRoot string, rootSubPath string, fileInfo os.FileInfo) {
	pathInfo := &pathInfo{RelativeRoot: relativeRoot, RootSubPath: rootSubPath, FileInfo: fileInfo}
	e.resultList.Append(pathInfo)
}
