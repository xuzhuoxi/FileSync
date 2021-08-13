package infra

import "strings"

const separator = ","

type ConfigTarget struct {
	Name    string `yaml:"name"`
	Mode    string `yaml:"mode"`
	Src     string `yaml:"src"`
	Tar     string `yaml:"tar"`
	Include string `yaml:"include"`
	Exclude string `yaml:"exclude"`
	Case    bool   `yaml:"case"`
	Args    string `yaml:"args"`
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
