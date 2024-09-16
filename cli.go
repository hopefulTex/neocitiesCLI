package main

import (
	"flag"
	"fmt"
	"io/fs"
	"neocitiesCli/api"
	"neocitiesCli/ui"
	"neocitiesCli/util"
	"os"
	"path/filepath"
)

type command struct {
	function     string
	args         []string
	useGitIgnore bool
}

func execute(conn *api.Connection, config api.Config, cmd command) error {
	var err error

	switch cmd.function {
	case "push":
		var files []string
		err = filepath.WalkDir(cmd.args[0], func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
		if cmd.useGitIgnore {
			files, err = util.RemoveGitIgnore(files)
			if err != nil {
				return err
			}
		}
		util.Upload(conn, files)
	case "upload":
		var errs []error
		var files []string
		for _, file := range cmd.args {
			_, err := os.Stat(file)
			if err != nil {
				errs = append(errs, err)
			} else {
				files = append(files, file)
			}
		}
		if len(errs) > 0 {
			errorList := ""
			for _, err := range errs {
				errorList += err.Error() + "\n"
			}
			return fmt.Errorf("multiple errors:\n%s", errorList)
		}
		if cmd.useGitIgnore {
			files, err = util.RemoveGitIgnore(cmd.args)
			if err != nil {
				return err
			}
		}
		util.Upload(conn, files)
	case "delete":
		files := cmd.args
		if cmd.useGitIgnore {
			files, err = util.RemoveGitIgnore(cmd.args)
			if err != nil {
				return err
			}
		}
		err = conn.Delete(files)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted %d files\n", len(files))
	case "list":
		args := ""
		if len(cmd.args) > 0 {
			args = cmd.args[0]
		}
		items, err := conn.List(args)
		if err != nil {
			return err
		}
		for _, item := range items {
			fmt.Println(item.View())
		}
	case "info":
		info, err := conn.Info(cmd.args[0])
		if err != nil {
			return fmt.Errorf("error: %s\nEnsure URL is valid", err)
		}
		fmt.Println(ui.InfoView(info))
	case "config":
		// change current account
		// list accounts
		// login
		// reset (scary)

		if len(cmd.args) == 0 {
			fmt.Println("no arguments provided")
			fmt.Print(CONFIG_HELP_STRING)
			return fmt.Errorf("invalid arguments")
		}
		switch cmd.args[0] {
		case "list":
			configList, err := util.ReadConfig()
			if err != nil {
				return err
			}
			for _, cfg := range configList.Configs {
				fmt.Printf("- %s\n", cfg.Domain)
			}
			return nil
		case "set":
			if len(cmd.args) < 2 {
				fmt.Println("Usage: neocities config set [domain]")
				return fmt.Errorf(("invalid arguments"))
			}
			configList, err := util.ReadConfig()
			if err != nil {
				return err
			}
			defaultSet := false
			var changed []int = []int{}
			for i, cfg := range configList.Configs {
				if cfg.IsDefault {
					configList.Configs[i].IsDefault = false
					changed = append(changed, i)
				}
				if !defaultSet && cfg.Domain == cmd.args[1] {
					configList.Configs[i].IsDefault = true
					defaultSet = true
					changed = append(changed, i)
				}
			}
			if !defaultSet {
				return fmt.Errorf("domain %s not found", cmd.args[1])
			}
			for _, i := range changed {
				err = util.WriteConfig(configList.Configs[i])
				if err != nil {
					return err
				}
			}
			return nil
		case "login":
			cfg, err := util.Login()
			if err != nil {
				return err
			}
			err = util.WriteConfig(cfg)
			return err
		case "reset":
			if len(cmd.args) == 0 || cmd.args[0] != "--force" {
				resetting := util.ResetPrompt()
				if !resetting {
					return fmt.Errorf("reset aborted")
				}
			}
			path, err := util.GetConfigPath()
			if err != nil {
				return err
			}
			err = os.Remove(path)
			if err != nil {
				return err
			}
			util.WriteConfig(util.DEFAULT_CONFIG.Configs[0])
		default:
			fmt.Printf("invalid subcommand: %s\n", cmd.args[0])
			fmt.Print(CONFIG_HELP_STRING)
			return fmt.Errorf("invalid subcommand")
		}

	case "tui":
		err := ui.Run(*conn, config)

		if err != nil {
			return err
		}
	case "version":
		fmt.Println(VERSION)
		return nil
	}
	return err
}

func setFlags() (command, error) {
	var cmd command
	args := os.Args
	if len(args) == 1 {
		cmd.function = "tui"
		return cmd, nil
	}
	switch args[1] {
	case "push":
		if len(args) < 3 {
			fmt.Println("Usage: neocities push directory")
			return cmd, fmt.Errorf("invalid arguments")
		}
		cmd = command{
			function: "push",
			args:     args[2:],
		}

	case "upload":
		if len(args) < 3 {
			fmt.Println("Usage: neocities upload [files]")
			return cmd, fmt.Errorf("invalid arguments")
		}
		cmd = command{
			function: "upload",
			args:     args[2:],
		}

	case "delete":
		if len(args) < 3 {
			fmt.Println("Usage: neocities delete [files]")
			return cmd, fmt.Errorf("invalid arguments")
		}
		cmd = command{
			function: "delete",
			args:     args[2:],
		}

	case "list":
		cmd = command{
			function: "list",
			args:     args[2:],
		}

	case "info":
		sitename := ""
		if len(args) > 2 {
			sitename = args[2]
		}
		cmd = command{
			function: "info",
			args:     []string{sitename},
		}

	case "config":
		cmd = command{
			function: "config",
			args:     args[2:],
		}

	case "--version":
		cmd = command{
			function: "version",
			args:     nil,
		}
	default:
		fmt.Print(MAIN_HELP_STRING)
		return cmd, fmt.Errorf("invalid subcommand")
	}

	ignore := flag.Bool("no-gitignore", false, "ignore exclusions listed in .gitignore file")
	cmd.useGitIgnore = !*ignore
	flag.Parse()

	return cmd, nil
}
