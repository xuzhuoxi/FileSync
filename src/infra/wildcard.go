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
	// 长度不足
	if len(value) < len(wildcard) {
		return false
	}
	// 绝对相等
	if value == wildcard {
		return true
	}
	// 非通配符且比较失败
	if strings.Index(wildcard, Wildcard) < 0 && value != wildcard {
		return false
	}
	// 通配符"*"
	if wildcard == Wildcard {
		return true
	}

	// 大小写处理
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
