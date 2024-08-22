package main

import (
	"encoding/json"
	"fmt"
	"neocitiesCli/api"
	"os"
	"strings"
)

func login() (api.Config, error) {
	config := api.Config{}

	domain, isSubdomain, err := getDomainPrompt()
	if err != nil {
		return config, err
	}
	config.Domain = domain
	config.IsSubdomain = isSubdomain

	api_key, err := getAPIkeyPrompt(domain)
	if err != nil {
		return config, err
	}
	config.APIKey = api_key

	return config, nil
}

func getDomainPrompt() (string, bool, error) {
	var isSubdomain bool = false
	var err error
	isSubdomainString := ""
	domain := os.Getenv("NEOCITIES_SITENAME")
	for domain == "" {
		prompt := "Enter your Neocities Domain: "
		fmt.Print(prompt)
		_, err = fmt.Scanln(&domain)
		if err != nil && err.Error() != "unexpected newline" {
			fmt.Print("error: \n", err)
			return "", isSubdomain, err
		}
	}
	for isSubdomainString != "y" && isSubdomainString != "n" {
		fmt.Printf("Is this a subdomain (i.e. neocities.org/%s) or a custom domain?", domain)
		fmt.Print(" [y/n]\ny:  subdomain\nn:  custom domain\n")
		_, err = fmt.Scanln(&isSubdomainString)
		if err != nil && err.Error() != "unexpected newline" {
			return "", isSubdomain, err
		}
		if isSubdomainString == "Y" || isSubdomainString == "s" {
			isSubdomainString = "y"
		}
		if isSubdomainString == "N" || isSubdomainString == "c" {
			isSubdomainString = "n"
		}
		if isSubdomainString == "y" {
			// if !strings.HasPrefix(domain, "neocities.org/") {
			// 	domain = "neocities.org/" + domain
			// }
			isSubdomain = true
		}
	}
	return domain, isSubdomain, nil
}

func getAPIkeyPrompt(domain string) (string, error) {
	var err error

	pass := os.Getenv("NEOCITIES_PASSWORD")
	for pass == "" {
		prompt := "Enter Password for " + domain + ": "
		fmt.Print(prompt)
		_, err = fmt.Scanln(&pass)
		if err != nil && err.Error() != "unexpected newline" {
			return "", err
		}
	}

	api_key, err := api.GetAPIkey(domain, pass)
	if err != nil {
		return "", err
	}
	fmt.Printf("Your API Key is: %s\n", api_key)

	if api_key == "" {
		err = fmt.Errorf("could not find API Key")
		return api_key, err
	}
	return api_key, nil
}

func resetPrompt() bool {
	var input string
	fmt.Print("Are you sure you want to reset the configuration file? [y/n]\n")
	_, err := fmt.Scanln(&input)
	if err != nil && err.Error() != "unexpected newline" {
		return false
	}
	if input == "y" || input == "Y" {
		return true
	}
	return false
}

func removeGitIgnore(files []string) ([]string, error) {
	var new_files []string
	var isIgnored bool = false
	ignore, err := os.ReadFile(".gitignore")
	if err != nil {
		return []string{}, err
	}
	ignore_lines := strings.Split(string(ignore), "\n")
	for _, file := range files {
		isIgnored = false
		for _, ignore := range ignore_lines {

			if strings.Contains(ignore, "*") {
				// start
				if strings.HasPrefix(ignore, "*") {
					ignore = strings.TrimPrefix(ignore, "*")
					if strings.HasSuffix(file, ignore) {
						isIgnored = true
						break
					}
				}
				// end
				if strings.HasSuffix(ignore, "*") {
					ignore = strings.TrimSuffix(ignore, "*")
					if strings.HasPrefix(file, ignore) {
						isIgnored = true
						break
					}
				}
				// middle
				first, last, _ := strings.Cut(ignore, "*")
				if strings.HasPrefix(file, first) {
					if strings.HasSuffix(file, last) {
						isIgnored = true
						break
					}
				}
				continue
			}
			if ignore == file {
				isIgnored = true
				break
			}
		}
		if !isIgnored {
			new_files = append(new_files, file)
		}
	}
	return new_files, nil
}

func upload(conn *api.Connection, paths []string) error {
	var errs []error
	var files []api.UploadFile

	for _, path := range paths {
		text, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		file := api.UploadFile{
			Name: path,
			File: text,
		}
		files = append(files, file)
	}
	conn.Upload(files)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

var CONFIG_PATH string = ""

var DEFAILT_CONFIG ConfigFile = ConfigFile{
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

func getDefaultConfig() (api.Config, error) {
	file, err := readConfig()

	if err != nil || len(file.Configs) == 0 {
		return DEFAILT_CONFIG.Configs[0], err
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

func readConfig() (ConfigFile, error) {
	cfgs := ConfigFile{}

	path, err := getConfigPath()
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

func writeConfig(cfg api.Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	cfgFile, err := readConfig()
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

func getConfigPath() (string, error) {
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
