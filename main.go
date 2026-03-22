package main

import (
	"astro/api"
	"astro/config"
	"astro/logger"
	"astro/models/list"
	"astro/state"
	"astro/token"
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.Init(cfg.LogFilePath); err != nil {
		log.Fatal(err)
	}

	const baseURL = "https://astro.joaofnds.com"

	tok, err := token.Init(cfg.TokenFilePath, baseURL)
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewClient(baseURL, tok)
	state.Init(client)
	state.GetAll()

	p := tea.NewProgram(list.NewList())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
