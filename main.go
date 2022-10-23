package main

import (
	"astro/habit"
	"astro/logger"
	"astro/models/list"
	"astro/state"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger.Init()
	habit.Init()
	state.GetAll()
	p := tea.NewProgram(list.NewList(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
