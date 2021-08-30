package internal

import (
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
)

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
