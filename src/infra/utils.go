package infra

import (
	"errors"
	"os"
)

func GetRuneCount(str string, r rune) int {
	runes := []rune(str)
	rs := 0
	for _, o := range runes {
		if r == o {
			rs += 1
		}
	}
	return rs
}

func CloneTime(tarPath string, srcInfo os.FileInfo) {
	os.Chtimes(tarPath, srcInfo.ModTime(), srcInfo.ModTime())
}

// 取路径属性
func GetFileInfo(path string) os.FileInfo {
	fileInfo, err := os.Stat(path)
	// 不存在
	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil
	}
	return fileInfo
}

// 比较源文件与目标文件的修改时间
// 如果源文件更新，则返回true，否则返回false
func CompareWithTime(srcInfo, tarInfo os.FileInfo) bool {
	return srcInfo.ModTime().UnixNano() > tarInfo.ModTime().UnixNano()
}

// 比较源文件与目标文件的大小
// 如果源文件更大，则返回true，否则返回false
func CompareWithSize(srcInfo, tarInfo os.FileInfo) bool {
	return srcInfo.Size() < tarInfo.Size()
}
