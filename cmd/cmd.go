package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/MarlonCorreia/ephemail/internal/clipb"
	email "github.com/MarlonCorreia/ephemail/internal/email"
	"github.com/MarlonCorreia/ephemail/utils"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	emailClient   email.EmailModel
	cursor        int
	selected      *email.Message
	attView       bool
	viewport      viewport.Model
	viewPortReady bool
	stopwatch     stopwatch.Model
	spinner       spinner.Model
	error         string
}

func initialModel() model {
	error := ""
	client := email.EmailModel{}
	err := client.BuildNewEmail()
	if err != nil {
		error = newEmailAdressErr
	}
	newSpinner := spinner.New()
	newSpinner.Spinner = spinner.Line
	newSpinner.Style = highlightStyle

	return model{
		emailClient:   client,
		selected:      nil,
		viewport:      viewport.New(100, 20),
		viewPortReady: false,
		stopwatch:     stopwatch.NewWithInterval(time.Millisecond),
		spinner:       newSpinner,
		error:         error,
		attView:       false,
	}
}

func (m model) Init() tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	cmd = m.spinner.Tick
	cmds = append(cmds, cmd)

	cmd = m.stopwatch.Init()
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		helpHeight := lipgloss.Height(m.helpView())
		verticalMarginHeight := headerHeight + footerHeight + helpHeight

		if !m.viewPortReady {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight + 1
			m.viewPortReady = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.selected == nil || m.attView {
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "down", "j":
			if m.selected == nil {
				if m.cursor < len(m.emailClient.Messages)-1 {
					m.cursor++
				}
			} else if m.attView {
				if m.cursor < len(m.selected.Content.Attachments)-1 {
					m.cursor++
				}
			}

		case "enter", " ":
			if m.cursor >= 0 || m.cursor < len(m.emailClient.Messages) {
				m.selected = m.emailClient.Messages[m.cursor]
				m.viewport.SetContent(m.messageContentView())
			}

		case "b":
			if m.attView {
				m.cursor = 0
				m.attView = false
			} else {
				m.selected = nil
				m.viewport.SetContent("")
			}

		case "c":
			err := clipb.SendToClipBoard(m.emailClient.GetEmail())
			if err != nil {
				m.error = copyEmailAddressErr
			} else if m.error == copyEmailAddressErr {
				m.error = ""
			}

		case "n":
			err := m.emailClient.BuildNewEmail()
			if err != nil {
				m.error = newEmailAdressErr
			} else {
				m.error = ""
			}
			m.emailClient.Messages = []*email.Message{}

		case "d":
			if m.selected != nil {
				if m.attView {
					if len(m.selected.Content.Attachments) > m.cursor {
						att := m.selected.Content.Attachments[m.cursor]
						err := m.emailClient.DownloadAttachment(m.selected, att)
						if err != nil {
							m.error = downloadAttchmentErr
						}
					}
				} else {
					htmlContent := m.selected.Content.HtmlBody
					err := utils.WriteFile("email.html", htmlContent)
					if err != nil {
						m.error = writeEmailToFileErr
					}
				}
			}

		case "a":
			if m.selected != nil {
				m.cursor = 0
				m.attView = true
			}
		}
	}

	if m.stopwatch.Elapsed().Seconds() >= 5 {
		cmd = m.stopwatch.Reset()
		cmds = append(cmds, cmd)

		if m.error != newEmailAdressErr {
			err := m.emailClient.UpdateEmailMessages()
			if err != nil {
				m.error = fetchNewEmailsErr
			} else {
				m.error = ""
			}
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := ""

	if m.error != "" {
		s += errorStateStyle(m.error)
	} else if m.selected == nil {
		s += m.listView()
	} else {
		if m.attView {
			s += m.attchmentsView()
		} else {
			s += m.viewport.View()
		}
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), s, m.footerView(), m.helpView())
}

func InitCmd() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
