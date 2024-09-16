package ui

import (
	"fmt"
	"neocitiesCli/api"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedColor   = lipgloss.Color("#ed49a3")
	unselectedColor = lipgloss.Color("#9349ed")
	mainStyle       = lipgloss.NewStyle().
			Width(40).
			Height(8).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(unselectedColor).
			Padding(1)
	headerStyle = lipgloss.NewStyle().PaddingLeft(16)
	// headerStyle = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(unselectedColor).Width(36)
	buttonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(unselectedColor).MarginLeft(14)
	// Padding(1)
)

type loginModel struct {
	Config   api.Config
	LoggedIn bool
	loginErr error

	username string
	pw       string
	// apiKey    string
	userInput textinput.Model
	keyInput  textinput.Model
	// useAPIKey toggle

	backgroundColor string
	accentColor     string
	textColor       string

	width  int
	height int

	// style lipgloss.Style

	focused       bool
	index         int
	awaitingLogin bool
	invalidLogin  bool
	staleView     bool
	viewCache     string

	indexMap []string
}

func newLogin() loginModel {
	// s := lipgloss.NewStyle().
	// 	Width(40).
	// 	Height(8).
	// 	Border(lipgloss.RoundedBorder()).
	// 	BorderForeground(lipgloss.Color("242")).
	// 	Padding(1)
	ut := textinput.New()
	ut.Focus()
	// ut.Prompt = ">"
	ut.Placeholder = "username"
	ut.CharLimit = 32 // neocities usernames are <= 32 chars

	kt := textinput.New()
	// kt.Prompt = ">"
	kt.EchoMode = textinput.EchoPassword
	kt.Placeholder = "password"
	kt.Validate = func(text string) error { // neocities passwords are >= 5 chars
		if len(text) < 5 {
			return fmt.Errorf("password must be at least 5 characters")
		}
		return nil
	}

	var indexList = []string{
		0: "username",
		1: "password",
		// 2: "apikeyToggle",
		2: "login",
	}

	return loginModel{
		Config:   api.Config{},
		LoggedIn: false,

		username: "",
		pw:       "",
		// apiKey:   "",

		userInput: ut,
		keyInput:  kt,
		// useAPIKey: newToggle(),

		backgroundColor: "",
		accentColor:     "",
		textColor:       "",

		width:  40,
		height: 8,

		// style: s,

		focused:       false,
		index:         0,
		awaitingLogin: false,
		invalidLogin:  false,
		staleView:     true,
		viewCache:     "",

		indexMap: indexList,
	}
}

func (m loginModel) Init() tea.Cmd {
	return nil
}

func (m loginModel) Update(msg tea.Msg) (loginModel, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case loginMsg:
		if msg.err != nil {
			m.invalidLogin = true
			m.loginErr = msg.err
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.index > 0 {
				m.index--
				m.staleView = true
			}
			switch m.indexMap[m.index] {

			case "username":
				m.userInput.Focus()
				m.keyInput.Blur()
			case "password":
				m.keyInput.Focus()
				// m.useAPIKey.Blur()
				// case "apikeyToggle":
				// 	m.useAPIKey.Focus()
				// m.loginButton.Blur()
			}

		case "down", "tab":
			if m.index < len(m.indexMap)-1 {
				m.index++
				m.staleView = true
			}
			switch m.indexMap[m.index] {

			case "username":
				m.userInput.Focus()
				m.keyInput.Blur()
			case "password":
				m.keyInput.Focus()
				m.userInput.Blur()
				// m.useAPIKey.Blur()
			// case "apikeyToggle":
			// 	m.useAPIKey.Focus()
			// 	m.keyInput.Blur()
			case "login":
				m.keyInput.Blur()
				// m.useAPIKey.Blur()
				// m.loginButton.Focus()
			}

		case "enter", " ":
			if m.invalidLogin && m.awaitingLogin {
				m.invalidLogin = false
				m.awaitingLogin = false
				break
			}
			m.awaitingLogin = true
			m.username = m.userInput.Value()
			m.pw = m.keyInput.Value()
			cmds = append(cmds, teaLogin(m.username, m.pw))
		}
	}
	m.userInput, cmd = m.userInput.Update(msg)
	cmds = append(cmds, cmd)

	m.keyInput, cmd = m.keyInput.Update(msg)
	cmds = append(cmds, cmd)

	// m.useAPIKey, cmd = m.useAPIKey.Update(msg)
	// cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m loginModel) View() string {
	// if !m.staleView {
	// 	return m.viewCache
	// }
	if m.invalidLogin {
		return fmt.Sprintf("error logging in:\n%s\nPress enter to try again", m.loginErr)
	}
	if m.awaitingLogin {
		return m.awaitingLoginView()
	}
	var view strings.Builder

	view.WriteString(headerStyle.Render("Login"))
	view.WriteString("\nUsername:\n")
	if m.indexMap[m.index] == "username" {
		view.WriteString(inputStyle.BorderForeground(selectedColor).Render(m.userInput.View()))
	} else {
		view.WriteString(inputStyle.Render(m.userInput.View()))
	}
	view.WriteString("\nPassword:\n")
	if m.indexMap[m.index] == "password" {
		view.WriteString(inputStyle.BorderForeground(selectedColor).Render(m.keyInput.View()))
	} else {
		view.WriteString(inputStyle.Render(m.keyInput.View()))
	}
	view.WriteRune('\n')
	// view.WriteString(m.useAPIKey.View() + " Use API Key")
	view.WriteString(m.loginButton())
	return mainStyle.Render(view.String())
}

func (m *loginModel) Focus() {
	m.focused = true
}

func (m *loginModel) Blur() {
	m.focused = false
}

func (m loginModel) loginButton() string {
	if m.indexMap[m.index] == "login" {
		return buttonStyle.BorderForeground(lipgloss.Color("#ed49a3")).Render("[Login]")
	}
	return buttonStyle.Render(" Login ")
}

func (m loginModel) awaitingLoginView() string {
	// TODO: show loading spinner
	return headerStyle.Render("Logging in...")
}
