package ui

import (
	"neocitiesCli/api"

	tea "github.com/charmbracelet/bubbletea"
)

type loginMsg struct {
	err error
	cfg api.Config
}

func teaLogin(username, password string) tea.Cmd {

	cfg := api.Config{
		Domain:      username,
		IsSubdomain: false,
	}

	return func() tea.Msg {
		key, err := api.GetAPIkey(username, password)
		if err != nil {
			return loginMsg{
				err: err,
			}
		}
		cfg.APIKey = key
		return loginMsg{
			cfg: cfg,
		}
	}
}
