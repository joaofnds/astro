package main

import (
	"astroapp/habit"
	"astroapp/models"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client := habit.NewClient()
	habits, err := client.List()
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(models.NewList(habits), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
