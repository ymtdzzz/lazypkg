package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type managersKeyMap struct {
	Toggle   key.Binding
	Select   key.Binding
	Check    key.Binding
	CheckAll key.Binding
	Update   key.Binding
}

func newManagersKeyMap() managersKeyMap {
	return managersKeyMap{
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle check"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "right", "l"),
			key.WithHelp("enter | l | →", "select"),
		),
		Check: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "check update"),
		),
		CheckAll: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "check update (all)"),
		),
		Update: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "update"),
		),
	}
}

type ManagersModel struct {
	keyMap     managersKeyMap
	spinner    spinner.Model
	spinnerStr *string
	mgrToIdx   map[string]int
	idxToMgr   map[int]string
	list       list.Model
	pkglists   map[string]*PackagesModel
	focus      *bool
	selection  map[int]bool
	loading    map[int]bool
}

func NewManagersModel(mgrs []string, pkglists map[string]*PackagesModel) ManagersModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	ss := s.View()

	selection := map[int]bool{}
	loading := map[int]bool{}
	focus := false

	var items []list.Item
	mgrToIdx := map[string]int{}
	idxToMgr := map[int]string{}
	for i, mgr := range mgrs {
		items = append(items, item{
			icon:  pkglists[mgr].Icon(),
			title: mgr,
		})
		mgrToIdx[mgr] = i
		idxToMgr[i] = mgr
	}
	l := list.New(
		items,
		newItemDelegate(&ss, selection, loading, &focus),
		0,
		0,
	)
	l.Title = fmt.Sprintf("Package Managers")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.Styles.Title = blurTitleStyle
	l.Styles.HelpStyle = helpStyle
	km := newManagersKeyMap()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{km.Toggle, km.Select}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{km.Toggle, km.Select}
	}

	return ManagersModel{
		keyMap:     km,
		spinner:    s,
		spinnerStr: &ss,
		mgrToIdx:   mgrToIdx,
		idxToMgr:   idxToMgr,
		list:       l,
		pkglists:   pkglists,
		selection:  selection,
		loading:    loading,
		focus:      &focus,
	}
}

func (m ManagersModel) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick}

	for _, pkg := range m.pkglists {
		cmds = append(cmds, pkg.getPackagesCmd())
	}

	return tea.Sequence(cmds...)
}

func (m ManagersModel) Update(msg tea.Msg) (ManagersModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{m.spinner.Tick}

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		*m.spinnerStr = m.spinner.View()
		return m, cmd
	case getPackageStartMsg:
		if idx, ok := m.mgrToIdx[msg.name]; ok {
			m.loading[idx] = true
		}
	case getPackageFinishMsg:
		if idx, ok := m.mgrToIdx[msg.name]; ok {
			m.loading[idx] = false
		}
	}

	if *m.focus {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keyMap.Select):
				if item := m.list.SelectedItem(); item != nil {
					cmds = append(cmds, func() tea.Msg {
						return FocusPackagesMsg{
							Name: item.FilterValue(),
						}
					})
				}
			case key.Matches(msg, m.keyMap.Toggle):
				idx := m.list.Index()
				if v, ok := m.selection[idx]; v && ok {
					m.selection[idx] = false
				} else {
					m.selection[idx] = true
				}
			case key.Matches(msg, m.keyMap.Check):
				if item := m.list.SelectedItem(); item != nil {
					if pkg, ok := m.pkglists[item.FilterValue()]; ok {
						cmds = append(cmds, pkg.getPackagesCmd())
					}
				}
			case key.Matches(msg, m.keyMap.CheckAll):
				for _, pkg := range m.pkglists {
					cmds = append(cmds, pkg.getPackagesCmd())
				}
			case key.Matches(msg, m.keyMap.Update):
				// Bulk update
				var mgrs []string
				for i, v := range m.selection {
					if v {
						mgrs = append(mgrs, m.idxToMgr[i])
					}
				}
				if len(mgrs) > 0 {
					for i := range m.selection {
						m.selection[i] = false
					}

					subcmds := []tea.Cmd{}
					for _, mgr := range mgrs {
						subcmds = append(subcmds, func() tea.Msg {
							return updateAllPackagesMsg{
								name:      mgr,
								confirmed: true,
							}
						})
					}
					cmds = append(cmds, showDialogCmd(
						fmt.Sprintf("All packages of selected %d managers will be updated", len(mgrs)),
						tea.Sequence(subcmds...),
					))
				} else {
					// Single update
					if item := m.list.SelectedItem(); item != nil {
						mgr := item.FilterValue()
						cmds = append(cmds, showDialogCmd(
							fmt.Sprintf("All %s package will be updated", mgr),
							func() tea.Msg {
								return updateAllPackagesMsg{
									name:      mgr,
									confirmed: true,
								}
							},
						))
					}
				}
			}
		}

		previous := m.list.SelectedItem()

		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)

		now := m.list.SelectedItem()

		if previous != now {
			cmds = append(cmds, func() tea.Msg {
				return ChangeManagerSelectionMsg{
					Name: now.FilterValue(),
				}
			})
		}
	}

	for mgr, l := range m.pkglists {
		if i, ok := m.mgrToIdx[mgr]; ok {
			desc := "✓"
			if l.Count() > 0 {
				desc = fmt.Sprintf("[%d]", l.Count())
			}
			m.list.SetItem(i, item{
				icon:  m.pkglists[mgr].icon,
				title: mgr,
				desc:  desc,
			})
		}
	}

	return m, tea.Batch(cmds...)
}

func (m ManagersModel) View() string {
	return m.list.View()
}

func (m ManagersModel) IsFocus() bool {
	return *m.focus
}

func (m ManagersModel) ShortHelp() []key.Binding {
	return m.list.ShortHelp()
}

func (m ManagersModel) FullHelp() [][]key.Binding {
	return m.list.FullHelp()
}

func (m *ManagersModel) SetSize(w, h int) {
	m.list.SetSize(w, h)
}

func (m *ManagersModel) Focus(focus bool) {
	*m.focus = focus
	if focus {
		m.list.Styles.Title = titleStyle
	} else {
		m.list.Styles.Title = blurTitleStyle
	}
}
