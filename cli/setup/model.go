package setup

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/dmorn/scanner-trampoline/cli"
	"github.com/dmorn/scanner-trampoline/cli/scan"
	"github.com/dmorn/scanner-trampoline/config"
)

type Input struct {
	Description string
	Default     string
	Field       textinput.Model
	OnSubmit    func(*Input, *config.Config)
}

func (i *Input) View() string {
	return fmt.Sprintf("%s\n%s\n", i.Description, i.Field.View())
}

func NewInput(description, placeholder string, onSubmit func(*Input, *config.Config)) *Input {
	t := textinput.New()
	t.CursorStyle = cli.CursorStyle
	t.SetValue(placeholder)
	return &Input{
		Description: description,
		Default:     placeholder,
		Field:       t,
		OnSubmit:    onSubmit,
	}
}

func (i *Input) Focus() tea.Cmd {
	i.Field.PromptStyle = cli.FocusedStyle
	i.Field.TextStyle = cli.FocusedStyle
	return i.Field.Focus()
}

func (i *Input) Blur() {
	// Remove focused state
	i.Field.PromptStyle = cli.NoStyle
	i.Field.TextStyle = cli.NoStyle
	i.Field.Blur()
}

func (m *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	input, cmd := m.Field.Update(msg)
	m.Field = input
	return m, cmd
}

type model struct {
	focusIndex int
	inputs     []*Input
}

func New() *model {
	c := config.Default()
	m := model{
		inputs: make([]*Input, 2),
	}

	for i := range m.inputs {
		var input *Input
		switch i {
		case 0:
			input = NewInput("Command to execute at each scan", strings.Join(c.OpenCmd, " "), func(i *Input, c *config.Config) {
				c.OpenCmd = strings.Split(i.Field.Value(), " ")
			})
			input.Focus()
		case 1:
			input = NewInput("String to remove from the leading (left) part of the scanned text", "", func(i *Input, c *config.Config) {
				c.TrimLeading = i.Field.Value()
			})
		}

		m.inputs[i] = input
	}

	return &m
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				cfg := &config.Config{}
				for _, v := range m.inputs {
					v.OnSubmit(v, cfg)
				}
				// Move to scan model
				return scan.New(*cfg), nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd

}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &cli.BlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &cli.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
