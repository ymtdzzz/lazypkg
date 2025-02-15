package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmModel struct {
	msg      string
	maxlen   int
	show     bool
	callback tea.Cmd
}

func NewConfirmModel() ConfirmModel {
	return ConfirmModel{}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	var cmds []tea.Cmd

	if m.show {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.show = false
				m.callback = nil
				cmds = append(cmds, func() tea.Msg {
					return BlurConfirmDialogMsg{}
				}, updateLayoutCmd())
			case "enter":
				m.show = false
				cmds = append(cmds, func() tea.Msg {
					return BlurConfirmDialogMsg{}
				}, m.callback, updateLayoutCmd())
				m.callback = nil
			}
		}
	}

	switch msg := msg.(type) {
	case showDialogMsg:
		m.msg, m.maxlen = wrapText(msg.msg, DIALOG_MAX_LINE_LENGTH)
		m.callback = msg.callback
		m.show = true
		cmds = append(cmds, func() tea.Msg {
			return FocusConfirmDialogMsg{}
		}, updateLayoutCmd())
	}

	return m, tea.Batch(cmds...)
}

func (m ConfirmModel) View() string {
	if !m.show {
		return ""
	}

	dialog := lipgloss.JoinVertical(lipgloss.Center,
		m.msg,
		"\n[Enter] OK  [Esc] Cancel",
	)

	return dialogStyle.Render(dialog)
}

func (m ConfirmModel) GetSize() (x, y int) {
	if !m.show {
		return 0, 0
	}
	fw, fh := dialogStyle.GetFrameSize()
	return DIALOG_WIDTH + fw, strings.Count(m.msg, "\n") + 3 + fh // 2 = line count + new line + button row
}
