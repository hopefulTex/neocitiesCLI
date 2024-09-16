package util

import (
	"encoding/json"
	"fmt"
	"neocitiesCli/api"
	"os"
	"strings"
)

var CONFIG_PATH string = ""

var DEFAULT_CONFIG ConfigFile = ConfigFile{
	Configs: []api.Config{
		{
			IsDefault:   true,
			Domain:      "",
			APIKey:      "",
			IsSubdomain: false,
		},
	},
}

type ConfigFile struct {
	Configs []api.Config `json:"configs"`
}

func GetDefaultConfig() (api.Config, error) {
	file, err := ReadConfig()

	if err != nil || len(file.Configs) == 0 {
		return DEFAULT_CONFIG.Configs[0], err
	}
	for _, cfg := range file.Configs {
		if cfg.IsDefault {

			if cfg.APIKey == "" {
				cfg.APIKey = os.Getenv("NEOCITIES_API_KEY")
			}
			if cfg.Domain == "" {
				cfg.Domain = os.Getenv("NEOCITIES_DOMAIN")
			}
			return cfg, nil
		}
	}

	if file.Configs[0].APIKey == "" {
		file.Configs[0].APIKey = os.Getenv("NEOCITIES_API_KEY")
	}
	if file.Configs[0].Domain == "" {
		file.Configs[0].Domain = os.Getenv("NEOCITIES_DOMAIN")
	}
	return file.Configs[0], nil
}

func ReadConfig() (ConfigFile, error) {
	cfgs := ConfigFile{}

	path, err := GetConfigPath()
	if err != nil {
		return cfgs, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return cfgs, err
	}
	err = json.Unmarshal(file, &cfgs)
	if err != nil {
		return cfgs, err
	}

	return cfgs, nil
}

func WriteConfig(cfg api.Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	cfgFile, err := ReadConfig()
	if err == nil {
		if len(cfgFile.Configs) == 0 {
			cfgFile.Configs = []api.Config{
				cfg,
			}
		}
		for i, block := range cfgFile.Configs {
			if block.Domain == cfg.Domain {
				cfgFile.Configs[i] = cfg
				break
			}
		}
	}
	cfgFile = ConfigFile{
		Configs: []api.Config{
			cfg,
		},
	}
	fmt.Printf("Writing config to %s\n", path)
	text, err := json.MarshalIndent(cfgFile, "", "  ")
	if err != nil {
		return err
	}

	dirPath := strings.TrimSuffix(path, "config.json")
	fmt.Printf("DirPath: %s\n", dirPath)
	_, err = os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("ReadDir Error: %s\n", err)
		if os.IsNotExist(err) {
			fmt.Printf("Creating directory %s\n", dirPath)
			err = os.Mkdir(dirPath, os.ModeDir)

		}
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = fmt.Fprintln(file, string(text))
	return err
}

func GetConfigPath() (string, error) {
	var path string
	if CONFIG_PATH == "" {
		path = os.Getenv("NEOCITIES_CONFIG_PATH")
		if path == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			path = home + "/.config/neocities/config.json"
		}
		CONFIG_PATH = path
	}

	return CONFIG_PATH, nil
}
