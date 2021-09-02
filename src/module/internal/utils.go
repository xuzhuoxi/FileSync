package internal

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"os"
)

// 生成目录路径，
// relativePath: 目标相对路径
// fullPath: 目标绝对路径
func GetTarPaths(pathInfo IPathInfo, stable bool, tarRelative string) (relativePath, fullPath string) {
	fileInfo := pathInfo.GetFileInfo()
	if stable { // 保持目录结构
		relativePath = filex.Combine(tarRelative, pathInfo.GetRootSubPath())
	} else { // 不保持目录
		relativePath = filex.Combine(tarRelative, fileInfo.Name())
	}
	fullPath = filex.Combine(infra.RunningDir, relativePath)
	return
}

// 复制
func DoCopy(srcPath, tarPath string, doFilter func(srcFileInfo, tarFileInfo os.FileInfo) bool) {
	srcFileInfo := infra.GetFileInfo(srcPath)
	if nil == srcFileInfo {
		return
	}
	tarFileInfo := infra.GetFileInfo(tarPath)
	if nil != doFilter && !doFilter(srcFileInfo, tarFileInfo) {
		return
	}
	if srcFileInfo.IsDir() {
		os.MkdirAll(tarPath, srcFileInfo.Mode())
	} else {
		filex.CopyAuto(srcPath, tarPath, srcFileInfo.Mode())
	}
	infra.SetModTime(tarPath, srcFileInfo.ModTime())
}

// 移动
func DoMove(srcPath, tarPath string, doFilter func(srcFileInfo, tarFileInfo os.FileInfo) bool) {
	srcFileInfo := infra.GetFileInfo(srcPath)
	if nil == srcFileInfo {
		return
	}
	tarFileInfo := infra.GetFileInfo(tarPath)
	if nil != doFilter && !doFilter(srcFileInfo, tarFileInfo) {
		return
	}

	if nil != tarFileInfo {
		filex.Remove(tarPath)
	} else {
		filex.CompleteParentPath(tarPath, srcFileInfo.Mode())
	}
	os.Rename(srcPath, tarPath)

	if !srcFileInfo.IsDir() || filex.IsEmptyDir(srcPath) {
		filex.Remove(srcPath)
	}
}
