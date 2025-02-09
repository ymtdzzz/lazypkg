package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ymtdzzz/lazypkg/components"
	"github.com/ymtdzzz/lazypkg/executors"
)

const (
	PACKAGE_MANAGER_APT      = "apt"
	PACKAGE_MANAGER_HOMEBREW = "homebrew"
)

var docStyle = lipgloss.NewStyle().
	Margin(1, 2)

var docStyleRightBorder = docStyle.
	Border(lipgloss.NormalBorder(), false, true, false, false)

var borderStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true, false, false, false).
	Padding(1)

type model struct {
	focusRight  bool
	selectedPkg string
	mgrlist     components.ManagersModel
	pkglists    map[string]*components.PackagesModel
	out         components.OutputModel
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.mgrlist.Init())
	for _, pkg := range m.pkglists {
		cmds = append(cmds, pkg.Init())
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.updateLayout(msg.Width, msg.Height)
	case components.ChangeManagerSelectionMsg:
		m.selectedPkg = msg.Name
	}

	m.mgrlist, cmd = m.mgrlist.Update(msg)
	cmds = append(cmds, cmd)
	for k, lptr := range m.pkglists {
		mm, cc := (*lptr).Update(msg)
		m.pkglists[k] = &mm
		cmd = cc
		cmds = append(cmds, cmd)
	}

	m.out, cmd = m.out.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	left := docStyleRightBorder.Render(m.mgrlist.View())

	rightTop := docStyle.Render("")
	if pkg, ok := m.pkglists[m.selectedPkg]; ok {
		rightTop = docStyle.Render(pkg.View())
	}
	rightBottom := borderStyle.Render(m.out.View())

	right := lipgloss.JoinVertical(
		lipgloss.Left,
		rightTop,
		rightBottom,
	)

	layout := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)

	return layout
}

func (m *model) updateLayout(w, h int) {
	leftWidth := int(float64(w) * 0.4)
	rightWidth := w - leftWidth
	rightHeight := h / 2

	dfw, dfh := docStyle.GetFrameSize()
	bfw, bfh := borderStyle.GetFrameSize()

	m.mgrlist.SetSize(leftWidth-dfw, h-dfh)

	for _, l := range m.pkglists {
		l.SetSize(rightWidth-dfw, rightHeight-dfh)
	}

	m.out.SetSize(rightWidth-bfw, rightHeight-bfh)
}

func main() {
	apt := components.NewPackageModel(PACKAGE_MANAGER_APT, &executors.AptExecutor{})
	homebrew := components.NewPackageModel(PACKAGE_MANAGER_HOMEBREW, &executors.HomebrewExecutor{})
	pkglists := map[string]*components.PackagesModel{
		PACKAGE_MANAGER_APT:      &apt,
		PACKAGE_MANAGER_HOMEBREW: &homebrew,
	}

	mgrlist := components.NewManagersModel([]string{
		PACKAGE_MANAGER_APT,
		PACKAGE_MANAGER_HOMEBREW,
	}, pkglists)
	mgrlist.Focus(true)

	out := components.NewOutputModel()
	log.SetOutput(out.GetLogWriter())

	m := model{false, PACKAGE_MANAGER_APT, mgrlist, pkglists, out}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
