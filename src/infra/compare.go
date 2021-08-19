package infra

type CompareType int

const (
	// 不比较
	CompareNone CompareType = 1 << iota
	// 按修改时间比较
	CompareTime
	// 按文件大小比较
	CompareSize
	// 按md5比较
	CompareMd5
)
