package cmd

import (
	"fmt"
	"strings"

	"github.com/MarlonCorreia/ephemail/utils"
	"github.com/charmbracelet/lipgloss"
)

func (m model) helpView() string {
	return fmt.Sprintf(" %s\n", m.help.View(m.GetKeyMap()))
}

func (m model) headerView() string {
	title := titleStyle.Render(fmt.Sprintf("EphEmail - %s", highlightStyle.Render(m.emailClient.GetEmail())))
	line := strings.Repeat("─", utils.MaxInt(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", utils.MaxInt(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m model) listView() string {
	s := "\n"
	if len(m.emailClient.Messages) == 0 {
		s += fmt.Sprintf("\n %s%s%s\n\n", m.spinner.View(), " ", "Waiting Emails")
	} else {
		for i, msg := range m.emailClient.Messages {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			emailDetails := fmt.Sprintf("from: %s at:%s", msg.From, msg.Date)
			s += fmt.Sprintf("%s %s\n  %s\n", highlightStyle.Render(cursor), msg.Subject, fadedTextStyle(emailDetails))
		}
	}
	return s
}

func (m model) attchmentsView() string {
	s := "\n"
	if len(m.selected.Content.Attachments) == 0 {
		return "\n No Attchments to this Email"
	}

	for i, att := range m.selected.Content.Attachments {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		attDetails := fmt.Sprintf("type: %s size: %d", att.ContentType, att.Size)
		s += fmt.Sprintf("%s %s\n  %s\n", highlightStyle.Render(cursor), att.FileName, fadedTextStyle(attDetails))
	}

	return s

}

func (m model) messageContentView() string {
	content := m.selected.DisplayCompleteEmail()
	style := lipgloss.NewStyle().Width(lipgloss.Width(m.headerView()))
	return style.Render(content)
}
