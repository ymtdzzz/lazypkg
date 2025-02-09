package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type managersKeyMap struct {
	Toggle key.Binding
}

func newManagersKeyMap() managersKeyMap {
	return managersKeyMap{
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle check"),
		),
	}
}

type ManagersModel struct {
	spinner    spinner.Model
	spinnerStr *string
	mgrToIdx   map[string]int
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
	for i, mgr := range mgrs {
		items = append(items, item{
			title: mgr,
		})
		mgrToIdx[mgr] = i
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
	l.Styles.Title = blurTitleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	km := newManagersKeyMap()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{km.Toggle}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{km.Toggle}
	}

	return ManagersModel{
		spinner:    s,
		spinnerStr: &ss,
		mgrToIdx:   mgrToIdx,
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
		cmds = append(cmds, func() tea.Msg {
			return getPackageStartMsg{name: pkg.name}
		})
		cmds = append(cmds, pkg.GetPackagesCmd())
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
	case focusManagersMsg:
		m.Focus(true)
		for _, pkg := range m.pkglists {
			pkg.Focus(false)
		}
	case focusPackagesMsg:
		m.Focus(false)
		for k, pkg := range m.pkglists {
			if k == msg.name {
				pkg.Focus(true)
			} else {
				pkg.Focus(false)
			}
		}
	}

	if *m.focus {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch key := msg.String(); key {
			case "enter", "right", "l":
				if item := m.list.SelectedItem(); item != nil {
					cmds = append(cmds, func() tea.Msg {
						return focusPackagesMsg{
							name: item.FilterValue(),
						}
					})
				}
			case " ":
				idx := m.list.Index()
				if v, ok := m.selection[idx]; v && ok {
					m.selection[idx] = false
				} else {
					m.selection[idx] = true
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

	return m, tea.Batch(cmds...)
}

func (m ManagersModel) View() string {
	return m.list.View()
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
