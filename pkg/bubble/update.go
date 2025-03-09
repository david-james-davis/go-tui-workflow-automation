package bubble

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case ResultMsg:
		return m.handleResult(msg)
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	}
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		if m.State == DefaultMenu {
			return m, tea.Quit
		}
	case "esc":
		return m.handleEscape()
	case "enter":
		return m.handleEnter()
	}

	// Handle list updates
	return m.handleListUpdate(msg)
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.State {
	case DefaultMenu:
		return m.handleDefaultMenuSelection()
	case GoMenu:
		return m.handleGoMenuSelection()
	case InputProjectName:
		m.State = InputDirectory
		m.projectInput.Blur()
		m.dirPathInput.Focus()
		return m, textinput.Blink
	case InputDirectory:
		return m.handleDirectoryInput()
	}
	return m, nil
}

func (m Model) handleDefaultMenuSelection() (tea.Model, tea.Cmd) {
	if i := m.list.Index(); i != -1 {
		selectedItem := m.list.Items()[i]
		if selectedItem.FilterValue() == string(Choice2) {
			m.State = GoMenu
			return m, nil
		}
		return m.transitionToDirectoryInput()
	}
	return m, nil
}

func (m Model) handleGoMenuSelection() (tea.Model, tea.Cmd) {
	if i := m.goList.Index(); i != -1 {
		return m.transitionToDirectoryInput()
	}
	return m, nil
}

func (m Model) handleDirectoryInput() (tea.Model, tea.Cmd) {
	m.State = Building

	// Check if we came from GoMenu (Choice2 path)
	if m.list.Index() != -1 && m.list.Items()[m.list.Index()].FilterValue() == string(Choice2) {
		if i := m.goList.Index(); i != -1 {
			selectedItem := m.goList.Items()[i]
			if handler, exists := m.goMenuHandlers[selectedItem.FilterValue()]; exists {
				return m, handler(m.projectInput.Value(), m.dirPathInput.Value())
			}
		}
		return m, nil
	}

	// Handle Default menu items (Choice1 and Choice3)
	if i := m.list.Index(); i != -1 {
		selectedItem := m.list.Items()[i]
		if handler, exists := m.defaultMenuHandlers[selectedItem.FilterValue()]; exists {
			return m, handler(m.projectInput.Value(), m.dirPathInput.Value())
		}
	}

	return m, nil
}

func (m Model) transitionToDirectoryInput() (tea.Model, tea.Cmd) {
	m.State = InputProjectName
	m.projectInput.Focus()
	return m, textinput.Blink
}

func (m Model) handleEscape() (tea.Model, tea.Cmd) {
	m.State = DefaultMenu
	return m, nil
}

func (m Model) handleListUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.State {
	case DefaultMenu:
		m.list, cmd = m.list.Update(msg)
	case GoMenu:
		m.goList, cmd = m.goList.Update(msg)
	case InputProjectName:
		m.projectInput, cmd = m.projectInput.Update(msg)
	case InputDirectory:
		m.dirPathInput, cmd = m.dirPathInput.Update(msg)
	}
	return m, cmd
}

func (m Model) handleResult(msg ResultMsg) (tea.Model, tea.Cmd) {
	m.State = ShowResult
	m.resultMessage = msg.Message
	m.isSuccess = msg.Success
	return m, nil
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	if m.State == DefaultMenu {
		m.list.SetSize(msg.Width, msg.Height)
	} else {
		m.goList.SetSize(msg.Width, msg.Height)
	}
	return m, nil
}
