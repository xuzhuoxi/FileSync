package foundation

type Config struct {
	Main    string        `yaml:"main"`
	Groups  []TargetGroup `yaml:"groups"`
	Targets []*Target     `yaml:"targets"`
}

// 配置主任务列表
func (c *Config) MainTargets() []*Target {
	return c.GetMainTargets(c.Main)
}

// 取任务列表
// main不区分Group与Target
func (c *Config) GetMainTargets(main string) []*Target {
	if main == "" {
		return nil
	}
	target := c.GetTarget(main)
	if nil != target {
		return []*Target{target}
	}
	if c.Groups == nil || len(c.Groups) == 0 {
		return nil
	}
	for index := range c.Groups {
		if c.Groups[index].Name == main {
			return c.GetTargets(c.Groups[index].Targets)
		}
	}
	return nil
}

// 取任务列表
func (c *Config) GetTargets(targetNames []string) []*Target {
	if nil == c.Targets || len(c.Targets) == 0 {
		return nil
	}
	if nil == targetNames || len(targetNames) == 0 {
		return nil
	}
	var rs []*Target
	for index := range targetNames {
		target := c.GetTarget(targetNames[index])
		if nil == target {
			continue
		}
		rs = append(rs, target)
	}
	return rs
}

// 取任务
func (c *Config) GetTarget(targetName string) *Target {
	if nil == c.Targets || len(c.Targets) == 0 {
		return nil
	}
	for index := range c.Targets {
		if c.Targets[index].Name == targetName {
			return c.Targets[index]
		}
	}
	return nil
}
