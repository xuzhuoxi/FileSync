package main

import (
	"flag"
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module"
	_ "github.com/xuzhuoxi/FileSync/src/module/clear"
	_ "github.com/xuzhuoxi/FileSync/src/module/copy"
	_ "github.com/xuzhuoxi/FileSync/src/module/delete"
	_ "github.com/xuzhuoxi/FileSync/src/module/move"
	_ "github.com/xuzhuoxi/FileSync/src/module/sync"
)

func main() {
	cfgTargets, err := parseFlags()
	if nil != err {
		infra.Logger.Errorln("Line:18", err)
		return
	}
	infra.Logger.Infoln(fmt.Sprintf("[main] Target number=%d", len(cfgTargets)))
	execTargets(cfgTargets)
}

func execTargets(cfgTargets []infra.ConfigTarget) {
	if len(cfgTargets) == 0 {
		return
	}
	for index := range cfgTargets {
		execTarget(cfgTargets[index])
	}
}

func execTarget(cfgTarget infra.ConfigTarget) {
	executor := module.GetExecutor(cfgTarget.GetMode())
	infra.Logger.Infoln(fmt.Sprintf("[main] Target=%v", cfgTarget))
	executor.ExecConfigTarget(cfgTarget)
}

// 处理命令行参数
// 得到任务配置
func parseFlags() (targets []infra.ConfigTarget, err error) {
	file := flag.String("file", "", "Running mode! ")
	main := flag.String("main", "", "Main! ")

	mode := flag.String("module", "", "Running mode! ")
	src := flag.String("src", "", "Src path or Src paths! ")
	tar := flag.String("tar", "", "Tar path! ")
	include := flag.String("include", "", "Include settings! ")
	exclude := flag.String("exclude", "", "exclude settings! ")
	args := flag.String("args", "", "Running args! ")
	wildcardCase := flag.Bool("case", true, "Whether include settings and exclude settings are case sensitive! ")

	flag.Parse()

	if *file != "" {
		targets, err = loadTargets(*file, *main)
	} else {
		target, errTarget := genTarget(fmt.Sprintf("Cmd.%s", *mode), *mode, *src, *tar, *include, *exclude, *wildcardCase, *args)
		if nil != errTarget {
			err = errTarget
		} else {
			targets = []infra.ConfigTarget{target}
		}
	}
	return
}
