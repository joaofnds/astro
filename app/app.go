package app

import (
	"astro/api"
	"astro/domain"

	tea "charm.land/bubbletea/v2"
)

// App is the root model. It owns all application state, the API client,
// the screen stack, and terminal dimensions.
type App struct {
	state  AppState
	client *api.Client
	stack  []tea.Model
	width  int
	height int
	ready  bool
}

// New creates a root model that starts in the loading state.
func New(client *api.Client) App {
	return App{
		state:  NewAppState(),
		client: client,
	}
}

// NewForTest creates an App with a pre-pushed screen, bypassing the normal
// loading flow. Test-only; not part of the public API contract.
func NewForTest(screen tea.Model) App {
	return App{
		state: NewAppState(),
		stack: []tea.Model{screen},
		ready: true,
	}
}

// SetStateForTest populates the app state. Test-only.
func (a *App) SetStateForTest(habits []*domain.Habit, groups []*domain.Group) {
	a.state.SetAll(habits, groups)
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return a, nil
}

func (a App) View() tea.View {
	return tea.NewView("")
}

func (a App) activeScreen() tea.Model {
	if len(a.stack) == 0 {
		return nil
	}
	return a.stack[len(a.stack)-1]
}
