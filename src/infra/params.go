package infra

import (
	"strings"
)

const (
	ParamValueDir     = "/d" //单向
	ParamValueDouble  = "/D" //双向
	ParamValueForce   = "/f" //强制(force)，若有重复或其它疑问时，不会询问用户，而强制复制
	ParamValueIgnore  = "/i" //忽略空目录
	ParamValueLog     = "/L" //记录日志
	ParamValueConsole = "/l" //控制台打印信息
	ParamValueRecurse = "/r" //递归
	ParamValueStable  = "/s" //保持文件目录结构
	ParamValueUpdate  = "/u" //若目标文件比源文件旧，更新目标文件
)

type ParamMark int

const (
	ParamMarDir ParamMark = 1 << iota
	ParamMarkDouble
	ParamMarkForce
	ParamMarkIgnore
	ParamMarkLog
	ParamMarkLogConsole
	ParamMarkRecurse
	ParamMarkStable
	ParamMarkUpdate
)

const DefaultParamMark = ParamMarkLog | ParamMarkLogConsole

var (
	mapValue2Mark = make(map[string]ParamMark)
	mapMark2Value = make(map[ParamMark]string)
)

func init() {
	mapMark2Value[ParamMarDir] = ParamValueDir
	mapMark2Value[ParamMarkDouble] = ParamValueDouble
	mapMark2Value[ParamMarkForce] = ParamValueForce
	mapMark2Value[ParamMarkIgnore] = ParamValueIgnore
	mapMark2Value[ParamMarkLog] = ParamValueLog
	mapMark2Value[ParamMarkLogConsole] = ParamValueConsole
	mapMark2Value[ParamMarkRecurse] = ParamValueRecurse
	mapMark2Value[ParamMarkStable] = ParamValueStable
	mapMark2Value[ParamMarkUpdate] = ParamValueUpdate

	mapValue2Mark[ParamValueDir] = ParamMarDir
	mapValue2Mark[ParamValueDouble] = ParamMarkDouble
	mapValue2Mark[ParamValueForce] = ParamMarkForce
	mapValue2Mark[ParamValueIgnore] = ParamMarkIgnore
	mapValue2Mark[ParamValueLog] = ParamMarkLog
	mapValue2Mark[ParamValueConsole] = ParamMarkLogConsole
	mapValue2Mark[ParamValueRecurse] = ParamMarkRecurse
	mapValue2Mark[ParamValueStable] = ParamMarkStable
	mapValue2Mark[ParamValueUpdate] = ParamMarkUpdate
}

// 检查参数码位的匹配情况
func (m ParamMark) MatchParam(param ParamMark) bool {
	return int(m)&int(param) > 0
}

// 字符串参数转换为码位
func ValuesToMarks(params string) ParamMark {
	values := splitParams(params)
	if nil == values {
		return 0
	}
	var rs ParamMark = 0
	for _, value := range values {
		rs = rs | value2Mark(value)
	}
	return ParamMark(rs)
}

// 拆分参数字符串为参数数组
func splitParams(params string) []string {
	if "" == params {
		return nil
	}
	rs := strings.Split(params, "/")
	if len(rs) == 0 {
		return nil
	}
	for index := range rs {
		rs[index] = "/" + rs[index]
	}
	return rs
}

func mark2Value(mark ParamMark) string {
	if v, ok := mapMark2Value[mark]; ok {
		return v
	}
	return ""
}

func value2Mark(value string) ParamMark {
	if v, ok := mapValue2Mark[value]; ok {
		return v
	}
	return 0
}
