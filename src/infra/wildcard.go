package infra

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/stringx"
	"strings"
)

const (
	WildcardCharStr    = "*"
	WildcardSepStr     = ","
	WildcardTypeSepStr = ";"
	WildcardTypeFile   = "file:"
	WildcardTypeDir    = "dir:"
)

type Wildcard string

// 是名称
func (w Wildcard) IsName() bool {
	return !w.IsWildcard()
}

// 是通配符
func (w Wildcard) IsWildcard() bool {
	return strings.Contains(string(w), WildcardCharStr)
}

// 匹配判断
func (w Wildcard) Match(value string) bool {
	// 长度不足
	wildcard := string(w)
	if len(value) < len(wildcard) {
		return false
	}
	// 绝对相等
	if value == wildcard {
		return true
	}
	// 非通配符且比较失败
	if strings.Index(wildcard, WildcardCharStr) < 0 && value != wildcard {
		return false
	}
	// 通配符"*"
	if wildcard == WildcardCharStr {
		return true
	}

	ws := strings.Split(wildcard, WildcardCharStr)
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

// 匹配判断
// 大小写相关
func (w Wildcard) MatchWithCase(value string) bool {
	return w.Match(value)
}

// 解释通配符信息
// value格式：file:*.jpg,*.png;dir:*,a*b
func ParseWildcards(value string) (fileWildcards []Wildcard, dirWildcards []Wildcard, err error) {
	if value == "" {
		return nil, nil, nil
	}
	if !strings.Contains(value, WildcardTypeSepStr) {
		return ParseTypeWildcards(value)
	}
	values := strings.Split(value, WildcardTypeSepStr)
	for _, val := range values {
		fws, dws, err := ParseTypeWildcards(val)
		if nil != err {
			return nil, nil, err
		}
		fileWildcards = append(fileWildcards, fws...)
		dirWildcards = append(dirWildcards, dws...)
	}
	return
}

// 解释通配符
// value格式: file:*.jpg,*.png 或 dir:*.jpg,*.png
func ParseTypeWildcards(value string) (fileWildcards []Wildcard, dirWildcards []Wildcard, err error) {
	if value == "" {
		return nil, nil, nil
	}
	fileVal, fileOk := checkStart(value, WildcardTypeFile)
	if fileOk {
		fileWildcards = ParseValueWildcards(fileVal)
		return
	}
	dirVal, dirOk := checkStart(value, WildcardTypeDir)
	if dirOk {
		dirWildcards = ParseValueWildcards(dirVal)
		return
	}
	return nil, nil, errors.New(fmt.Sprintf("Type wildcars parse errror:%s", value))
}

// 解释通配符
// value格式: *.jpg,*.png
func ParseValueWildcards(value string) []Wildcard {
	if value == "" {
		return nil
	}

	if !strings.Contains(value, WildcardSepStr) {
		return []Wildcard{Wildcard(value)}
	}
	values := strings.Split(value, WildcardSepStr)
	rs := make([]Wildcard, len(values))
	for index := range values {
		rs[index] = Wildcard(values[index])
	}
	return rs
}

// 检查是否以start开关，并返回除start外内容
func checkStart(value string, start string) (listStr string, ok bool) {
	if "" == value {
		return "", false
	}
	index := strings.Index(value, start)
	if 0 != index {
		return "", false
	}
	listStr = value[len(start):]
	return listStr, true
}
