package main

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/FileSync/src/infra"
	"github.com/xuzhuoxi/infra-go/slicex"
	"strings"
)

func genTarget(name, mode, src, tar, include, exclude string, wildcardCase bool,
	args string) (target infra.ConfigTarget, err error) {
	m, errMode := checkMode(mode)
	if nil != errMode {
		err = errMode
		return
	}
	src, err = checkSrc(src)
	if nil != err {
		return
	}
	tar, err = checkTar(tar)
	if nil != err {
		return
	}
	err = checkArgs(args, infra.GetSupportArgs(m))
	if nil != err {
		return
	}
	return infra.ConfigTarget{
		Name: name, Mode: mode, Src: src, Tar: tar, Include: include, Exclude: exclude, Case: wildcardCase, Args: args}, nil
}

func checkMode(modeValue string) (mode infra.RuntimeMode, err error) {
	if m, ok := infra.CheckModeValue(modeValue); ok {
		return m, nil
	}
	return infra.ModeNone, errors.New(fmt.Sprintf("Undefined mode:%v", modeValue))
}

func checkSrc(srcValue string) (src string, err error) {
	srcValue = strings.TrimSpace(srcValue)
	if "" == srcValue {
		return "", errors.New(fmt.Sprintf("Src Empty! "))
	}
	if !strings.Contains(srcValue, infra.PathSep) {
		return srcValue, nil
	}
	srcArr := strings.Split(srcValue, infra.PathSep)
	for index := range srcArr {
		if "" == srcArr[index] || "" == strings.TrimSpace(srcArr[index]) {
			return "", errors.New(fmt.Sprintf("Src[%d] Empty! ", index))
		}
	}
	return
}

func checkTar(tarValue string) (tar string, err error) {
	tarValue = strings.TrimSpace(tarValue)
	if "" == tarValue {
		return "", errors.New(fmt.Sprintf("Tar Empty! "))
	}
	if strings.Contains(tarValue, infra.PathSep) {
		return "", errors.New(fmt.Sprintf("Tar does not support multi paths! "))
	}
	return tarValue, nil
}

func checkArgs(value string, supports []string) (err error) {
	if "" == value {
		return
	}
	if len(supports) == 0 {
		return errors.New(fmt.Sprintf("Unsupport Args:'%s'", value))
	}
	args := infra.SplitArgs(value)
	for index := range args {
		if !slicex.ContainsString(supports, args[index]) {
			return errors.New(fmt.Sprintf("Unsupport Arg[%d]:'%s'", index, args[index]))
		}
	}
	return nil
}
