package infra

import "strings"

func NewRuntimeTarget(target ConfigTarget) *RuntimeTarget {
	mode := GetMode(target.Mode)
	srcArr := strings.Split(target.Src, separator)
	tarArr := strings.Split(target.Tar, separator)
	includes := strings.Split(target.Include, separator)
	excludes := strings.Split(target.Exclude, separator)
	argMarks := ValuesToMarks(target.Args)
	return &RuntimeTarget{
		Name:     target.Name,
		Mode:     mode,
		SrcArr:   srcArr,
		TarArr:   tarArr,
		Includes: includes,
		Excludes: excludes,
		Case:     target.Case,
		ArgMarks: argMarks}
}

type RuntimeTarget struct {
	Name     string
	Mode     RuntimeMode
	SrcArr   []string
	TarArr   []string
	Includes []string
	Excludes []string
	Case     bool
	ArgMarks ParamMark
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

func (t *RuntimeTarget) MatchParam(param ParamMark) bool {
	return t.ArgMarks.MatchParam(param)
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
