package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestPasswordOutput(t *testing.T) {
	t.Run("show up and submit password", func(t *testing.T) {
		wm, msgChan := newWrappedModel(NewPasswordModel())
		defer close(msgChan)
		tm := teatest.NewTestModel(
			t, wm,
			teatest.WithInitialTermSize(300, 100),
		)
		tm.Send(passwordInputStartMsg{
			callback: func(password string) tea.Cmd {
				return func() tea.Msg {
					return callbackMsg{
						value: password,
					}
				}
			},
		})
		wm.waitForMsgs(t, []any{
			FocusPasswordDialogMsg{},
			UpdateLayoutMsg{},
		})

		out := waitForString(t, tm, "Cancel")
		teatest.RequireEqualOutput(t, out)

		// Simulate being called from another process
		tm.Send(passwordInputStartMsg{
			callback: func(password string) tea.Cmd {
				return func() tea.Msg {
					return callbackMsg{
						value: password,
					}
				}
			},
		})

		t.Run("input password", func(t *testing.T) {
			tm.Type("password")
			out := waitForString(t, tm, "********")
			teatest.RequireEqualOutput(t, out)
		})

		t.Run("submit password", func(t *testing.T) {
			tm.Send(tea.KeyMsg{
				Type: tea.KeyEnter,
			})
			wm.waitForMsgs(t, []any{
				BlurPasswordDialogMsg{},
				callbackMsg{
					value: "password",
				},
				callbackMsg{
					value: "password",
				},
				UpdateLayoutMsg{},
			})
		})

		tm.Send(quitMsg{})

		tm.WaitFinished(t)
	})

	t.Run("show up and cancel", func(t *testing.T) {
		wm, msgChan := newWrappedModel(NewPasswordModel())
		defer close(msgChan)
		tm := teatest.NewTestModel(
			t, wm,
			teatest.WithInitialTermSize(300, 100),
		)
		tm.Send(passwordInputStartMsg{
			callback: func(password string) tea.Cmd {
				return func() tea.Msg {
					return callbackMsg{
						value: password,
					}
				}
			},
		})
		wm.waitForMsgs(t, []any{
			FocusPasswordDialogMsg{},
			UpdateLayoutMsg{},
		})

		out := waitForString(t, tm, "Cancel")
		teatest.RequireEqualOutput(t, out)

		t.Run("input password", func(t *testing.T) {
			tm.Type("password")
			out := waitForString(t, tm, "********")
			teatest.RequireEqualOutput(t, out)
		})

		t.Run("cancel password input", func(t *testing.T) {
			tm.Send(tea.KeyMsg{
				Type: tea.KeyEsc,
			})
			wm.waitForMsgs(t, []any{
				BlurPasswordDialogMsg{},
				UpdateLayoutMsg{},
			})
		})

		t.Run("show up again and input form is reset", func(t *testing.T) {
			tm.Send(passwordInputStartMsg{
				callback: func(password string) tea.Cmd {
					return func() tea.Msg {
						return callbackMsg{
							value: password,
						}
					}
				},
			})
			wm.waitForMsgs(t, []any{
				FocusPasswordDialogMsg{},
				UpdateLayoutMsg{},
			})

			out := waitForString(t, tm, ">")
			teatest.RequireEqualOutput(t, out)
		})

		tm.Send(quitMsg{})

		tm.WaitFinished(t)
	})
}
