package infra

const (
	ArgDoubleStr      = "/d" //双向
	ArgIgnoreEmptyStr = "/i" //忽略空目录
	ArgLogFileStr     = "/l" //记录日志
	ArgNoCaseStr      = "/n" //大小写无关
	ArgPrintStr       = "/p" //控制台打印信息
	ArgRecurseStr     = "/r" //递归
	ArgStableStr      = "/s" //保持文件目录结构
	ArgUpdateStr      = "/u" //若目标文件比源文件旧，更新目标文件
)

const argStart = '/'

type ArgMark int

const (
	ArgDouble ArgMark = 1 << iota
	ArgIgnoreEmpty
	ArgLogFile
	ArgNoCase
	ArgPrint
	ArgRecurse
	ArgStable
	ArgUpdate
)

const (
	ClearArgMark  = ArgLogFile | ArgNoCase | ArgPrint | ArgRecurse
	CopyArgMark   = ArgIgnoreEmpty | ArgLogFile | ArgNoCase | ArgPrint | ArgRecurse | ArgStable | ArgUpdate
	DeleteArgMark = ArgLogFile | ArgNoCase | ArgPrint | ArgRecurse
	MoveArgMark   = ArgIgnoreEmpty | ArgLogFile | ArgNoCase | ArgPrint | ArgRecurse | ArgStable | ArgUpdate
	SyncArgMark   = ArgIgnoreEmpty | ArgLogFile | ArgNoCase | ArgPrint | ArgRecurse | ArgUpdate
)

var (
	ClearArgs  = []string{ArgNoCaseStr, ArgLogFileStr, ArgPrintStr, ArgRecurseStr}
	CopyArgs   = []string{ArgNoCaseStr, ArgIgnoreEmptyStr, ArgLogFileStr, ArgPrintStr, ArgRecurseStr, ArgStableStr, ArgUpdateStr}
	DeleteArgs = []string{ArgNoCaseStr, ArgIgnoreEmptyStr, ArgLogFileStr, ArgPrintStr, ArgRecurseStr}
	MoveArgs   = []string{ArgNoCaseStr, ArgIgnoreEmptyStr, ArgLogFileStr, ArgPrintStr, ArgRecurseStr, ArgStableStr, ArgUpdateStr}
	SyncArgs   = []string{ArgNoCaseStr, ArgDoubleStr, ArgIgnoreEmptyStr, ArgLogFileStr, ArgPrintStr, ArgRecurseStr, ArgUpdateStr}
)

const DefaultArgMark = ArgLogFile | ArgPrint

var (
	mapValue2Mark = make(map[string]ArgMark)
	mapMark2Value = make(map[ArgMark]string)
)

func init() {
	mapMark2Value[ArgDouble] = ArgDoubleStr
	mapMark2Value[ArgIgnoreEmpty] = ArgIgnoreEmptyStr
	mapMark2Value[ArgLogFile] = ArgLogFileStr
	mapMark2Value[ArgNoCase] = ArgNoCaseStr
	mapMark2Value[ArgPrint] = ArgPrintStr
	mapMark2Value[ArgRecurse] = ArgRecurseStr
	mapMark2Value[ArgStable] = ArgStableStr
	mapMark2Value[ArgUpdate] = ArgUpdateStr

	mapValue2Mark[ArgDoubleStr] = ArgDouble
	mapValue2Mark[ArgIgnoreEmptyStr] = ArgIgnoreEmpty
	mapValue2Mark[ArgLogFileStr] = ArgLogFile
	mapValue2Mark[ArgNoCaseStr] = ArgNoCase
	mapValue2Mark[ArgPrintStr] = ArgPrint
	mapValue2Mark[ArgRecurseStr] = ArgRecurse
	mapValue2Mark[ArgStableStr] = ArgStable
	mapValue2Mark[ArgUpdateStr] = ArgUpdate
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
