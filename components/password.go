package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PasswordModel struct {
	textinput textinput.Model
	show      bool
	callbacks []func(password string) tea.Cmd
}

func NewPasswordModel() PasswordModel {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.CharLimit = 50
	ti.EchoMode = textinput.EchoPassword

	return PasswordModel{
		textinput: ti,
	}
}

func (m PasswordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PasswordModel) Update(msg tea.Msg) (PasswordModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.show {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.show = false
				m.FlushCallbacks()
				cmds = append(cmds, func() tea.Msg {
					return BlurPasswordDialogMsg{}
				}, updateLayoutCmd())
				m.textinput.Reset()
			case "enter":
				m.show = false
				cmds = append(cmds, func() tea.Msg {
					return BlurPasswordDialogMsg{}
				}, m.CallbackInBatch(), updateLayoutCmd())
				m.FlushCallbacks()
				m.textinput.Reset()
			}
		case passwordInputStartMsg:
			m.PushCallback(msg.callback)
		}

		m.textinput, cmd = m.textinput.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case passwordInputStartMsg:
		m.show = true
		m.PushCallback(msg.callback)
		cmds = append(cmds, func() tea.Msg {
			return FocusPasswordDialogMsg{}
		}, updateLayoutCmd())
	}

	return m, tea.Batch(cmds...)
}

func (m PasswordModel) View() string {
	if !m.show {
		return ""
	}

	dialog := lipgloss.JoinVertical(lipgloss.Center,
		"Enter your password to proceed",
		m.textinput.View(),
		"[Enter] OK  [Esc] Cancel",
	)

	return dialogStyle.Render(dialog)
}

func (m PasswordModel) GetSize() (x, y int) {
	if !m.show {
		return 0, 0
	}
	fw, fh := dialogStyle.GetFrameSize()
	return m.textinput.Width + fw, 2 + fh
}

func (m PasswordModel) CallbackInBatch() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.callbacks))
	for _, c := range m.callbacks {
		cmds = append(cmds, c(m.textinput.Value()))
	}

	return tea.Batch(cmds...)
}

func (m *PasswordModel) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m *PasswordModel) Blur() {
	m.textinput.Blur()
}

func (m *PasswordModel) PushCallback(callback func(password string) tea.Cmd) {
	m.callbacks = append(m.callbacks, callback)
}

func (m *PasswordModel) FlushCallbacks() {
	m.callbacks = []func(password string) tea.Cmd{}
}
