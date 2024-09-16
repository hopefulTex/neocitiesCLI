package ui

import (
	"neocitiesCli/api"

	tea "github.com/charmbracelet/bubbletea"
)

// file browser using api - list && https - get

type directoryBrowser struct {
	root    string
	dirs    []string
	items   api.ListItem
	focused bool
}

func newDirectoryBrowser() directoryBrowser {
	return directoryBrowser{}
}

func (m directoryBrowser) Update(msg tea.Msg) (directoryBrowser, tea.Cmd) {
	var cmds []tea.Cmd
	// var cmd tea.Cmd
	if !m.focused {
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m directoryBrowser) help() string {
	return `
	help
	`
}
