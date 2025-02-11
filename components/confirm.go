package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmModel struct {
	msg      string
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

		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case showDialogMsg:
		m.msg = msg.msg
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
		"[Enter] OK  [Esc] Cancel",
	)

	return dialogStyle.Render(dialog)
}

func (m ConfirmModel) GetSize() (x, y int) {
	if !m.show {
		return 0, 0
	}
	fw, fh := dialogStyle.GetFrameSize()
	return len(m.msg) + fw, 2 + fh
}
