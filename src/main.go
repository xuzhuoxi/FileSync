package main

import (
	"flag"
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
)

func main() {
	targets, err := parseFlags()
	if nil != err {
		infra.Logger.Error(err)
		return
	}
	infra.Logger.Infoln(fmt.Sprintf("[main] Target len=%d", len(targets)))
	execTargets(targets)
}

func execTargets(targets []infra.ConfigTarget) {
	if len(targets) == 0 {
		return
	}
	for index := range targets {
		execTarget(targets[index])
	}
}

func execTarget(target infra.ConfigTarget) {

}

func parseFlags() (targets []infra.ConfigTarget, err error) {
	file := flag.String("mode", "", "Running Mode! ")
	if *file == "" {
		main := flag.String("main", "", "Main! ")
		targets, err = loadTargets(*file, *main)
	} else {
		mode := flag.String("mode", "", "Running Mode! ")
		src := flag.String("src", "", "Use Languages! ")
		tar := flag.String("tar", "", "Use Fields! ")
		include := flag.String("include", "", "Output Files! ")
		exclude := flag.String("exclude", "", "Source Redefine! ")
		wildcardCase := flag.Bool("case", true, "Source Redefine! ")
		args := flag.String("args", "", "Target Redefine! ")
		target, errTarget := genTarget("Main", *mode, *src, *tar, *include, *exclude, *wildcardCase, *args)
		if nil != errTarget {
			err = errTarget
		} else {
			targets = []infra.ConfigTarget{target}
		}
	}
	return
}
