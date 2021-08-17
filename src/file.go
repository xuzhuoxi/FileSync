package main

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/filex"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

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
