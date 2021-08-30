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

func CloneTime(tarPath string, srcInfo os.FileInfo) {
	os.Chtimes(tarPath, srcInfo.ModTime(), srcInfo.ModTime())
}

func CheckPathByTime(tarPath string, mt time.Time) bool {
	fileInfo, err := os.Stat(tarPath)
	// 不存在
	if err != nil && !errors.Is(err, os.ErrExist) {
		return true
	}
	return fileInfo.ModTime().UnixNano() < mt.UnixNano()
}
