package infra

import (
	"errors"
	"os"
	"time"
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

// 修改文件或目录的访问时间和修改时间
func SetModTime(tarPath string, time time.Time) {
	os.Chtimes(tarPath, time, time)
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
