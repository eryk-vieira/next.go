package build_tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/eryk-vieira/next.go/cli/nextgo/build"
	"github.com/eryk-vieira/next.go/cli/nextgo/types"
)

type model struct {
	settings *types.Settings
	spinner  spinner.Model
	done     bool
}

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 1)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func newModel(settings *types.Settings) model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return model{
		spinner:  s,
		settings: settings,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.buildApplication)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case buildDone:
		m.done = true
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		return doneStyle.Render(fmt.Sprintf("%s Build complete!.", checkMark))
	}

	spin := m.spinner.View() + " "

	return spin + "Next.go is building your application..."
}

type buildDone bool

func (m model) buildApplication() tea.Msg {
	builder := build.NewBuilder(m.settings)
	builder.Build()

	return buildDone(true)
}

func Run(s *types.Settings) {
	if _, err := tea.NewProgram(newModel(s)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
