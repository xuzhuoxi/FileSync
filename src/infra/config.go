package infra

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/slicex"
	"strings"
)

const (
	PathListSeparatorStr = filex.PathListSeparatorStr
	DirSeparator         = filex.UnixSeparator
	NameSeq              = ","
)

type ConfigTask struct {
	Name    string `yaml:"name"`    // 任务名称，用于唯一标识配置
	Mode    string `yaml:"mode"`    // 任务模式，RuntimeMode相应的string值
	Src     string `yaml:"src"`     // 任务源路径，多个用";"分隔
	Tar     string `yaml:"tar"`     // 任务目标路径，单一目标
	Include string `yaml:"include"` // 处理包含的文件名通配符
	Exclude string `yaml:"exclude"` // 处理排除的文件名通配符
	Args    string `yaml:"args"`    // 任务管理参数
}

func (ct ConfigTask) String() string {
	return fmt.Sprintf("Task{Name=%s,Mode=%s}", ct.Name, ct.Mode)
}

func (ct ConfigTask) ToFullString() string {
	return fmt.Sprintf("ConfigTask[Name='%s',Mode='%s',Src='%s',Tar='%s',Include='%s',Exclude='%s',Args='%s']",
		ct.Name, ct.Mode, ct.Src, ct.Tar, ct.Include, ct.Exclude, ct.Args)
}

func (ct ConfigTask) ToShortString() string {
	return fmt.Sprintf("ConfigTask{Name='%s',Mode='%s',Include='%s',Exclude='%s',Args='%s'}",
		ct.Name, ct.Mode, ct.Include, ct.Exclude, ct.Args)
}

func (ct ConfigTask) ToPathString() string {
	return fmt.Sprintf("ConfigTask{Name='%s',Src='%s',Tar='%s'}",
		ct.Name, ct.Src, ct.Tar)
}

func (ct ConfigTask) GetMode() RuntimeMode {
	return GetMode(ct.Mode)
}

func (ct ConfigTask) GetSrcArr() []SrcInfo {
	if ct.Src == "" {
		return nil
	}
	if !strings.Contains(ct.Src, PathListSeparatorStr) {
		return []SrcInfo{NewSrcInfo(ct.Src)}
	}
	srcArr := strings.Split(ct.Src, PathListSeparatorStr)
	rs := make([]SrcInfo, len(srcArr))
	for index := range srcArr {
		rs[index] = NewSrcInfo(srcArr[index])
	}
	return rs
}

func (ct ConfigTask) GetIncludeArr() (fws []Wildcard, dws []Wildcard, err error) {
	return ParseWildcards(ct.Include)
}

func (ct ConfigTask) GetExcludeArr() (fws []Wildcard, dws []Wildcard, err error) {
	return ParseWildcards(ct.Exclude)
}

func (ct ConfigTask) CheckSelf() (err error) {
	m, errMode := checkMode(ct.Mode)
	if nil != errMode {
		err = errMode
		return
	}
	err = checkSrc(ct.Src)
	if nil != err {
		return
	}
	if m == ModeCopy || m == ModeMove || m == ModeSync {
		err = checkTar(ct.Tar)
		if nil != err {
			return
		}
	}
	err = checkArgs(ct.Args, GetSupportArgs(m))
	return
}

func (ct ConfigTask) GetArgsMark() ArgMark {
	return ValuesToMarks(ct.Args)
}

type ConfigGroup struct {
	Name  string `yaml:"name"`
	Tasks string `yaml:"tasks"`
}

func (o ConfigGroup) GetTasks() (tasks []string) {
	str := strings.TrimSpace(o.Tasks)
	if len(str) == 0 {
		return nil
	}
	return strings.Split(str, NameSeq)
}

type Config struct {
	RelativeRoot string        `yaml:"root"`      // 相对路径的根目录
	Main         string        `yaml:"main"`      // 任务默认入口
	Groups       []ConfigGroup `yaml:"groups"`    // 任务组
	Sequences    []ConfigGroup `yaml:"sequences"` // 预制任务组合
	Tasks        []ConfigTask  `yaml:"tasks"`     // 任务列表
}

