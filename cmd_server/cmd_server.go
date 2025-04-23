package cmd_server

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"root/http_server"
	"root/manager"
)

func Work() error {

	setupLog()
	useConsoleLogIfDebug()

	dir, _ := os.Getwd()

	config, err := loadConfig(dir)
	if err != nil {
		return fmt.Errorf("load config, %w", err)
	}

	manager.SetupGroupChannels(config.GroupList)

	useJsonLogIfRelease(dir)

	log.Info().Msgf("http server start at [%s]", config.HttpListenAddress)

	return http_server.StartRouter(config.HttpListenAddress)
}
