package infra

const (
	ArgValueDouble  = "/d" //双向
	ArgValueIgnore  = "/i" //忽略空目录
	ArgValueLog     = "/L" //记录日志
	ArgValueConsole = "/l" //控制台打印信息
	ArgValueRecurse = "/r" //递归
	ArgValueStable  = "/s" //保持文件目录结构
	ArgValueUpdate  = "/u" //若目标文件比源文件旧，更新目标文件
)

const argStart = '/'

type ArgMark int

const (
	ArgMarkDouble ArgMark = 1 << iota
	ArgMarkIgnore
	ArgMarkLogFile
	ArgMarkLogConsole
	ArgMarkRecurse
	ArgMarkStable
	ArgMarkUpdate
)

const (
	ClearArgMark  = ArgMarkLogFile | ArgMarkLogConsole | ArgMarkRecurse
	CopyArgMark   = ArgMarkIgnore | ArgMarkLogFile | ArgMarkLogConsole | ArgMarkRecurse | ArgMarkStable | ArgMarkUpdate
	DeleteArgMark = ArgMarkLogFile | ArgMarkLogConsole | ArgMarkRecurse
	MoveArgMark   = ArgMarkIgnore | ArgMarkLogFile | ArgMarkLogConsole | ArgMarkRecurse | ArgMarkStable | ArgMarkUpdate
	SyncArgMark   = ArgMarkIgnore | ArgMarkLogFile | ArgMarkLogConsole | ArgMarkRecurse | ArgMarkUpdate
)

var (
	ClearArgs = []string{ArgValueLog, ArgValueConsole, ArgValueRecurse}
	CopyArgs  = []string{ArgValueIgnore, ArgValueLog, ArgValueConsole,
		ArgValueRecurse, ArgValueStable, ArgValueUpdate}
	DeleteArgs = []string{ArgValueIgnore, ArgValueLog, ArgValueConsole, ArgValueRecurse}
	MoveArgs   = []string{ArgValueIgnore, ArgValueLog, ArgValueConsole,
		ArgValueRecurse, ArgValueStable, ArgValueUpdate}
	SyncArgs = []string{ArgValueDouble, ArgValueIgnore, ArgValueLog,
		ArgValueConsole, ArgValueRecurse, ArgValueUpdate}
)

const DefaultArgMark = ArgMarkLogFile | ArgMarkLogConsole

var (
	mapValue2Mark = make(map[string]ArgMark)
	mapMark2Value = make(map[ArgMark]string)
)

func init() {
	mapMark2Value[ArgMarkDouble] = ArgValueDouble
	mapMark2Value[ArgMarkIgnore] = ArgValueIgnore
	mapMark2Value[ArgMarkLogFile] = ArgValueLog
	mapMark2Value[ArgMarkLogConsole] = ArgValueConsole
	mapMark2Value[ArgMarkRecurse] = ArgValueRecurse
	mapMark2Value[ArgMarkStable] = ArgValueStable
	mapMark2Value[ArgMarkUpdate] = ArgValueUpdate

	mapValue2Mark[ArgValueDouble] = ArgMarkDouble
	mapValue2Mark[ArgValueIgnore] = ArgMarkIgnore
	mapValue2Mark[ArgValueLog] = ArgMarkLogFile
	mapValue2Mark[ArgValueConsole] = ArgMarkLogConsole
	mapValue2Mark[ArgValueRecurse] = ArgMarkRecurse
	mapValue2Mark[ArgValueStable] = ArgMarkStable
	mapValue2Mark[ArgValueUpdate] = ArgMarkUpdate
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
