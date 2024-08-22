package main

import (
	"fmt"
	"neocitiesCli/api"
	"os"
)

const VERSION = "0.0.1"

func main() {
	cmd, err := setFlags()
	if err != nil {
		fmt.Println("error setting flags")
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	config, err := getDefaultConfig()
	if err != nil {
		config = DEFAILT_CONFIG.Configs[0]
	}

	if config.Domain == "" || config.APIKey == "" {
		fmt.Println("no login found")
		config, err = login()
		if err != nil {
			fmt.Println("unable to login")
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		} else {
			fmt.Println("login successful")
		}
		err = writeConfig(config)
		if err != nil {
			fmt.Println("unable to write config")
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		// login runs if not logged in,
		// if cmd is to login then its not needed
		// so we can just return. EZ
		if cmd.function == "config" && cmd.args[0] == "login" {
			return
		}
	}

	conn := api.NewConnection(config)

	err = execute(conn, cmd)
	if err != nil {
		fmt.Println("error executing command")
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
