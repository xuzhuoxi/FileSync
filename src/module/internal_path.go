package module

// detailPath & detailPathList
//
//func newDetailPathList(ln int, cap int) detailPathList {
//	return make(detailPathList, ln, cap)
//}
//
//
//type detailPath struct {
//	Index    int
//	SrcInfo  infra.SrcInfo // Src信息
//	FileInfo os.FileInfo
//
//	SrcRelativePath string // 运行时 源 相对路径
//	TarRelativePath string // 运行时 目录 相对路径
//
//	SrcAbsPath string // 源 绝对路径
//	TarAbsPath string // 目标 绝对路径
//}
//
//type detailPathList []detailPath
//
//func (l detailPathList) Len() int {
//	return len(l)
//}
//
//func (l detailPathList) Less(i, j int) bool {
//	if l[i].Index != l[j].Index {
//		return l[i].Index < l[j].Index
//	}
//	lenI := GetRuneCount(l[i].SrcRelativePath, infra.DirSeparator)
//	lenJ := GetRuneCount(l[j].SrcRelativePath, infra.DirSeparator)
//	if lenI != lenJ {
//		return lenI > lenJ
//	}
//	return l[i].SrcRelativePath < l[j].SrcRelativePath
//}
//
//func (l detailPathList) Swap(i, j int) {
//	l[i], l[j] = l[j], l[i]
//}
//
//func (l detailPathList) Sort() { sort.Sort(l) }
