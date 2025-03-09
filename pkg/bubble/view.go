package bubble

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	switch m.State {
	case DefaultMenu:
		return m.list.View()
	case GoMenu:
		return m.goList.View()
	case InputProjectName:
		return m.projectInputView()
	case InputDirectory:
		return m.dirInputView()
	case Building:
		return m.buildingView()
	case ShowResult:
		return m.resultView()
	default:
		return "Unknown state"
	}
}

func (m Model) buildingView() string {
	return fmt.Sprintf("\n\n   %s\n\n", "building...")
}

func (m Model) projectInputView() string {
	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s",
		"Enter Project Name:",
		m.projectInput.View(),
		"(Press Enter to confirm, Esc to cancel)",
	)
}

func (m Model) dirInputView() string {
	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s",
		"Enter Directory Path:",
		m.dirPathInput.View(),
		"(Press Enter to confirm, Esc to cancel)",
	)
}

func (m Model) resultView() string {
	style := lipgloss.NewStyle().MarginLeft(2).MarginTop(1)
	if m.isSuccess {
		style = style.Foreground(lipgloss.Color("green"))
	} else {
		style = style.Foreground(lipgloss.Color("red"))
	}

	return fmt.Sprintf("\n%s\n\nPress ESC to return to the main menu\n",
		style.Render(m.resultMessage))
}
