package internal

import "os"

type iModuleExecutor interface {
	// 初始化Log
	initArgs()
	// 初始化处理列表
	initExecuteList()
	// 处理列表中文件
	execList()

	// 文件过滤
	fileFitting(fileInfo os.FileInfo) bool
	// 目录过滤
	dirFitting(dirInfo os.FileInfo) bool
}
