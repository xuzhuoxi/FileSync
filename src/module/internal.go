package module

import (
	"errors"
	"os"
	"time"
)

type iModuleExecutor interface {
	// 初始化Log
	initArgs()
	// 初始化处理列表
	initExecuteList()
	// 处理列表中文件
	execList()
}

// 内部函数

func getRuneCount(str string, r rune) int {
	runes := []rune(str)
	rs := 0
	for _, o := range runes {
		if r == o {
			rs += 1
		}
	}
	return rs
}

func cloneTime(tarPath string, srcInfo os.FileInfo) {
	os.Chtimes(tarPath, srcInfo.ModTime(), srcInfo.ModTime())
}

func compareTime(tarPath string, mt time.Time) bool {
	fileInfo, err := os.Stat(tarPath)
	// 不存在
	if err != nil && !errors.Is(err, os.ErrExist) {
		return true
	}
	return fileInfo.ModTime().UnixNano() < mt.UnixNano()
}
