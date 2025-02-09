package components

import "github.com/charmbracelet/bubbles/list"

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

type updatePackageStartMsg struct {
	name string
	pkg  string
}

type updatePackageFinishMsg struct {
	name string
	pkg  string
}

type focusManagersMsg struct{}

type focusPackagesMsg struct {
	name string
}
