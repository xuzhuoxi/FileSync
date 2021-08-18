package infra

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/filex"
	"strings"
)

const (
	WildcardSep = ","
)

type ConfigTarget struct {
	Name    string `yaml:"name"`
	Mode    string `yaml:"mode"`
	Src     string `yaml:"src"`
	Tar     string `yaml:"tar"`
	Include string `yaml:"include"`
	Exclude string `yaml:"exclude"`
	Args    string `yaml:"args"`
	Case    bool   `yaml:"case"`
}

func (ct ConfigTarget) String() string {
	return fmt.Sprintf("ConfigTarget[Name='%s',Mode='%s',Src='%s',Tar='%s',Include='%s',Exclude='%s',Args='%s',Case='%v']",
		ct.Name, ct.Mode, ct.Src, ct.Tar, ct.Include, ct.Exclude, ct.Args, ct.Case)
}

func (ct ConfigTarget) GetMode() RuntimeMode {
	return GetMode(ct.Mode)
}

func (ct ConfigTarget) GetSrcArr() []string {
	if ct.Src == "" {
		return nil
	}
	if !strings.Contains(ct.Src, filex.PathListSeparatorStr) {
		return []string{ct.Src}
	}
	return strings.Split(ct.Src, filex.PathListSeparatorStr)
}

func (ct ConfigTarget) GetIncludeArr() []string {
	return ct.splitWildcard(ct.Include)
}

func (ct ConfigTarget) GetExcludeArr() []string {
	return ct.splitWildcard(ct.Exclude)
}

func (ct ConfigTarget) splitWildcard(value string) []string {
	if value == "" {
		return nil
	}
	if !strings.Contains(value, WildcardSep) {
		return []string{value}
	}
	return strings.Split(value, WildcardSep)
}

func (ct ConfigTarget) GetArgsMark() ArgMark {
	return ValuesToMarks(ct.Args)
}

type ConfigGroup struct {
	Name    string `yaml:"name"`
	Targets string `yaml:"targets"`
}

type Config struct {
	Main    string         `yaml:"main"`
	Groups  []ConfigGroup  `yaml:"groups"`
	Targets []ConfigTarget `yaml:"targets"`
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
