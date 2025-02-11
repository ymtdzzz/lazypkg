package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ymtdzzz/lazypkg/components"
	"github.com/ymtdzzz/lazypkg/executors"
)

const (
	PACKAGE_MANAGER_APT      = "apt"
	PACKAGE_MANAGER_HOMEBREW = "homebrew"
)

type mainKeyMap struct {
	quit key.Binding
}

func newMainKeyMap() mainKeyMap {
	return mainKeyMap{
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

var docStyle = lipgloss.NewStyle().
	Margin(1, 2)

var docStyleRightBorder = docStyle.
	Border(lipgloss.NormalBorder(), false, true, false, false)

var borderStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true, false, false, false).
	Padding(1)

var helpStyle = lipgloss.NewStyle().PaddingLeft(2)

type model struct {
	keyMap       mainKeyMap
	w, h         int
	focusRight   bool
	selectedPkg  string
	mgrlist      components.ManagersModel
	pkglists     map[string]*components.PackagesModel
	out          components.OutputModel
	dialog       components.PasswordModel
	prevCmd      tea.Cmd
	globalKeyMap globalKeyMap
	help         help.Model
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.mgrlist.Init())
	for _, pkg := range m.pkglists {
		cmds = append(cmds, pkg.Init())
	}
	cmds = append(cmds, m.out.Init())
	cmds = append(cmds, m.dialog.Init())

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
		m.w, m.h = msg.Width, msg.Height
		m.updateLayout(msg.Width, msg.Height)
	case components.UpdateLayoutMsg:
		m.updateLayout(m.w, m.h)
	case components.ChangeManagerSelectionMsg:
		m.selectedPkg = msg.Name
	case components.FocusManagersMsg:
		m.mgrlist.Focus(true)
		for _, pkg := range m.pkglists {
			pkg.Focus(false)
		}
		m.globalKeyMap = newGlobalKeyMap(m.keyMap, m.mgrlist, m.out)
	case components.FocusPackagesMsg:
		m.mgrlist.Focus(false)
		for k, pkg := range m.pkglists {
			if k == msg.Name {
				pkg.Focus(true)
				m.globalKeyMap = newGlobalKeyMap(m.keyMap, pkg, m.out)
			} else {
				pkg.Focus(false)
			}
		}
	case components.FocusDialogMsg:
		m.storePrevCmd()
		m.mgrlist.Focus(false)
		for _, pkg := range m.pkglists {
			pkg.Focus(false)
		}
		cmds = append(cmds, m.dialog.Focus())
		m.globalKeyMap = newGlobalKeyMap(m.keyMap, m.out)
	case components.BlurDialogMsg:
		m.dialog.Blur()
		cmds = append(cmds, m.prevCmd)
		m.prevCmd = nil
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

	m.dialog, cmd = m.dialog.Update(msg)
	cmds = append(cmds, cmd)

	m.help, cmd = m.help.Update(msg)
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
		m.dialog.View(),
		rightBottom,
	)

	layout := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)

	layoutWithHelp := lipgloss.JoinVertical(
		lipgloss.Left,
		layout,
		helpStyle.Render(m.help.View(m.globalKeyMap)),
	)

	return layoutWithHelp
}

func (m *model) updateLayout(w, h int) {
	leftWidth := int(float64(w) * 0.4)
	rightWidth := w - leftWidth
	rightHeight := h / 2

	dfw, dfh := docStyle.GetFrameSize()
	bfw, bfh := borderStyle.GetFrameSize()
	_, dh := m.dialog.GetSize()

	m.mgrlist.SetSize(leftWidth-dfw, h-dfh)

	for _, l := range m.pkglists {
		l.SetSize(rightWidth-dfw, rightHeight-dfh-dh)
	}

	m.out.SetSize(rightWidth-bfw, rightHeight-bfh)
}

func (m *model) storePrevCmd() {
	m.prevCmd = nil
	if m.mgrlist.IsFocus() {
		m.prevCmd = func() tea.Msg {
			return components.FocusManagersMsg{}
		}
		return
	}
	for k, pkg := range m.pkglists {
		if pkg.IsFocus() {
			m.prevCmd = func() tea.Msg {
				return components.FocusPackagesMsg{
					Name: k,
				}
			}
		}
	}
}

type globalKeyMap struct {
	shortHelp []key.Binding
	fullHelp  [][]key.Binding
}

func newGlobalKeyMap(km mainKeyMap, kms ...help.KeyMap) globalKeyMap {
	short := []key.Binding{km.quit}
	full := [][]key.Binding{{km.quit}}

	for _, km := range kms {
		short = append(short, km.ShortHelp()...)
		full = append(full, km.FullHelp()...)
	}

	return globalKeyMap{short, full}
}

func (m globalKeyMap) ShortHelp() []key.Binding {
	return m.shortHelp
}

// TODO: currently not implemented properly
func (m globalKeyMap) FullHelp() [][]key.Binding {
	return m.fullHelp
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

	dialog := components.NewPasswordModel()

	km := newMainKeyMap()
	globalKeyMap := newGlobalKeyMap(km, mgrlist, out)
	help := help.New()

	m := model{km, 0, 0, false, PACKAGE_MANAGER_APT, mgrlist, pkglists, out, dialog, nil, globalKeyMap, help}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
