package main

import (
	"bubble/pkg/bubble"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xanzy/go-gitlab"
)

func main() {
	git, err := gitlab.NewClient(os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	m := bubble.NewModel(git)

	// Initialize with your choices

	m.InitializeList(bubble.AllDefaultChoices())
	m.InitializeGoList(bubble.AllGoChoices())
	m.RegisterHandlers()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
