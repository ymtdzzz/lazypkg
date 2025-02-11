package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateLayoutMsg struct{}

type ChangeManagerSelectionMsg struct {
	Name string
}

type packageUpdateMsg struct {
	name  string
	items []list.Item
}

type getPackageStartMsg struct {
	name string
}

type getPackageFinishMsg struct {
	name string
}

type updatePackagesStartMsg struct {
	name string
	pkgs []string
}

type updatePackagesFinishMsg struct {
	name string
	pkgs []string
}

type passwordInputStartMsg struct {
	callback func(password string) tea.Cmd
}

type showDialogMsg struct {
	msg      string
	callback tea.Cmd
}

type FocusManagersMsg struct{}

type FocusPackagesMsg struct {
	Name string
}

type FocusPasswordDialogMsg struct{}

type BlurPasswordDialogMsg struct{}

type FocusConfirmDialogMsg struct{}

type BlurConfirmDialogMsg struct{}
