package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestConfirmOutput(t *testing.T) {
	t.Run("show up with short message and press enter", func(t *testing.T) {
		wm, msgChan := newWrappedModel(NewConfirmModel())
		defer close(msgChan)
		tm := teatest.NewTestModel(
			t, wm,
			teatest.WithInitialTermSize(300, 100),
		)
		msg := "test msg"
		tm.Send(showDialogMsg{
			msg: msg,
			callback: tea.Sequence(
				func() tea.Msg {
					return callbackMsg{}
				},
			),
		})

		out := waitForString(t, tm, "Cancel")
		teatest.RequireEqualOutput(t, out)

		tm.Send(tea.KeyMsg{
			Type: tea.KeyEnter,
		})

		wm.waitForMsgs(t, []any{
			BlurConfirmDialogMsg{},
			callbackMsg{},
			UpdateLayoutMsg{},
		})
		waitForEmpty(t, tm)

		tm.Send(quitMsg{})

		tm.WaitFinished(t)
	})

	t.Run("show up with long message and press esc", func(t *testing.T) {
		wm, msgChan := newWrappedModel(NewConfirmModel())
		defer close(msgChan)
		tm := teatest.NewTestModel(
			t, wm,
			teatest.WithInitialTermSize(300, 100),
		)

		msg := "Update the super-super-long-long-package-name package?"
		tm.Send(showDialogMsg{
			msg:      msg,
			callback: nil,
		})

		out := waitForString(t, tm, "Cancel")
		teatest.RequireEqualOutput(t, out)

		tm.Send(tea.KeyMsg{
			Type: tea.KeyEsc,
		})

		wm.waitForMsgs(t, []any{
			BlurConfirmDialogMsg{},
			UpdateLayoutMsg{},
		})
		waitForEmpty(t, tm)

		tm.Send(quitMsg{})

		tm.WaitFinished(t)
	})
}
