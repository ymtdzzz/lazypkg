package components

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ymtdzzz/lazypkg/executors"
)

type packagesKeyMap struct {
	Toggle    key.Binding
	Back      key.Binding
	Update    key.Binding
	UpdateAll key.Binding
}

func newPackagesKeyMap() packagesKeyMap {
	return packagesKeyMap{
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle check"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace", "left", "h"),
			key.WithHelp("backspace | h | ←", "back"),
		),
		Update: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "update"),
		),
		UpdateAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "update all"),
		),
	}
}

type PackagesModel struct {
	keyMap     packagesKeyMap
	spinner    spinner.Model
	spinnerStr *string
	name       string
	executor   executors.Executor
	pkgToIdx   map[string]int
	idxToPkg   map[int]string
	list       list.Model
	focus      *bool
	selection  map[int]bool
	loading    map[int]bool
}

func NewPackageModel(name string, executor executors.Executor) PackagesModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	ss := s.View()

	selection := map[int]bool{}
	loading := map[int]bool{}
	focus := false
	l := list.New(
		[]list.Item{},
		newItemDelegate(&ss, selection, loading, &focus),
		0,
		0,
	)
	l.Title = fmt.Sprintf("Packages [%s]", name)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.Styles.Title = blurTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	km := newPackagesKeyMap()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{km.Toggle, km.Back, km.Update, km.UpdateAll}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{km.Toggle, km.Back, km.Update, km.UpdateAll}
	}

	return PackagesModel{
		keyMap:     km,
		name:       name,
		spinner:    s,
		spinnerStr: &ss,
		executor:   executor,
		pkgToIdx:   map[string]int{},
		list:       l,
		selection:  selection,
		loading:    loading,
		focus:      &focus,
	}
}

func (m PackagesModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m PackagesModel) Update(msg tea.Msg) (PackagesModel, tea.Cmd) {
	var cmd tea.Cmd

	cmds := []tea.Cmd{m.spinner.Tick}

	switch msg := msg.(type) {
	case packageUpdateMsg:
		if msg.name == m.name {
			pkgToIdx := map[string]int{}
			idxToPkg := map[int]string{}
			for k := range m.selection {
				delete(m.selection, k)
			}
			for i, item := range msg.items {
				pkgToIdx[item.FilterValue()] = i
				idxToPkg[i] = item.FilterValue()
			}
			m.pkgToIdx = pkgToIdx
			m.idxToPkg = idxToPkg
			return m, tea.Sequence(
				m.list.SetItems(msg.items),
				func() tea.Msg {
					return getPackageFinishMsg{name: m.name}
				},
			)
		}
	case updatePackagesStartMsg:
		if msg.name == m.name {
			for _, pkg := range msg.pkgs {
				if i, ok := m.pkgToIdx[pkg]; ok {
					m.loading[i] = true
				}
			}
		}
	case updatePackagesFinishMsg:
		if msg.name == m.name {
			for _, pkg := range msg.pkgs {
				if i, ok := m.pkgToIdx[pkg]; ok {
					m.loading[i] = false
				}
			}
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		*m.spinnerStr = m.spinner.View()
		return m, cmd
	}

	if *m.focus {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keyMap.Toggle):
				idx := m.list.Index()
				if v, ok := m.selection[idx]; v && ok {
					m.selection[idx] = false
				} else {
					m.selection[idx] = true
				}
			case key.Matches(msg, m.keyMap.Back):
				cmds = append(cmds, func() tea.Msg {
					return FocusManagersMsg{}
				})
			case key.Matches(msg, m.keyMap.Update):
				// Bulk update
				var pkgs []string
				for i, v := range m.selection {
					if v {
						pkgs = append(pkgs, m.idxToPkg[i])
					}
				}
				if len(pkgs) > 0 {
					for i := range m.selection {
						m.selection[i] = false
					}
					cmds = append(cmds, tea.Sequence(
						func() tea.Msg {
							return updatePackagesStartMsg{name: m.name, pkgs: pkgs}
						},
						m.BulkUpdatePackageCmd(pkgs),
					))
				} else {
					// Single update
					if item := m.list.SelectedItem(); item != nil {
						pkg := item.FilterValue()
						cmds = append(cmds, tea.Sequence(
							func() tea.Msg {
								return updatePackagesStartMsg{name: m.name, pkgs: []string{pkg}}
							},
							m.UpdatePackageCmd(pkg),
						))
					}
				}
			case key.Matches(msg, m.keyMap.UpdateAll):
				var pkgs []string
				for k := range m.pkgToIdx {
					pkgs = append(pkgs, k)
				}
				if len(pkgs) > 0 {
					cmds = append(cmds, tea.Sequence(
						func() tea.Msg {
							return updatePackagesStartMsg{name: m.name, pkgs: pkgs}
						},
						m.BulkUpdatePackageCmd(pkgs),
					))
				}
			}
		}

		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m PackagesModel) View() string {
	return m.list.View()
}

