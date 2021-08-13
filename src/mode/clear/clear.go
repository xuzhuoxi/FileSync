package clear

import "github.com/xuzhuoxi/FileSync/src/infra"

func Clear(basePath string, src string, include string, exclude string, mathCase bool, args string) {
	target := infra.ConfigTarget{Name: "Clear", Mode: infra.ModeClearValue, Src: src,
		Include: include, Exclude: exclude, Case: mathCase, Args: args}
	runtimeTarget := infra.NewRuntimeTarget(target)
	ClearWithTarget(basePath, runtimeTarget)
}

func ClearWithTarget(basePath string, target *infra.RuntimeTarget) {

}
