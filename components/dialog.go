package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	DIALOG_WIDTH           = 60
	DIALOG_MAX_LINE_LENGTH = 50
)

var dialogStyle = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	Padding(1).
	Width(DIALOG_WIDTH).
	Align(lipgloss.Center)

func wrapText(input string, maxLineLength int) (string, int) {
	words := strings.Fields(input)
	var (
		result      []string
		currentLine string
	)

	maxlen := 0
	for _, word := range words {
		if len(currentLine)+len(word)+1 <= maxLineLength {
			if len(currentLine) > 0 {
				currentLine += " "
			}
			currentLine += word
		} else {
			result = append(result, currentLine)
			maxlen = max(maxlen, len(currentLine))
			currentLine = word
		}
	}

	result = append(result, currentLine)
	maxlen = max(maxlen, len(currentLine))

	return strings.Join(result, "\n"), maxlen
}

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
