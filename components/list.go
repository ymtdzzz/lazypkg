package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	blurTitleStyle    = titleStyle.Foreground(lipgloss.Color("#777777"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type itemDelegate struct {
	spinnerStr *string
	// TODO: set via item struct
	selection map[int]bool
	loading   map[int]bool
	focus     *bool
}

func newItemDelegate(spinnerStr *string, selection, loading map[int]bool, focus *bool) itemDelegate {
	return itemDelegate{
		spinnerStr: spinnerStr,
		selection:  selection,
		loading:    loading,
		focus:      focus,
	}
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	check := " "
	if v, ok := d.selection[index]; ok && v {
		check = "*"
	}
	str := fmt.Sprintf("%s %s %s", check, i.title, i.desc)
	if v, ok := d.loading[index]; ok && v {
		str = str + " " + *d.spinnerStr
	}

	var style lipgloss.Style
	if index == m.Index() {
		str = "> " + str
		style = selectedItemStyle
	} else {
		style = itemStyle
	}

	if !(*d.focus) {
		style = style.Foreground(lipgloss.Color("#777777"))
	}

	fmt.Fprint(w, style.Render(str))
}
