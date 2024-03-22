package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"root/extend/my_env"
	"root/http_server"
	"time"
)

func main() {

	//启动时，打印日志到console,
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	err := work()
	if err != nil {
		log.Error().AnErr("work", err).Send()
	}
}

type Config struct {
	GinHttp string `yaml:"http_server_listen"`
}

func work() error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	//load config
	data, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		return err
	}

	var config = new(Config)
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	log.Info().Interface("config", config).Send()

	//启动之后，如果是生产环境，日志输入到lumberjack文件进行分隔
	if my_env.ReleaseFlag {

		//设置日志分隔插件, 把日志输入到文件
		log.Info().Str("lumberjack redirect path", "log/server.log").Send()
		log.Info().Str("http server start at address", config.GinHttp).Send()

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zerolog.TimeFieldFormat = "01-02 15:04:05.000"
		zerolog.MessageFieldName = "msg"

		var l = &lumberjack.Logger{
			Filename:   filepath.Join(dir, "/log/server.log"),
			MaxBackups: 100, // files
			MaxSize:    100, // megabytes
			MaxAge:     100, // days
			Compress:   true,
		}

		var defaultLogger = log.Logger
		log.Logger = defaultLogger.Output(l)

		log.Info().Msg("----")
		log.Info().Str("http server start at address", config.GinHttp).Send()
		log.Info().Msg("----")
	} else {
		log.Info().Str("http server start at address", config.GinHttp).Send()
	}

	http_server.SetupRouter(config.GinHttp)

	return nil
}
