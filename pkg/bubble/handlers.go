package bubble

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xanzy/go-gitlab"
)

// RegisterHandlers sets up all menu handlers for the application
func (m *Model) RegisterHandlers() {
	// Register default menu handlers
	m.RegisterDefaultHandler(string(Choice1), m.handleChoice1)
	m.RegisterDefaultHandler(string(Choice3), m.handleChoice3)

	// Register go menu handlers
	m.RegisterGoMenuHandler(string(GoTemplate1), m.handleTemplate1)
	m.RegisterGoMenuHandler(string(GoTemplate2), m.handleTemplate2)
}

// Individual handlers
func (m *Model) handleTemplate1(projectName, dirPath string) tea.Cmd {
	return func() tea.Msg {
		// Create new project
		opt := &gitlab.CreateProjectOptions{
			Name: gitlab.Ptr(projectName),
			Path: gitlab.Ptr(projectName),
		}

		newProject, _, err := m.client.Projects.CreateProject(opt)
		if err != nil {
			return ResultMsg{
				Success: false,
				Message: fmt.Sprintf("Failed to create project: %v", err),
			}
		}

		// Get source repository tree
		sourceProjectID := 123 // Replace with your template project ID
		tree, _, err := m.client.Repositories.ListTree(sourceProjectID, &gitlab.ListTreeOptions{
			Recursive: gitlab.Ptr(true),
			Ref:       gitlab.Ptr("main"),
		})
		if err != nil {
			return ResultMsg{
				Success: false,
				Message: fmt.Sprintf("Failed to get source repository tree: %v", err),
			}
		}

		// First, create all directories
		for _, item := range tree {
			if item.Type == "tree" {
				// Create directory by creating a .gitkeep file
				commitOpts := &gitlab.CreateFileOptions{
					Branch:        gitlab.Ptr("main"),
					Content:       gitlab.Ptr(""),
					CommitMessage: gitlab.Ptr(fmt.Sprintf("Create directory %s", item.Path)),
				}

				_, _, err = m.client.RepositoryFiles.CreateFile(
					newProject.ID,
					item.Path+"/.gitkeep",
					commitOpts,
				)
				if err != nil {
					return ResultMsg{
						Success: false,
						Message: fmt.Sprintf("Failed to create directory %s: %v", item.Path, err),
					}
				}
			}
		}

		// Then copy all files
		for _, item := range tree {
			if item.Type == "blob" {
				// Get file content from source
				fileContent, _, err := m.client.RepositoryFiles.GetRawFile(
					sourceProjectID,
					item.Path,
					&gitlab.GetRawFileOptions{Ref: gitlab.Ptr("main")},
				)
				if err != nil {
					return ResultMsg{
						Success: false,
						Message: fmt.Sprintf("Failed to get file %s: %v", item.Path, err),
					}
				}

				// Create file in new project
				commitOpts := &gitlab.CreateFileOptions{
					Branch:        gitlab.Ptr("main"),
					Content:       gitlab.Ptr(string(fileContent)),
					CommitMessage: gitlab.Ptr(fmt.Sprintf("Add %s", item.Path)),
				}

				_, _, err = m.client.RepositoryFiles.CreateFile(
					newProject.ID,
					item.Path,
					commitOpts,
				)
				if err != nil {
					return ResultMsg{
						Success: false,
						Message: fmt.Sprintf("Failed to create file %s: %v", item.Path, err),
					}
				}
			}
		}

		return ResultMsg{
			Success: true,
			Message: fmt.Sprintf("Successfully created project %s and copied template files", projectName),
		}
	}
}

func (m *Model) handleTemplate2(projectName, dirPath string) tea.Cmd {
	return func() tea.Msg {
		// Simulate some work
		time.Sleep(1 * time.Second)
		// Simulate an error
		return ResultMsg{
			Success: false,
			Message: "Successfuly called handler for GoTemplate2",
		}
	}
}

func (m *Model) handleChoice1(projectName, dirPath string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		// Use m.client here
		// project, _, err := m.client.Projects.GetProject(123, nil)
		// if err != nil {
		// 	return ResultMsg{
		// 		Success: false,
		// 		Message: fmt.Sprintf("Failed to get project: %v", err),
		// 	}
		// }
		if dirPath == "" {
			return ResultMsg{
				Success: false,
				Message: "Directory path cannot be empty",
			}
		}

		return ResultMsg{
			Success: true,
			Message: fmt.Sprintf("Successfully called handler for Choice1: %s", dirPath),
		}
	}
}

func (m *Model) handleChoice3(projectName, dirPath string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return ResultMsg{
			Success: true,
			Message: fmt.Sprintf("Successfully called handler for Choice3: %s", dirPath),
		}
	}
}