// 配置主任务列表
func (c *Config) MainTasks() []ConfigTask {
	return c.GetMainTasks(c.Main)
}

// 取任务列表
// main不区分Group与Task
func (c *Config) GetMainTasks(main string) []ConfigTask {
	if main == "" {
		return nil
	}
	if task, ok := c.GetTask(main); ok {
		return []ConfigTask{task}
	}
	if tasks, ok := c.GetTasksFromGroup(main); ok {
		return tasks
	}
	return nil
}

func (c *Config) GetTasksFromGroup(groupName string) (tasks []ConfigTask, ok bool) {
	if len(c.Groups) == 0 {
		return nil, false
	}
	for index := range c.Groups {
		if c.Groups[index].Name == groupName {
			ts := c.Groups[index].GetTasks()
			for _, name := range ts {
				if task, ok := c.GetTask(name); ok { // 普通任务
					tasks = append(tasks, task)
					continue
				}
				if seqTasks, ok := c.GetTasksFromSequence(name); ok { // 预制任务
					tasks = append(tasks, seqTasks...)
					continue
				}
			}
			return tasks, true
		}
	}
	return nil, false
}

func (c *Config) GetTasksFromSequence(seqName string) (tasks []ConfigTask, ok bool) {
	if len(c.Sequences) == 0 {
		return nil, false
	}
	for index := range c.Sequences {
		if c.Sequences[index].Name == seqName {
			return c.GetTasks(c.Sequences[index].GetTasks()), true
		}
	}
	return nil, false
}

// 取任务列表
func (c *Config) GetTasks(taskNames []string) []ConfigTask {
	if nil == c.Tasks || len(c.Tasks) == 0 {
		return nil
	}
	if nil == taskNames || len(taskNames) == 0 {
		return nil
	}
	var rs []ConfigTask
	for index := range taskNames {
		if task, ok := c.GetTask(taskNames[index]); ok {
			rs = append(rs, task)
		}
	}
	return rs
}

// 取任务
func (c *Config) GetTask(taskName string) (task ConfigTask, ok bool) {
	if nil == c.Tasks || len(c.Tasks) == 0 {
		return ConfigTask{}, false
	}
	for index := range c.Tasks {
		if c.Tasks[index].Name == taskName {
			return c.Tasks[index], true
		}
	}
	return ConfigTask{}, false
}

//-------------------------

func checkMode(modeValue string) (mode RuntimeMode, err error) {
	if m, ok := CheckModeValue(modeValue); ok {
		return m, nil
	}
	return ModeNone, errors.New(fmt.Sprintf("Undefined module:%v", modeValue))
}

func checkSrc(srcValue string) (err error) {
	if "" == srcValue || "" == strings.TrimSpace(srcValue) {
		return errors.New(fmt.Sprintf("Src Empty! "))
	}
	if !strings.Contains(srcValue, PathListSeparatorStr) {
		return nil
	}
	srcArr := strings.Split(srcValue, PathListSeparatorStr)
	for index := range srcArr {
		if "" == srcArr[index] || "" == strings.TrimSpace(srcArr[index]) {
			return errors.New(fmt.Sprintf("Src[%d] Empty! ", index))
		}
	}
	return
}

func checkTar(tarValue string) (err error) {
	if "" == tarValue || "" == strings.TrimSpace(tarValue) {
		return errors.New(fmt.Sprintf("Tar Empty! "))
	}
	if strings.Contains(tarValue, PathListSeparatorStr) {
		return errors.New(fmt.Sprintf("Tar does not support multi paths! "))
	}
	return nil
}

func checkArgs(value string, supports []string) (err error) {
	if "" == value {
		return
	}
	if len(supports) == 0 {
		return errors.New(fmt.Sprintf("Unsupport Args:'%s'", value))
	}
	args := SplitArgs(value)
	for index := range args {
		if !slicex.ContainsString(supports, args[index]) {
			return errors.New(fmt.Sprintf("Unsupport Arg[%d]:'%s'", index, args[index]))
		}
	}
	return nil
}
