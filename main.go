package main

import (
	"astro/habit"
	"astro/logger"
	"astro/models/list"
	"astro/state"
	"astro/token"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := logger.Init(); err != nil {
		log.Fatal(err)
	}
	tok, err := token.Init()
	if err != nil {
		log.Fatal(err)
	}
	habit.InitClient(tok)
	if err := state.GetAll(); err != nil {
		log.Fatal(err)
	}
	p := tea.NewProgram(list.NewList(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
