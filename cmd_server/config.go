package cmd_server

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	GinHttp   string   `yaml:"http_server_listen"`
	GroupList []string `yaml:"group_list"`
}

func loadConfig(dir string) (*Config, error) {
	var config = new(Config)

	data, err2 := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err2 != nil {
		return nil, fmt.Errorf("read config file err: %w", err2)
	}

	println("----")
	println(string(data))
	println("----")

	err := yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file err: %w", err)
	}

	return config, nil
}
