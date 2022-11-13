package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/dmorn/scanner-trampoline/config"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	ShowConfig key.Binding
	Clear      key.Binding
	Help       key.Binding
	Quit       key.Binding

	// Not to be listed in help
	Enter key.Binding
}

func (k keyMap) help() [][]key.Binding {
	return [][]key.Binding{
		{k.ShowConfig, k.Clear}, // first column
		{k.Help, k.Quit},        // second column
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return k.help()[1]
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return k.help()
}

type model struct {
	keys   keyMap
	help   help.Model
	input  textinput.Model
	area   textarea.Model
	config config.Config

	lastScan   string
	err        error
	showConfig bool
}

func newModel(cfg config.Config) *model {
	input := textinput.New()
	input.Placeholder = "scan it!"
	input.Focus()
	input.SetCursorMode(textinput.CursorBlink)

	help := help.New()
	help.ShowAll = true

	area := textarea.New()
	configString, err := config.Marshal(&cfg)
	if err != nil {
		panic(err)
	}
	area.InsertString(string(configString))

	keys := keyMap{
		ShowConfig: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "show configuration"),
		),
		Clear: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "reset view (clear error, config, last scanned text)"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "attempt opening scanned string"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit the program"),
		),
	}

	return &model{
		keys:   keys,
		help:   help,
		input:  input,
		area:   area,
		config: cfg,
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.ShowConfig):
			m.showConfig = true
		case key.Matches(msg, m.keys.Clear):
			m.reset()
			m.input.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			m.lastScan = m.input.Value()
			m.input.Reset()
			return m, func() tea.Msg {
				return open(m.lastScan)
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
			return clearMsg{}
		})
	case clearMsg:
		m.reset()

		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) reset() {
	m.err = nil
	m.lastScan = ""
	m.showConfig = false
}

func (m *model) View() string {
	var body string
	header := fmt.Sprintf("%s", m.input.View())

	body += header + "\n"
	body += m.lastScan + "\n"

	if m.err != nil {
		body += fmt.Sprintf("error! %v\n", m.err)
	} else {
		body += "\n"
	}

	if m.showConfig {
		body += fmt.Sprintf("%s\n", m.area.View())
	}

	helpView := m.help.View(m.keys)
	body += helpView

	return body
}

type (
	errMsg   error
	clearMsg struct{}
)

func open(path string) tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	cmd := exec.CommandContext(ctx, "open", path)
	if err := cmd.Run(); err != nil {
		return errMsg(err)
	}
	return nil
}

func errorf(err error) {
	fmt.Printf("%s: error: %v", os.Args[0], err)
	os.Exit(1)
}

func main() {
	cfg := config.Default()
	if len(os.Args) > 1 {
		path := os.Args[1]
		raw, err := os.ReadFile(path)
		if err != nil {
			errorf(fmt.Errorf("read configuration file: %w", err))
		}

		if err := config.Unmarshal(raw, &cfg); err != nil {
			errorf(fmt.Errorf("unmarshal configuration file: %w", err))
		}
	}

	m := newModel(cfg)
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		errorf(err)
	}
}
