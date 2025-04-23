package cmd_server

import (
	"fmt"
	"os"
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

	return nil
}
