package components

import (
	"container/ring"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type LogWriter struct {
	mu   sync.Mutex
	logs *ring.Ring
}

func (w *LogWriter) Write(b []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logs.Value = string(b)
	w.logs = w.logs.Next()

	return len(b), nil
}

func (w *LogWriter) getLog() string {
	var sb strings.Builder
	w.logs.Do(func(p any) {
		if p != nil {
			sb.WriteString(p.(string))
		}
	})

	return sb.String()
}

type outputKeyMap struct {
	up   key.Binding
	down key.Binding
}

func newOutputKeyMap() outputKeyMap {
	return outputKeyMap{
		up: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "[Logs] up"),
		),
		down: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "[Logs] down"),
		),
	}
}

type OutputModel struct {
	keyMap   outputKeyMap
	viewport viewport.Model
	writer   *LogWriter
	content  string
}

func newViewPort(w, h int) viewport.Model {
	vp := viewport.New(w, h)
	vp.KeyMap = viewport.KeyMap{}
	return vp
}

func NewOutputModel() OutputModel {
	r := ring.New(200)
	return OutputModel{
		keyMap:   newOutputKeyMap(),
		viewport: newViewPort(0, 0),
		writer: &LogWriter{
			logs: r,
		},
	}
}

func (m *OutputModel) GetLogWriter() *LogWriter {
	return m.writer
}

func (m *OutputModel) setContent() {
	prev := m.content
	m.content = m.writer.getLog()
	if prev != m.content {
		m.viewport.SetContent(m.content)
		m.viewport.GotoBottom()
	}
}

func (m *OutputModel) SetSize(w, h int) {
	m.viewport = newViewPort(w, h)
	m.viewport.SetContent(m.content)
}

func (m OutputModel) Init() tea.Cmd {
	return nil
}

func (m OutputModel) Update(msg tea.Msg) (OutputModel, tea.Cmd) {
	var cmd tea.Cmd
	m.setContent()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.down):
			m.viewport.ScrollDown(1)
		case key.Matches(msg, m.keyMap.up):
			m.viewport.ScrollUp(1)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

func (m OutputModel) View() string {
	return m.viewport.View()
}

func (m OutputModel) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keyMap.up,
		m.keyMap.down,
	}
}

func (m OutputModel) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.up,
		m.keyMap.down,
	}}
}
