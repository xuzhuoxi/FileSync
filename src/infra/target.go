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
		ArgMarks: argMarks}
}

type RuntimeTarget struct {
	Name     string
	Mode     RuntimeMode
	SrcArr   []string
	Tar      string
	Includes []string
	Excludes []string
	Case     bool
	ArgMarks ArgMark
}

func (t *RuntimeTarget) CheckNameFitting(filename string) bool {
	if len(t.Includes) == 0 && len(t.Excludes) == 0 {
		return true
	}
	if t.checkInWildcard(t.Excludes, filename) {
		return false
	}
	if t.checkInWildcard(t.Includes, filename) {
		return true
	}
	return false
}

func (t *RuntimeTarget) MatchParam(param ArgMark) bool {
	return t.ArgMarks.MatchArg(param)
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
