package ui

import (
	"fmt"
	"neocitiesCli/api"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var icon string = "🥰"

var statusStyle lipgloss.Style = lipgloss.NewStyle().
	Background(lipgloss.Color(ACCENT_COLOR)).
	Foreground(lipgloss.Color(TEXT_COLOR)).Width(38)

func (m model) statusView(info api.SiteInfo) string {
	var view strings.Builder

	view.WriteString(icon)
	view.WriteRune(' ')
	view.WriteString(info.Name)
	view.WriteString("\tSite Visits: ")
	view.WriteString(fmt.Sprintf("%d", info.Hits))
	view.WriteString("\tTags: ")
	view.WriteString(strings.Join(info.Tags, ", "))
	return statusStyle.Render(view.String())
}