func (m PackagesModel) IsFocus() bool {
	return *m.focus
}

func (m PackagesModel) ShortHelp() []key.Binding {
	return m.list.ShortHelp()
}

func (m PackagesModel) FullHelp() [][]key.Binding {
	return m.list.FullHelp()
}

func (m *PackagesModel) SetSize(w, h int) {
	m.list.SetSize(w, h)
}

func (m *PackagesModel) Focus(focus bool) {
	*m.focus = focus
	if focus {
		m.list.Styles.Title = titleStyle
	} else {
		m.list.Styles.Title = blurTitleStyle
	}
}

func (m *PackagesModel) log(text string) {
	log.Printf("[%s] %s", m.name, text)
}

func (m *PackagesModel) GetPackagesCmd() tea.Cmd {
	return func() tea.Msg {
		pkgs, err := m.executor.GetPackages()
		if err != nil {
			m.log(fmt.Sprintf("Error fetching packages: %v", err))
			pkgs = []*executors.PackageInfo{}
		}

		return packageUpdateMsg{m.name, getPackageItems(pkgs)}
	}
}

func (m *PackagesModel) UpdatePackageCmd(pkg string) tea.Cmd {
	return func() tea.Msg {
		err := m.executor.Update(pkg, "")
		if err == executors.PasswordErr {
			return passwordInputStartMsg{
				Callback: func(password string) tea.Cmd {
					return func() tea.Msg {
						err := m.executor.Update(pkg, password)
						if err != nil {
							m.log(fmt.Sprintf("Error update pacakge (after password input): %v", err))
						}
						return updatePackagesFinishMsg{
							name: m.name,
							pkgs: []string{pkg},
						}
					}
				},
			}
		} else if err != nil {
			m.log(fmt.Sprintf("Error update pacakge: %v", err))
		}

		return updatePackagesFinishMsg{
			name: m.name,
			pkgs: []string{pkg},
		}
	}
}

func (m *PackagesModel) BulkUpdatePackageCmd(pkgs []string) tea.Cmd {
	return func() tea.Msg {
		err := m.executor.BulkUpdate(pkgs, "")
		if err == executors.PasswordErr {
			return passwordInputStartMsg{
				Callback: func(password string) tea.Cmd {
					return func() tea.Msg {
						err := m.executor.BulkUpdate(pkgs, password)
						if err != nil {
							m.log(fmt.Sprintf("Error update pacakge (after password input): %v", err))
						}
						return updatePackagesFinishMsg{
							name: m.name,
							pkgs: pkgs,
						}
					}
				},
			}
		} else if err != nil {
			m.log(fmt.Sprintf("Error update pacakge: %v", err))
		}

		return updatePackagesFinishMsg{
			name: m.name,
			pkgs: pkgs,
		}
	}
}

func getPackageItems(pkgs []*executors.PackageInfo) []list.Item {
	rows := []list.Item{}
	for _, pkg := range pkgs {
		desc := fmt.Sprintf("%s %s", pkg.Version, pkg.Arch)
		rows = append(rows, item{
			title: pkg.Name,
			desc:  desc,
		})
	}

	return rows
}
