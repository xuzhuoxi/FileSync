package clear

import "github.com/xuzhuoxi/FileSync/src/infra"

func Clear(src []string, include string, exclude string, mathCase bool, args string) {
	target := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Case: mathCase, Args: args}
	runtimeTarget := infra.NewRuntimeTarget(target)
	ClearWithTarget(runtimeTarget)
}

func ClearWithTarget(target *infra.RuntimeTarget) {

}
