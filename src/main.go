package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/FileSync/src/module"
	"github.com/xuzhuoxi/infra-go/filex"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
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
	err := cfgTarget.CheckTarget()
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[main] Target=%v Err=%v", cfgTarget.Name, err))
		return
	}
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
		target := genTarget(fmt.Sprintf("Cmd.%s", *mode), *mode, *src, *tar, *include, *exclude, *wildcardCase, *args)
		targets = []infra.ConfigTarget{target}
	}
	return
}

func genTarget(name, mode, src, tar, include, exclude string, wildcardCase bool,
	args string) (target infra.ConfigTarget) {
	mode = strings.ToLower(mode)
	return infra.ConfigTarget{
		Name: name, Mode: mode, Src: src, Tar: tar, Include: include, Exclude: exclude, Case: wildcardCase, Args: args}
}

func loadTargets(relativeFilePath string, main string) (targets []infra.ConfigTarget, err error) {
	absPath := filex.Combine(infra.RunningDir, relativeFilePath)
	config := &infra.Config{}
	err = loadConfigFile(absPath, config)
	if nil != err {
		return
	}
	if "" == main {
		return config.MainTargets(), nil
	}
	targets = config.GetMainTargets(main)
	if len(targets) == 0 {
		err = errors.New(fmt.Sprintf("No targets with name '%s'", main))
	}
	return
}

func loadConfigFile(configPath string, dataRef interface{}) error {
	bs, err := ioutil.ReadFile(configPath)
	if nil != err {
		return err
	}
	err = yaml.Unmarshal(bs, dataRef)
	if nil != err {
		return err
	}
	return nil
}
