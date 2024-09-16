package ui

import (
	"neocitiesCli/api"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ACCENT_COLOR     = "#986AFC"
	BACKGROUND_COLOR = "#000000"
	TEXT_COLOR       = "#FFFFFF"
)

type model struct {
	directory    directoryBrowser
	login        loginModel
	config       api.Config
	conn         api.Connection
	info         api.SiteInfo
	invalidLogin bool
	isLoggedIn   bool
}

func newModel(conn api.Connection, config api.Config) model {
	l := newLogin()
	isLoggedIn := true
	siteInfo := api.SiteInfo{}
	if config.Domain == "" || config.APIKey == "" {
		l.Focus()
		isLoggedIn = false
	} else {
		info, err := conn.Info("")
		if err != nil {
			isLoggedIn = false
		}
		siteInfo = info
	}

	return model{
		directory:    newDirectoryBrowser(),
		login:        l,
		conn:         conn,
		config:       config,
		invalidLogin: false,
		isLoggedIn:   isLoggedIn,
		info:         siteInfo,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.login.LoggedIn {
		m.isLoggedIn = true
		m.config = m.login.Config
		m.login.Blur()
	}

	switch msg := msg.(type) {
	case loginMsg:
		if msg.err == nil && msg.cfg.APIKey != "" {
			m.config = msg.cfg
			m.isLoggedIn = true
			m.login.Blur()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	// m.directory, cmd = m.directory.Update(msg)
	// cmds = append(cmds, cmd)

	m.login, cmd = m.login.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var view strings.Builder

	if !m.isLoggedIn {
		view.WriteString(m.login.View())
		return view.String()
	}

	view.WriteString(m.statusView(m.info))
	// view.WriteString(InfoView(m.info))
	return view.String()
}

func Run(conn api.Connection, config api.Config) error {
	m := newModel(conn, config)
	p := tea.NewProgram(m)

	_, err := p.Run()
	return err
}
