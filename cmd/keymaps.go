package cmd

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up           key.Binding
	Down         key.Binding
	Select       key.Binding
	Back         key.Binding
	Download     key.Binding
	RefreshEmail key.Binding
	Attchments   key.Binding
	Quit         key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up,
		k.Down,
		k.Select,
		k.Back,
		k.RefreshEmail,
		k.Download,
		k.Attchments,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.Up,
		k.Down,
		k.Select,
		k.Back,
		k.RefreshEmail,
		k.Download,
		k.Attchments,
		k.Quit,
	}}
}

func (m model) GetKeyMap() keyMap {
	var keys = keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}

	if m.selected != nil {
		keys.Back = key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "back"),
		)
		keys.Download = key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Download"),
		)

		if !m.attView {
			keys.Attchments = key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "Attchments"),
			)
		}

	} else {
		keys.Select = key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "select"),
		)
		keys.RefreshEmail = key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "New Address"),
		)
	}

	return keys
}
