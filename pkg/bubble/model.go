package bubble

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	gitlab "github.com/xanzy/go-gitlab"
)

type State int

const (
	DefaultMenu State = iota
	GoMenu
	InputProjectName
	InputDirectory
	Building
	ShowResult
)

type item struct {
	title string
}

type SelectionHandler func(projectName, dirPath string) tea.Cmd

type ResultMsg struct {
	Success bool
	Message string
}

type Model struct {
	State         State
	previousState State
	list          list.Model
	goList        list.Model
	selected      map[int]struct{}
	client        *gitlab.Client
	projectInput  textinput.Model
	dirPathInput  textinput.Model
	// Add handlers for different menus
	goMenuHandlers      map[string]SelectionHandler
	resultMessage       string
	isSuccess           bool
	defaultMenuHandlers map[string]SelectionHandler
}

func NewModel(client *gitlab.Client) Model {
	pi := textinput.New()
	pi.Placeholder = "Enter project name..."
	pi.Width = 40
	// pi.Focus()

	ti := textinput.New()
	ti.Placeholder = "Enter directory path..."
	ti.Width = 40
	// ti.Focus()

	return Model{
		State:               DefaultMenu,
		previousState:       DefaultMenu,
		list:                newDefaultList("Available Templates"),
		goList:              newDefaultList("Go Templates"),
		projectInput:        pi,
		dirPathInput:        ti,
		selected:            make(map[int]struct{}),
		goMenuHandlers:      make(map[string]SelectionHandler),
		defaultMenuHandlers: make(map[string]SelectionHandler),
	}
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

func newDefaultList(title string) list.Model {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#AAFF00")).
		// BorderBackground(lipgloss.Color("#AAFF00")).
		BorderForeground(lipgloss.Color("#AAFF00"))
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(true)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true) // Hide the help menu
	l.SetHeight(10)     // Set a reasonable height
	l.SetWidth(30)      // Set a reasonable width

	// Optional: Add some styling to the title
	l.Styles.Title = lipgloss.NewStyle().
		MarginLeft(2).
		MarginBottom(1).
		Bold(true)

	return l
}

func (m *Model) InitializeList(choices []string) {
	var items []list.Item
	for _, choice := range choices {
		items = append(items, item{title: choice})
	}
	m.list.SetItems(items)
}

func (m *Model) InitializeGoList(choices []string) {
	var items []list.Item
	for _, choice := range choices {
		items = append(items, item{title: choice})
	}
	m.goList.SetItems(items)
}

func (m *Model) RegisterGoMenuHandler(choice string, handler SelectionHandler) {
	m.goMenuHandlers[choice] = handler
}

func (m *Model) RegisterDefaultHandler(choice string, handler SelectionHandler) {
	m.defaultMenuHandlers[choice] = handler
}
