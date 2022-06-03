package infra

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/cryptox"
	"os"
)

//
//type CompareType int
//
//const (
//	// 不比较
//	CompareNone CompareType = 1 << iota
//	// 按修改时间比较
//	CompareTime
//	// 按文件大小比较
//	CompareSize
//	// 按md5比较
//	CompareMd5
//)

// 比较源文件与目标文件的修改时间
// >0：	源文件新
// =0：	时间相同
// <0：	源文件旧
func CompareWithTime(srcInfo, tarInfo os.FileInfo) int64 {
	return srcInfo.ModTime().UnixNano() - tarInfo.ModTime().UnixNano()
}

// 比较源文件与目标文件的大小
// >0：	源文件新
// =0：	时间相同
// <0：	源文件旧
func CompareWithSize(srcInfo, tarInfo os.FileInfo) int64 {
	return srcInfo.Size() - tarInfo.Size()
}

// 比较源文件与目标文件的md5
// true：	相同
// false：	不同
func CompareWithMd5(srcPath, tarPath string) bool {
	md5Src := cryptox.Md5File(srcPath)
	md5Tar := cryptox.Md5File(tarPath)
	fmt.Println(srcPath, md5Src, "|", tarPath, md5Tar, "|", md5Src == md5Tar)
	return md5Src == md5Tar
}
