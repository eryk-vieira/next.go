package build_tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/eryk-vieira/next.go/internal/cli/build"
	"github.com/eryk-vieira/next.go/internal/cli/types"
)

var p *tea.Program

type model struct {
	settings *types.Settings
	spinner  spinner.Model
	done     bool
	errors   []build.Errors
}

var (
	opacityStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	doneStyle    = lipgloss.NewStyle().Margin(1, 1)
	checkMark    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	errorMark    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f57c67")).SetString("x")
	error        = lipgloss.NewStyle().Margin(1, 1).Foreground(lipgloss.Color("#f57c67"))
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
	return tea.Batch(m.buildApplication, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errorType:
		m.errors = msg
		return m, tea.Quit
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
	if len(m.errors) > 0 {
		strBuilder := strings.Builder{}

		for _, e := range m.errors {
			strBuilder.WriteString(e.FilePath)
			strBuilder.WriteString("\n")
			strBuilder.WriteString("\n")
			strBuilder.WriteString(opacityStyle.Render(e.Error.Error()))
			strBuilder.WriteString("\n")
			strBuilder.WriteString("\n")
		}

		strBuilder.WriteString("\n")
		strBuilder.WriteString(opacityStyle.Render(fmt.Sprintf("%s Build Failed!.", errorMark)))
		strBuilder.WriteString("\n")

		return error.Render(strBuilder.String())
	}

	if m.done && len(m.errors) <= 0 {
		return doneStyle.Render(fmt.Sprintf("%s Build complete!.", checkMark))
	}

	spin := m.spinner.View() + " "

	return spin + "Next.go is building your application..."
}

type buildDone bool
type errorType []build.Errors

func (m model) buildApplication() tea.Msg {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)

			p.Quit()
		}
	}()

	builder := build.NewBuilder(m.settings)
	_, errorList := builder.Build()

	time.Sleep(2 * time.Second)

	if len(errorList) > 0 {
		m.errors = errorList
		return errorType(m.errors)
	}

	time.Sleep(2 * time.Second)

	return buildDone(true)
}

func Run(s *types.Settings) {
	p = tea.NewProgram(newModel(s))

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
