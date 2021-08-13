package infra

import "strings"

type RuntimeMode int

const (
	ModeNone RuntimeMode = iota
	ModeClear
	ModeCopy
	ModeDelete
	ModeMove
	ModeSync
)

const (
	ModeClearValue  = "clear"
	ModeCopyValue   = "copy"
	ModeDeleteValue = "delete"
	ModeMoveValue   = "move"
	ModeSyncValue   = "sync"
)

var (
	modMap  = make(map[RuntimeMode]string)
	modMap2 = make(map[string]RuntimeMode)
)

func init() {
	modMap[ModeClear] = ModeClearValue
	modMap[ModeCopy] = ModeCopyValue
	modMap[ModeDelete] = ModeDeleteValue
	modMap[ModeMove] = ModeMoveValue
	modMap[ModeSync] = ModeSyncValue

	modMap2[ModeClearValue] = ModeClear
	modMap2[ModeCopyValue] = ModeCopy
	modMap2[ModeDeleteValue] = ModeDelete
	modMap2[ModeMoveValue] = ModeMove
	modMap2[ModeSyncValue] = ModeSync
}

// 通过字符串，查找模式
func GetMode(value string) RuntimeMode {
	if value == "" {
		return ModeNone
	}
	value = strings.ToLower(value)
	if mod, ok := modMap2[value]; ok {
		return mod
	}
	return ModeNone
}

// 取模式的字符串值
func GetModeValue(mod RuntimeMode) string {
	if mod == ModeNone {
		return ""
	}
	if value, ok := modMap[mod]; ok {
		return value
	}
	return ""
}
