package infra

const (
	ArgDouble      = "/d"    //双向
	ArgIgnoreEmpty = "/i"    //忽略空目录
	ArgLogFile     = "/Lf"   //记录日志
	ArgLogPrint    = "/Lp"   //控制台打印信息
	ArgRecurse     = "/r"    //递归
	ArgStable      = "/s"    //保持文件目录结构
	ArgFile        = "/file" //单文件模式
	ArgTimeUpdate  = "/time" //若目标文件比源文件旧，更新目标文件
	ArgSizeUpdate  = "/size" //若目标文件比源文件大，更新目标文件
	ArgMd5Update   = "/md5"  //若目标文件md5与源文件不一致，更新目标文件
)

const argStart = '/'

type ArgMark int

const (
	MarkDouble ArgMark = 1 << iota
	MarkIgnoreEmpty
	MarkLogFile
	MarkLogPrint
	MarkRecurse
	MarkStable
	MarkFile
	MarkTimeUpdate
	MarkSizeUpdate
	MarkMd5Update
)

const (
	ClearArgMark  = MarkLogFile | MarkLogPrint | MarkRecurse
	CopyArgMark   = MarkIgnoreEmpty | MarkLogFile | MarkLogPrint | MarkRecurse | MarkStable | MarkFile | MarkTimeUpdate | MarkSizeUpdate | MarkMd5Update
	DeleteArgMark = MarkLogFile | MarkLogPrint | MarkRecurse
	MoveArgMark   = MarkIgnoreEmpty | MarkLogFile | MarkLogPrint | MarkRecurse | MarkStable | MarkFile | MarkTimeUpdate | MarkSizeUpdate | MarkMd5Update
	SyncArgMark   = MarkDouble | MarkIgnoreEmpty | MarkLogFile | MarkLogPrint | MarkRecurse | MarkTimeUpdate | MarkSizeUpdate | MarkMd5Update
)

var (
	ClearArgs  = []string{ArgLogFile, ArgLogPrint, ArgRecurse}
	CopyArgs   = []string{ArgIgnoreEmpty, ArgLogFile, ArgLogPrint, ArgRecurse, ArgStable, ArgFile, ArgTimeUpdate, ArgSizeUpdate, ArgMd5Update}
	DeleteArgs = []string{ArgLogFile, ArgLogPrint, ArgRecurse}
	MoveArgs   = []string{ArgIgnoreEmpty, ArgLogFile, ArgLogPrint, ArgRecurse, ArgStable, ArgFile, ArgTimeUpdate, ArgSizeUpdate, ArgMd5Update}
	SyncArgs   = []string{ArgDouble, ArgIgnoreEmpty, ArgLogFile, ArgLogPrint, ArgRecurse, ArgTimeUpdate, ArgSizeUpdate, ArgMd5Update}
)

const DefaultArgMark = MarkLogFile | MarkLogPrint

var (
	mapValue2Mark = make(map[string]ArgMark)
	mapMark2Value = make(map[ArgMark]string)
)

func init() {
	mapMark2Value[MarkDouble] = ArgDouble
	mapMark2Value[MarkIgnoreEmpty] = ArgIgnoreEmpty
	mapMark2Value[MarkLogFile] = ArgLogFile
	mapMark2Value[MarkLogPrint] = ArgLogPrint
	mapMark2Value[MarkRecurse] = ArgRecurse
	mapMark2Value[MarkStable] = ArgStable
	mapMark2Value[MarkFile] = ArgFile
	mapMark2Value[MarkTimeUpdate] = ArgTimeUpdate
	mapMark2Value[MarkSizeUpdate] = ArgSizeUpdate
	mapMark2Value[MarkMd5Update] = ArgMd5Update

	mapValue2Mark[ArgDouble] = MarkDouble
	mapValue2Mark[ArgIgnoreEmpty] = MarkIgnoreEmpty
	mapValue2Mark[ArgLogFile] = MarkLogFile
	mapValue2Mark[ArgLogPrint] = MarkLogPrint
	mapValue2Mark[ArgRecurse] = MarkRecurse
	mapValue2Mark[ArgStable] = MarkStable
	mapValue2Mark[ArgFile] = MarkFile
	mapValue2Mark[ArgTimeUpdate] = MarkTimeUpdate
	mapValue2Mark[ArgSizeUpdate] = MarkSizeUpdate
	mapValue2Mark[ArgMd5Update] = MarkMd5Update
}

// 检查参数码位的匹配情况
func (m ArgMark) MatchArg(arg ArgMark) bool {
	return int(m&arg) > 0
}

// 字符串参数转换为码位
func ValuesToMarks(params string) ArgMark {
	values := SplitArgs(params)
	if nil == values {
		return 0
	}
	var rs ArgMark = 0
	for _, value := range values {
		rs = rs | value2Mark(value)
	}
	return ArgMark(rs)
}

// 拆分参数字符串为参数数组
func SplitArgs(args string) []string {
	if "" == args {
		return nil
	}

	start := 0
	var rs []string
	for index := range args {
		if args[index] == argStart {
			if index > start {
				rs = append(rs, args[start:index])
			}
			start = index
		}
	}
	rs = append(rs, args[start:])
	return rs
}

// 取对应模式下支持的参数范围
func GetSupportArgs(mode RuntimeMode) []string {
	switch mode {
	case ModeClear:
		return ClearArgs
	case ModeCopy:
		return CopyArgs
	case ModeDelete:
		return DeleteArgs
	case ModeMove:
		return MoveArgs
	case ModeSync:
		return SyncArgs
	}
	return nil
}

func mark2Value(mark ArgMark) string {
	if v, ok := mapMark2Value[mark]; ok {
		return v
	}
	return ""
}

func value2Mark(value string) ArgMark {
	if v, ok := mapValue2Mark[value]; ok {
		return v
	}
	return 0
}
