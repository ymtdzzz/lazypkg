package components

import (
	"reflect"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

type wrappedModel struct {
	model   any
	msgChan chan tea.Msg
}

func newWrappedModel(model any) (wrappedModel, chan tea.Msg) {
	msgChan := make(chan tea.Msg, 100)
	return wrappedModel{
		model:   model,
		msgChan: msgChan,
	}, msgChan
}

func (wm wrappedModel) Init() tea.Cmd {
	var cmd tea.Cmd

	switch m := wm.model.(type) {
	case ConfirmModel:
		return m.Init()
	}

	return cmd
}
func (wm wrappedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	wm.msgChan <- msg

	switch msg.(type) {
	case quitMsg:
		return wm, tea.Quit
	}

	switch m := wm.model.(type) {
	case ConfirmModel:
		wm.model, cmd = m.Update(msg)
	}

	return wm, cmd
}

func (wm wrappedModel) View() string {
	switch m := wm.model.(type) {
	case ConfirmModel:
		return m.View()
	}

	return ""
}

func (wm wrappedModel) waitForMsgs(t *testing.T, targets []any) {
	t.Helper()

	receivedMsgs := make([]any, 0, len(targets))

	for len(targets) > 0 {
		select {
		case msg := <-wm.msgChan:
			for i, target := range targets {
				if reflect.TypeOf(target) == reflect.TypeOf(msg) {
					receivedMsgs = append(receivedMsgs, msg)
					targets = append(targets[:i], targets[i+1:]...)
					break
				}
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for message")
			return
		}
	}
}

type quitMsg struct{}
type callbackMsg struct{}

func waitForString(t *testing.T, tm *teatest.TestModel, s string) (result []byte) {
	t.Helper()
	teatest.WaitFor(
		t,
		tm.Output(),
		func(b []byte) bool {
			matched := strings.Contains(string(b), s)
			if matched {
				result = b
			}
			return matched
		},
	)
	return
}

func waitForEmpty(t *testing.T, tm *teatest.TestModel) {
	t.Helper()
	teatest.WaitFor(
		t,
		tm.Output(),
		func(b []byte) bool {
			return string(b) == ""
		},
	)
	return
}
