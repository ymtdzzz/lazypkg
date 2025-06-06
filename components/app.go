package components

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ymtdzzz/lazypkg/executors"
)

const (
	PACKAGE_MANAGER_APT      = "apt"
	PACKAGE_MANAGER_HOMEBREW = "homebrew"
	PACKAGE_MANAGER_DOCKER   = "docker"
	PACKAGE_MANAGER_NPM      = "npm"
	PACKAGE_MANAGER_GEM      = "gem"

	ICON_APT      = '\uebc6'
	ICON_HOMEBREW = '\uf0fc'
	ICON_DOCKER   = '\uf21f'
	ICON_NPM      = '\ued0d'
	ICON_GEM      = '\uf219'
)

var (
	docStyle = lipgloss.NewStyle().
			Margin(1, 2)
	docStyleRightBorder = docStyle.
				Border(lipgloss.NormalBorder(), false, true, false, false)
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			Padding(1)
	helpStyle = lipgloss.NewStyle().PaddingLeft(2)
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

type AppModel struct {
	config       Config
	keyMap       mainKeyMap
	w, h         int
	selectedPkg  string
	mgrlist      ManagersModel
	pkglists     map[string]*PackagesModel
	out          OutputModel
	pdialog      PasswordModel
	cdialog      ConfirmModel
	prevCmd      tea.Cmd
	globalKeyMap globalKeyMap
	help         help.Model
}

func NewAppModel(config Config) (AppModel, error) {
	apt := NewPackageModel(config, PACKAGE_MANAGER_APT, ICON_APT, &executors.AptExecutor{})
	homebrew := NewPackageModel(config, PACKAGE_MANAGER_HOMEBREW, ICON_HOMEBREW, &executors.HomebrewExecutor{})
	de, err := executors.NewDockerExecutor()
	if err != nil {
		return AppModel{}, err
	}
	docker := NewPackageModel(config, PACKAGE_MANAGER_DOCKER, ICON_DOCKER, de)
	npm := NewPackageModel(config, PACKAGE_MANAGER_NPM, ICON_NPM, &executors.NpmExecutor{})
	gem := NewPackageModel(config, PACKAGE_MANAGER_GEM, ICON_GEM, &executors.GemExecutor{})

	baseMgrs := map[string]*PackagesModel{
		PACKAGE_MANAGER_APT:      &apt,
		PACKAGE_MANAGER_HOMEBREW: &homebrew,
		PACKAGE_MANAGER_NPM:      &npm,
		PACKAGE_MANAGER_GEM:      &gem,
	}
	if config.Demo {
		baseMgrs = getDemoMgrs(config)
	}
	optionalMgrs := map[string]*PackagesModel{
		PACKAGE_MANAGER_DOCKER: &docker,
	}
	for exclude := range config.Excludes {
		if _, ok := baseMgrs[exclude]; !ok {
			validValues := make([]string, 0, len(baseMgrs))
			for value := range baseMgrs {
				validValues = append(validValues, value)
			}
			fmt.Printf("Invalid exclude option. Valid values: %s\n", strings.Join(validValues, ", "))
			os.Exit(1)
		}
	}
	for enable := range config.EnableFeatures {
		if _, ok := optionalMgrs[enable]; !ok {
			validValues := make([]string, 0, len(optionalMgrs))
			for value := range optionalMgrs {
				validValues = append(validValues, value)
			}
			fmt.Printf("Invalid enable-feature option. Valid values: %s\n", strings.Join(validValues, ", "))
			os.Exit(1)
		}
	}

	var (
		pkglists = map[string]*PackagesModel{}
		mgrs     = make([]string, 0, len(baseMgrs))
	)
	for k, m := range baseMgrs {
		if !m.Valid() {
			continue
		}
		if v, ok := config.Excludes[k]; ok && v {
			continue
		}
		pkglists[k] = m
		mgrs = append(mgrs, k)
	}
	for k, m := range optionalMgrs {
		if !m.Valid() {
			continue
		}
		if v, ok := config.EnableFeatures[k]; !ok || !v {
			continue
		}
		pkglists[k] = m
		mgrs = append(mgrs, k)
	}
	if len(mgrs) == 0 {
		fmt.Println("No package managers are available")
		os.Exit(0)
	}
	sort.Slice(mgrs, func(i, j int) bool {
		return mgrs[i] < mgrs[j]
	})
	mgrlist := NewManagersModel(mgrs, pkglists)
	mgrlist.Focus(true)

	out := NewOutputModel()
	log.SetOutput(out.GetLogWriter())

	pdialog := NewPasswordModel()
	cdialog := NewConfirmModel()

	km := newMainKeyMap()
	globalKeyMap := newGlobalKeyMap(km, mgrlist, out)
	help := help.New()

	return AppModel{
		config: config,
		keyMap: km,
		w:      0, h: 0,
		selectedPkg:  mgrs[0],
		mgrlist:      mgrlist,
		pkglists:     pkglists,
		out:          out,
		pdialog:      pdialog,
		cdialog:      cdialog,
		prevCmd:      nil,
		globalKeyMap: globalKeyMap,
		help:         help,
	}, nil
}

func (m AppModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, 100)

	cmds = append(cmds, m.mgrlist.Init())
	for _, pkg := range m.pkglists {
		cmds = append(cmds, pkg.Init())
	}
	cmds = append(cmds, m.out.Init())
	cmds = append(cmds, m.pdialog.Init())

	return tea.Batch(cmds...)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0, 100)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.updateLayout(msg.Width, msg.Height)
	case UpdateLayoutMsg:
		m.updateLayout(m.w, m.h)
	case ChangeManagerSelectionMsg:
		m.selectedPkg = msg.Name
	case FocusManagersMsg:
		m.mgrlist.Focus(true)
		for _, pkg := range m.pkglists {
			pkg.Focus(false)
		}
		m.globalKeyMap = newGlobalKeyMap(m.keyMap, m.mgrlist, m.out)
	case FocusPackagesMsg:
		m.mgrlist.Focus(false)
		for k, pkg := range m.pkglists {
			if k == msg.Name {
				pkg.Focus(true)
				m.globalKeyMap = newGlobalKeyMap(m.keyMap, pkg, m.out)
			} else {
				pkg.Focus(false)
			}
		}
	case FocusPasswordDialogMsg:
		m.storePrevCmd()
		m.mgrlist.Focus(false)
		for _, pkg := range m.pkglists {
			pkg.Focus(false)
		}
		cmds = append(cmds, m.pdialog.Focus())
		m.globalKeyMap = newGlobalKeyMap(m.keyMap, m.out)
	case BlurPasswordDialogMsg:
		m.pdialog.Blur()
		cmds = append(cmds, m.prevCmd)
		m.prevCmd = nil
	case FocusConfirmDialogMsg:
		m.storePrevCmd()
		m.mgrlist.Focus(false)
		for _, pkg := range m.pkglists {
			pkg.Focus(false)
		}
		m.globalKeyMap = newGlobalKeyMap(m.keyMap, m.out)
	case BlurConfirmDialogMsg:
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

	m.cdialog, cmd = m.cdialog.Update(msg)
	cmds = append(cmds, cmd)

	m.pdialog, cmd = m.pdialog.Update(msg)
	cmds = append(cmds, cmd)

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m AppModel) View() string {
	left := docStyleRightBorder.Render(m.mgrlist.View())

	rightTop := docStyle.Render("")
	if pkg, ok := m.pkglists[m.selectedPkg]; ok {
		rightTop = docStyle.Render(pkg.View())
	}
	rightBottom := borderStyle.Render(m.out.View())

	dialog := ""
	if x, _ := m.cdialog.GetSize(); x > 0 {
		dialog = m.cdialog.View()
	}
	if x, _ := m.pdialog.GetSize(); x > 0 {
		dialog = m.pdialog.View()
	}

	right := lipgloss.JoinVertical(
		lipgloss.Left,
		rightTop,
		dialog,
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

func (m *AppModel) Close() {
	for _, pkg := range m.pkglists {
		pkg.executor.Close()
	}
}

func (m *AppModel) updateLayout(w, h int) {
	leftWidth := int(float64(w) * 0.4)
	rightWidth := w - leftWidth
	rightHeight := h / 2
	if h%2 == 0 {
		rightHeight -= 1
	} else {
		rightHeight -= 2
	}

	dfw, dfh := docStyle.GetFrameSize()
	bfw, bfh := borderStyle.GetFrameSize()
	_, pdh := m.pdialog.GetSize()
	_, cdh := m.cdialog.GetSize()

	m.mgrlist.SetSize(leftWidth-dfw, h-dfh)

	for _, l := range m.pkglists {
		l.SetSize(rightWidth-dfw, rightHeight-dfh-pdh-cdh)
	}

	m.out.SetSize(rightWidth-bfw, rightHeight-bfh)
}

func (m *AppModel) storePrevCmd() {
	m.prevCmd = nil
	if m.mgrlist.IsFocus() {
		m.prevCmd = func() tea.Msg {
			return FocusManagersMsg{}
		}
		return
	}
	for k, pkg := range m.pkglists {
		if pkg.IsFocus() {
			m.prevCmd = func() tea.Msg {
				return FocusPackagesMsg{
					Name: k,
				}
			}
		}
	}
}

func getDemoMgrs(config Config) map[string]*PackagesModel {
	apt := NewPackageModel(config, PACKAGE_MANAGER_APT, ICON_APT, executors.NewDemoExecutor(
		"apt",
		[]*executors.PackageInfo{
			{
				Name:       "curl",
				OldVersion: "7.68.0",
				NewVersion: "7.85.0",
			},
			{
				Name:       "git",
				OldVersion: "2.25.1",
				NewVersion: "2.39.0",
			},
		},
	))
	brew := NewPackageModel(config, PACKAGE_MANAGER_HOMEBREW, ICON_HOMEBREW, executors.NewDemoExecutor(
		"brew",
		[]*executors.PackageInfo{
			{
				Name:       "node",
				OldVersion: "16.14.0",
				NewVersion: "18.15.0",
			},
			{
				Name:       "python",
				OldVersion: "3.9.7",
				NewVersion: "3.11.2",
			},
			{
				Name:       "ffmpeg",
				OldVersion: "4.4",
				NewVersion: "5.1",
			},
			{
				Name:       "terraform",
				OldVersion: "1.0.11",
				NewVersion: "1.4.0",
			},
			{
				Name:       "wget",
				OldVersion: "1.21.1",
				NewVersion: "1.21.4",
			},
		},
	))
	npm := NewPackageModel(config, PACKAGE_MANAGER_NPM, ICON_NPM, executors.NewDemoExecutor(
		"brew",
		[]*executors.PackageInfo{
			{
				Name:       "react",
				OldVersion: "17.0.2",
				NewVersion: "18.2.0",
			},
			{
				Name:       "express",
				OldVersion: "4.17.1",
				NewVersion: "4.18.2",
			},
			{
				Name:       "lodash",
				OldVersion: "4.17.21",
				NewVersion: "4.18.0",
			},
		},
	))

	return map[string]*PackagesModel{
		PACKAGE_MANAGER_APT:      &apt,
		PACKAGE_MANAGER_HOMEBREW: &brew,
		PACKAGE_MANAGER_NPM:      &npm,
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
