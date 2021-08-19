package infra

func NewRuntimeTarget(target ConfigTarget) *RuntimeTarget {
	mode := target.GetMode()
	srcArr := target.GetSrcArr()
	tar := target.Tar
	includes := target.GetIncludeArr()
	excludes := target.GetExcludeArr()
	argMarks := target.GetArgsMark()
	return &RuntimeTarget{
		Name:     target.Name,
		Mode:     mode,
		SrcArr:   srcArr,
		Tar:      tar,
		Includes: includes,
		Excludes: excludes,
		Case:     target.Case,
		ArgsMark: argMarks}
}

type RuntimeTarget struct {
	Name     string
	Mode     RuntimeMode
	SrcArr   []string
	Tar      string
	Includes []string
	Excludes []string
	Case     bool
	ArgsMark ArgMark
}

func (t *RuntimeTarget) HasIncludeLimit() bool {
	return len(t.Includes) != 0
}

func (t *RuntimeTarget) HasExcludeLimit() bool {
	return len(t.Excludes) != 0
}

func (t *RuntimeTarget) CheckNameFitting(filename string) bool {
	if !t.HasExcludeLimit() && !t.HasIncludeLimit() {
		return true
	}
	if t.HasExcludeLimit() && t.checkInWildcard(t.Excludes, filename) {
		return false
	}
	if t.HasIncludeLimit() && !t.checkInWildcard(t.Includes, filename) {
		return false
	}
	return true
}

func (t *RuntimeTarget) MatchParam(param ArgMark) bool {
	return t.ArgsMark.MatchArg(param)
}

func (t *RuntimeTarget) checkInWildcard(wildcards []string, value string) bool {
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
