package ui

import tea "github.com/charmbracelet/bubbletea"

type toggle struct {
	on      bool
	focused bool
}

func newToggle() toggle {
	return toggle{
		on:      false,
		focused: false,
	}
}

func (m toggle) Update(msg tea.Msg) (toggle, tea.Cmd) {
	var cmd tea.Cmd
	if !m.focused {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.on = !m.on
		}
	}
	return m, cmd
}

func (m toggle) View() string {
	var view string
	if m.on {
		view = "(x)"
	} else {
		view = "( )"
	}
	if m.focused {
		view = "[" + view + "]"
	}
	return view
}

func (m *toggle) Focus() {
	m.focused = true
}

func (m *toggle) Blur() {
	m.focused = false
}

func (m toggle) state() bool {
	return m.on
}
