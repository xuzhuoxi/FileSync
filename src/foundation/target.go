package foundation

type Target struct {
	Name     string   `yaml:"name"`
	Mode     string   `yaml:"mode"`
	Src      string   `yaml:"src"`
	Tar      string   `yaml:"tar"`
	Includes []string `yaml:"include"`
	Excludes []string `yaml:"exclude"`
	Case     bool     `yaml:"case"`
	Args     string   `yaml:"args"`
}

func (t *Target) GetMode() Mode {
	return GetMode(t.Mode)
}

func (t *Target) MatchTarget(value string) bool {
	if len(t.Includes) == 0 && len(t.Excludes) == 0 {
		return true
	}
	if t.checkInWildcard(t.Excludes, value) {
		return false
	}
	if t.checkInWildcard(t.Includes, value) {
		return true
	}
	return false
}

func (t *Target) MatchParam(param TargetParam) bool {
	return IncludeParam(t.Args, param)
}

func (t *Target) checkInWildcard(wildcards []string, value string) bool {
	if len(wildcards) == 0 {
		return false
	}
	for index := range wildcards {
		if MatchWildcard(value, wildcards[index], t.Case) {
			return true
		}
	}
	return false
}

type TargetGroup struct {
	Name    string   `yaml:"name"`
	Targets []string `yaml:"targets"`
}
