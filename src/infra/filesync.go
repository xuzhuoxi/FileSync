package infra

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/mathx"
	"github.com/xuzhuoxi/infra-go/osxu"
)

const (
	ApplicationName    = "FileSync"
	ApplicationVersion = "1.0.6"
)

var (
	Logger     = logx.NewLogger()
	RunningDir = osxu.GetRunningDir()
)

func init() {
	Logger = GenLogger(MarkLogFile | MarkLogPrint)
}

func SetRunningDir(dir string) {
	RunningDir = dir
}

func GenLogger(mark ArgMark) logx.ILogger {
	logger := logx.NewLogger()
	if mark.MatchArg(MarkLogPrint) {
		logger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
	}
	if mark.MatchArg(MarkLogFile) {
		logger.SetConfig(logx.LogConfig{Type: logx.TypeRollingFile, Level: logx.LevelAll,
			FileDir: osxu.GetRunningDir(), FileName: ApplicationName, FileExtName: ".log", MaxSize: 4 * mathx.MB})
	}
	return logger
}
