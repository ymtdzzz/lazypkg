package components

import (
	"container/ring"
	"strings"
	"sync"

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

type OutputModel struct {
	viewport viewport.Model
	writer   *LogWriter
	content  string
}

func newViewPort(w, h int) viewport.Model {
	return viewport.New(w, h)
}

func NewOutputModel() OutputModel {
	r := ring.New(200)
	return OutputModel{
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
	m.viewport.SetContent(m.writer.getLog())
	m.viewport.GotoBottom()
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
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

func (m OutputModel) View() string {
	return m.viewport.View()
}
