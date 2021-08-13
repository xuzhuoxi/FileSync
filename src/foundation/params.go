package foundation

import (
	"strings"
)

type TargetParam string

const (
	ParamDir        TargetParam = "/d" //单向
	ParamDouble     TargetParam = "/D" //双向
	ParamForce      TargetParam = "/f" //强制(force)，若有重复或其它疑问时，不会询问用户，而强制复制
	ParamIgnore     TargetParam = "/i" //忽略空目录
	ParamLog        TargetParam = "/L" //记录日志
	ParamLogConsole TargetParam = "/l" //控制台打印信息
	ParamRecurse    TargetParam = "/r" //递归
	ParamStable     TargetParam = "/s" //保持文件目录结构
	ParamUpdate     TargetParam = "/u" //若目标文件比源文件旧，更新目标文件
)

// 检查参数字符串中是否包含某个参数
func IncludeParam(params string, param TargetParam) bool {
	return strings.Contains(params, string(param))
}

// 拆分参数字符串为参数数组
func SplitParams(params string) []string {
	if "" == params {
		return nil
	}
	rs := strings.Split(params, "/")
	if len(rs) == 0 {
		return nil
	}
	for index, _ := range rs {
		rs[index] = "/" + rs[index]
	}
	return rs
}
