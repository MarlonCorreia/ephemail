package cmd

import "github.com/charmbracelet/lipgloss"

var (
	fadedTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#eeb902"))

	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)
