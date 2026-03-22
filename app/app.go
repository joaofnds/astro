package app

import (
	"astro/api"
	"astro/domain"
	"astro/models/list"
	"astro/msgs"

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

// Init kicks off the asynchronous initial data load. The app shows a
// loading view until DataLoadedMsg arrives.
func (a App) Init() tea.Cmd {
	return msgs.LoadAll(a.client)
}

func (a App) activeScreen() tea.Model {
	if len(a.stack) == 0 {
		return nil
	}
	return a.stack[len(a.stack)-1]
}

// Update is the root message router. Navigation messages are consumed here
// (never forwarded). State-mutation messages update AppState first, then
// forward to the active screen. All other messages forward directly.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// --- Key input ---
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

	// --- Navigation (consumed, never forwarded) ---
	case msgs.PushScreenMsg:
		a.stack = append(a.stack, msg.Screen)
		return a, msg.Screen.Init()

	case msgs.PopScreenMsg:
		if len(a.stack) <= 1 {
			return a, tea.Quit
		}
		a.stack = a.stack[:len(a.stack)-1]
		if msg.Cmd != nil {
			cmds = append(cmds, msg.Cmd)
		}
		return a, tea.Batch(cmds...)

	// --- Terminal dimensions ---
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		// Fall through to forward to active screen.

	// --- Initial data load ---
	case msgs.DataLoadedMsg:
		a.state.SetAll(msg.Habits, msg.Groups)
		a.ready = true
		// Push the initial list screen now that data is available.
		// NewList() uses its current (no-param) signature; Plan 03 will
		// change it to accept data params directly.
		screen := list.NewList()
		a.stack = append(a.stack, screen)
		return a, screen.Init()

	// --- Fatal error ---
	case msgs.FatalErrorMsg:
		return a, tea.Sequence(
			tea.Println("Fatal: "+msg.Err.Error()),
			tea.Quit,
		)

	// --- Async results: mutate state, then forward ---
	case msgs.CheckInResultMsg:
		a.state.MergeHabit(msg.Habit)

	case msgs.HabitCreatedMsg:
		a.state.AddHabit(msg.Habit)

	case msgs.HabitDeletedMsg:
		a.state.RemoveHabit(msg.ID)

	case msgs.HabitUpdatedMsg:
		a.state.MergeHabit(msg.Habit)

	case msgs.GroupCreatedMsg:
		a.state.AddGroup(msg.Group)

	case msgs.GroupDeletedMsg:
		a.state.RemoveGroup(msg.ID)

	// Forward-only messages (state managed by screens in Plan 03):
	case msgs.AddedToGroupMsg,
		msgs.RemovedFromGroupMsg,
		msgs.ActivityUpdatedMsg,
		msgs.ActivityDeletedMsg,
		msgs.APIErrorMsg:
		// Fall through to forward to active screen.
	}

	// Forward to active screen.
	if screen := a.activeScreen(); screen != nil {
		updated, cmd := screen.Update(msg)
		a.stack[len(a.stack)-1] = updated
		cmds = append(cmds, cmd)
	}
	return a, tea.Batch(cmds...)
}

// View returns the loading screen when data hasn't arrived yet, or
// delegates to the active screen.
func (a App) View() tea.View {
	if !a.ready {
		v := tea.NewView("Loading...")
		v.AltScreen = true
		return v
	}
	if screen := a.activeScreen(); screen != nil {
		return screen.View()
	}
	v := tea.NewView("")
	v.AltScreen = true
	return v
}
