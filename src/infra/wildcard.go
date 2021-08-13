package infra

import (
	"github.com/xuzhuoxi/infra-go/stringx"
	"strings"
)

const (
	Wildcard = "*"
)

// 检查是否为通配符
func CheckWildcard(value string) bool {
	return strings.Contains(value, Wildcard)
}

// 检查是否匹配
func MatchWildcard(value string, wildcard string, matchCase bool) bool {
	if value == wildcard || wildcard == Wildcard {
		return true
	}
	if len(value) < len(wildcard) {
		return false
	}

	if !matchCase {
		value = strings.ToLower(value)
		wildcard = strings.ToLower(wildcard)
	}

	ws := strings.Split(wildcard, Wildcard)
	checkValue := value
	for _, w := range ws {
		index := strings.Index(checkValue, w)
		if index == -1 {
			return false
		}
		checkValue = stringx.SubSuffix(checkValue, len(w))
	}
	return true
}
