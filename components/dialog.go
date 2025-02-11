package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var dialogStyle = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	Padding(1).
	Width(60).
	Align(lipgloss.Center)

func updateLayoutCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateLayoutMsg{}
	}
}

func showDialogCmd(msg string, callback tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return showDialogMsg{
			msg:      "Are you sure?\n" + msg,
			callback: callback,
		}
	}
}
