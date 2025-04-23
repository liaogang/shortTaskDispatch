package cmd_server

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"root/extend/my_env"
)

func setupLog() {
	//启动时，打印日志到console,
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	zerolog.TimeFieldFormat = "01-02 15:04:05.000"
	zerolog.MessageFieldName = "msg"
	zerolog.LevelFieldName = "lvl"
	zerolog.TimestampFieldName = "tm"
}

func useConsoleLogIfDebug() {
	if my_env.DebugFlag {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"}
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}
}

func useJsonLogIfRelease(dir string) {

	if my_env.DebugFlag {
		return
	}

	//设置日志分隔插件, 把日志输入到文件
	log.Info().Str("lumberjack redirect path", "log/server.log").Send()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = "01-02 15:04:05.000"
	zerolog.MessageFieldName = "msg"
	zerolog.LevelFieldName = "" //don't write level field, only in json format

	var l = &lumberjack.Logger{
		Filename:   filepath.Join(dir, "/log/server.log"),
		MaxBackups: 100, // files
		MaxSize:    100, // megabytes
		MaxAge:     100, // days
		Compress:   true,
	}

	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(l)
}
