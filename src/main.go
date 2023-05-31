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
	cfgTasks, err := parseFlags()
	if nil != err {
		infra.Logger.Errorln("Flags Error(Line:18):", err)
		return
	}
	infra.Logger.Println()
	infra.Logger.Println("********************新命令分割线********************")
	infra.Logger.Infoln(fmt.Sprintf("[main] Tasks[%d]:%v", len(cfgTasks), cfgTasks))
	infra.Logger.Infoln(fmt.Sprintf("[main] RunningRelativeRoot='%s'", infra.RunningDir))
	execTasks(cfgTasks)
}

func execTasks(cfgTasks []infra.ConfigTask) {
	if len(cfgTasks) == 0 {
		return
	}
	for index := range cfgTasks {
		execTask(cfgTasks[index])
	}
}

func execTask(cfgTask infra.ConfigTask) {
	err := cfgTask.CheckSelf()
	if nil != err {
		infra.Logger.Errorln(fmt.Sprintf("[main] Task=%v Err=%v", cfgTask.Name, err))
		return
	}
	executor := module.GetExecutor(cfgTask.GetMode())
	infra.Logger.Println(fmt.Sprintf("--------------------任务(%s)--------------------", cfgTask.Name))
	infra.Logger.Infoln(fmt.Sprintf("[main] TaskInfo=%v", cfgTask.ToShortString()))
	infra.Logger.Infoln(fmt.Sprintf("[main] TaskPath=%v", cfgTask.ToPathString()))
	executor.ExecConfigTask(cfgTask)
}

// 处理命令行参数
// 得到任务配置
func parseFlags() (tasks []infra.ConfigTask, err error) {
	file := flag.String("file", "", "Running mode! ")
	main := flag.String("main", "", "Main! ")

	mode := flag.String("mode", "", "Running mode! ")
	src := flag.String("src", "", "Src path or Src paths! ")
	tar := flag.String("tar", "", "Tar path! ")
	include := flag.String("include", "", "Include settings! ")
	exclude := flag.String("exclude", "", "exclude settings! ")
	args := flag.String("args", "", "Running args! ")

	flag.Parse()

	if *file != "" {
		tasks, err = loadConfigTasks(*file, *main)
	} else {
		task := genTask(fmt.Sprintf("Cmd.%s", *mode), *mode, *src, *tar, *include, *exclude, *args)
		tasks = []infra.ConfigTask{task}
	}
	return
}

func genTask(name, mode, src, tar, include, exclude string, args string) (task infra.ConfigTask) {
	mode = strings.ToLower(mode)
	return infra.ConfigTask{
		Name: name, Mode: mode, Src: src, Tar: tar, Include: include, Exclude: exclude, Args: args}
}

func loadConfigTasks(filePath string, main string) (tasks []infra.ConfigTask, err error) {
	cfgPath := filePath
	if !filex.IsFile(filePath) {
		cfgPath = filex.Combine(infra.RunningDir, filePath)
	}
	config := &infra.Config{}
	err = loadConfigFile(cfgPath, config)
	if nil != err {
		return
	}
	if "" != config.RelativeRoot {
		infra.SetRunningDir(filex.FormatPath(config.RelativeRoot))
	} else {
		upDir, _ := filex.GetUpDir(cfgPath)
		infra.SetRunningDir(upDir)
	}
	if "" == main {
		tasks = config.MainTasks()
	} else {
		tasks = config.GetMainTasks(main)
	}
	if len(tasks) == 0 {
		err = errors.New(fmt.Sprintf("No tasks with name '%s'", main))
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
