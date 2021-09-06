package infra

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/slicex"
	"strings"
)

const (
	PathListSeparatorStr = filex.PathListSeparatorStr
	DirSeparator         = filex.UnixSeparator
)

type ConfigTarget struct {
	Name    string `yaml:"name"`    // 任务名称，用于唯一标识配置
	Mode    string `yaml:"mode"`    // 任务模式，RuntimeMode相应的string值
	Src     string `yaml:"src"`     // 任务源路径，多个用";"分隔
	Tar     string `yaml:"tar"`     // 任务目标路径，单一目标
	Include string `yaml:"include"` // 处理包含的文件名通配符
	Exclude string `yaml:"exclude"` // 处理排除的文件名通配符
	Args    string `yaml:"args"`    // 任务管理参数
}

func (ct ConfigTarget) String() string {
	return fmt.Sprintf("ConfigTarget[Name='%s',Mode='%s',Src='%s',Tar='%s',Include='%s',Exclude='%s',Args='%s']",
		ct.Name, ct.Mode, ct.Src, ct.Tar, ct.Include, ct.Exclude, ct.Args)
}

func (ct ConfigTarget) GetMode() RuntimeMode {
	return GetMode(ct.Mode)
}

func (ct ConfigTarget) GetSrcArr() []SrcInfo {
	if ct.Src == "" {
		return nil
	}
	if !strings.Contains(ct.Src, PathListSeparatorStr) {
		return []SrcInfo{NewSrcInfo(ct.Src)}
	}
	srcArr := strings.Split(ct.Src, PathListSeparatorStr)
	rs := make([]SrcInfo, len(srcArr))
	for index := range srcArr {
		rs[index] = NewSrcInfo(srcArr[index])
	}
	return rs
}

func (ct ConfigTarget) GetIncludeArr() (fws []Wildcard, dws []Wildcard, err error) {
	return ParseWildcards(ct.Include)
}

func (ct ConfigTarget) GetExcludeArr() (fws []Wildcard, dws []Wildcard, err error) {
	return ParseWildcards(ct.Exclude)
}

func (ct ConfigTarget) CheckTarget() (err error) {
	m, errMode := checkMode(ct.Mode)
	if nil != errMode {
		err = errMode
		return
	}
	err = checkSrc(ct.Src)
	if nil != err {
		return
	}
	if m == ModeCopy || m == ModeMove || m == ModeSync {
		err = checkTar(ct.Tar)
		if nil != err {
			return
		}
	}
	err = checkArgs(ct.Args, GetSupportArgs(m))
	return
}

func (ct ConfigTarget) GetArgsMark() ArgMark {
	return ValuesToMarks(ct.Args)
}

type ConfigGroup struct {
	Name    string `yaml:"name"`
	Targets string `yaml:"targets"`
}

type Config struct {
	RelativeRoot string         `yaml:"root"` // 相对路径的根目录
	Main         string         `yaml:"main"`
	Groups       []ConfigGroup  `yaml:"groups"`
	Targets      []ConfigTarget `yaml:"targets"`
}

// 配置主任务列表
func (c *Config) MainTargets() []ConfigTarget {
	return c.GetMainTargets(c.Main)
}

// 取任务列表
// main不区分Group与Target
func (c *Config) GetMainTargets(main string) []ConfigTarget {
	if main == "" {
		return nil
	}
	if target, ok := c.GetTarget(main); ok {
		return []ConfigTarget{target}
	}
	if c.Groups == nil || len(c.Groups) == 0 {
		return nil
	}
	for index := range c.Groups {
		if c.Groups[index].Name == main {
			targetNames := strings.Split(c.Groups[index].Targets, ",")
			return c.GetTargets(targetNames)
		}
	}
	return nil
}

// 取任务列表
func (c *Config) GetTargets(targetNames []string) []ConfigTarget {
	if nil == c.Targets || len(c.Targets) == 0 {
		return nil
	}
	if nil == targetNames || len(targetNames) == 0 {
		return nil
	}
	var rs []ConfigTarget
	for index := range targetNames {
		if target, ok := c.GetTarget(targetNames[index]); ok {
			rs = append(rs, target)
		}
	}
	return rs
}

// 取任务
func (c *Config) GetTarget(targetName string) (target ConfigTarget, ok bool) {
	if nil == c.Targets || len(c.Targets) == 0 {
		return ConfigTarget{}, false
	}
	for index := range c.Targets {
		if c.Targets[index].Name == targetName {
			return c.Targets[index], true
		}
	}
	return ConfigTarget{}, false
}

//-------------------------

func checkMode(modeValue string) (mode RuntimeMode, err error) {
	if m, ok := CheckModeValue(modeValue); ok {
		return m, nil
	}
	return ModeNone, errors.New(fmt.Sprintf("Undefined module:%v", modeValue))
}

func checkSrc(srcValue string) (err error) {
	if "" == srcValue || "" == strings.TrimSpace(srcValue) {
		return errors.New(fmt.Sprintf("Src Empty! "))
	}
	if !strings.Contains(srcValue, PathListSeparatorStr) {
		return nil
	}
	srcArr := strings.Split(srcValue, PathListSeparatorStr)
	for index := range srcArr {
		if "" == srcArr[index] || "" == strings.TrimSpace(srcArr[index]) {
			return errors.New(fmt.Sprintf("Src[%d] Empty! ", index))
		}
	}
	return
}

func checkTar(tarValue string) (err error) {
	if "" == tarValue || "" == strings.TrimSpace(tarValue) {
		return errors.New(fmt.Sprintf("Tar Empty! "))
	}
	if strings.Contains(tarValue, PathListSeparatorStr) {
		return errors.New(fmt.Sprintf("Tar does not support multi paths! "))
	}
	return nil
}

func checkArgs(value string, supports []string) (err error) {
	if "" == value {
		return
	}
	if len(supports) == 0 {
		return errors.New(fmt.Sprintf("Unsupport Args:'%s'", value))
	}
	args := SplitArgs(value)
	for index := range args {
		if !slicex.ContainsString(supports, args[index]) {
			return errors.New(fmt.Sprintf("Unsupport Arg[%d]:'%s'", index, args[index]))
		}
	}
	return nil
}
