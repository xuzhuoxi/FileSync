package infra

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/filex"
)

func NewSrcInfo(srcPath string) SrcInfo {
	originalSrc := filex.FormatPath(srcPath)
	dir, filename := filex.Split(srcPath)
	if "" == filename { // 根目录
		return SrcInfo{OriginalSrc: srcPath, FormattedSrc: originalSrc, IncludeSelf: false, Wildcard: ""}
	}
	wildcard := Wildcard(filename)
	if wildcard.IsWildcard() { // 通配符
		return SrcInfo{OriginalSrc: srcPath, FormattedSrc: dir, IncludeSelf: false, Wildcard: wildcard}
	}
	// 文件名 或 目录名
	return SrcInfo{OriginalSrc: srcPath, FormattedSrc: originalSrc, IncludeSelf: true, Wildcard: ""}
}

type SrcInfo struct {
	OriginalSrc  string   // 原始信息
	FormattedSrc string   // 处理后信息
	IncludeSelf  bool     // 目录路径
	Wildcard     Wildcard // 文件通配符
}

func (si SrcInfo) String() string {
	return fmt.Sprintf("SrcInfo{%s, %s, %v, %v}", si.OriginalSrc, si.FormattedSrc, si.IncludeSelf, si.Wildcard)
}

func (si SrcInfo) CheckFitting(filename string) bool {
	return "" == si.Wildcard || si.Wildcard.Match(filename)
}

func NewRuntimeTask(task ConfigTask) (runtimeTask *RuntimeTask, err error) {
	mode := task.GetMode()
	srcArr := task.GetSrcArr()
	tar := task.Tar
	fileIncludes, dirIncludes, iErr := task.GetIncludeArr()
	if nil != iErr {
		return nil, iErr
	}
	fileExcludes, dirExcludes, eErr := task.GetExcludeArr()
	if nil != eErr {
		return nil, eErr
	}
	argMarks := task.GetArgsMark()
	return &RuntimeTask{
		Name:         task.Name,
		Mode:         mode,
		SrcArr:       srcArr,
		Tar:          tar,
		FileIncludes: fileIncludes,
		DirIncludes:  dirIncludes,
		FileExcludes: fileExcludes,
		DirExcludes:  dirExcludes,
		ArgsMark:     argMarks,
	}, nil
}

type RuntimeTask struct {
	RelativeRoot string      // 相对路径的根目录
	Name         string      // 任务名称，用于唯一标识配置
	Mode         RuntimeMode // 任务模式，RuntimeMode
	SrcArr       []SrcInfo   // 任务源信息
	Tar          string      // 任务目标路径
	FileIncludes []Wildcard  // 处理包含的文件名通配符
	DirIncludes  []Wildcard  // 处理包含的目录通配符
	FileExcludes []Wildcard  // 处理排除的文件名通配符
	DirExcludes  []Wildcard  // 处理排除的目录通配符
	ArgsMark     ArgMark     // 任务管理参数
}

func (t *RuntimeTask) CheckFileFitting(filename string) bool {
	return t.checkFitting(filename, t.FileIncludes, t.FileExcludes)
}

func (t *RuntimeTask) CheckDirFitting(filename string) bool {
	return t.checkFitting(filename, t.DirIncludes, t.DirExcludes)
}

func (t *RuntimeTask) MatchArg(param ArgMark) bool {
	return t.ArgsMark.MatchArg(param)
}

func (t *RuntimeTask) checkFitting(filename string, includes []Wildcard, excludes []Wildcard) bool {
	iLen := len(includes)
	eLen := len(excludes)
	if eLen == 0 && iLen == 0 { // 没有include和exclude配置
		return true
	}
	if eLen > 0 && t.checkInWildcard(excludes, filename) { //有exclude配置
		return false
	}
	if iLen > 0 && !t.checkInWildcard(includes, filename) { //有include配置
		return false
	}
	return true
}

func (t *RuntimeTask) checkInWildcard(wildcards []Wildcard, value string) bool {
	for index := range wildcards {
		if wildcards[index].Match(value) {
			return true
		}
	}
	return false
}
