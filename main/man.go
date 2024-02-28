package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
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
	//RpcListenAddress string `yaml:"rpc_server_listen"`
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

	go http_server.SetupRouter(config.GinHttp)

	//err = rpc_server.DefaultRpcServer.Run(config.RpcListenAddress)

	return err
}
