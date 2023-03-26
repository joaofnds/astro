package main

import (
	"astro/config"
	"astro/habit"
	"astro/logger"
	"astro/models/list"
	"astro/state"
	"astro/token"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	if err := logger.Init(); err != nil {
		log.Fatal(err)
	}

	tok, err := token.Init()
	if err != nil {
		log.Fatal(err)
	}

	habit.InitClient(tok)
	state.GetAll()
	p := tea.NewProgram(list.NewList(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
