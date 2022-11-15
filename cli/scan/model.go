package scan

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/dmorn/scanner-trampoline/cli"
	"github.com/dmorn/scanner-trampoline/config"
)

type (
	errMsg error
)

type model struct {
	input  textinput.Model
	config config.Config

	lastScan string
	err      error
}

func New(cfg config.Config) *model {
	t := textinput.New()
	t.CursorStyle = cli.CursorStyle
	t.SetCursorMode(textinput.CursorBlink)
	t.Placeholder = "scan scan scan"
	t.PromptStyle = cli.FocusedStyle
	t.TextStyle = cli.FocusedStyle
	t.Focus()

	return &model{
		input:  t,
		config: cfg,
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.lastScan = m.input.Value()
			m.input.Reset()
			return m, m.open(m.lastScan)
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	header := "Use your scanner (or keyboard) to input a string.\n"
	header += "<enter> attemps opening the scanned text\n"
	header += "<ctrl+c> or <esc> terminate program execution\n"
	header += "\n"
	header += "*** config:\n"
	header += fmt.Sprintf("open command: %v\n", m.config.OpenCmd)
	header += fmt.Sprintf("trim leading: %v\n", m.config.TrimLeading)

	header += "\n"
	header += "***\n"
	header += fmt.Sprintf("last scan: %s\n", m.lastScan)
	header += fmt.Sprintf("last err: %v\n", m.err)

	header += "\n"
	header += "***\n"

	return header + m.input.View()
}

func (m *model) open(scan string) tea.Cmd {
	path := strings.TrimLeft(scan, m.config.TrimLeading)

	if _, err := os.Stat(path); err != nil {
		return func() tea.Msg {
			return errMsg(fmt.Errorf("stat path [%s]: %v", path, err))
		}

	}

	args := append(m.config.OpenCmd, path)
	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Start(); err != nil {
		return func() tea.Msg {
			return errMsg(fmt.Errorf("opening scanned text: %v", err))
		}
	}

	return func() tea.Msg {
		err := cmd.Wait()
		if err != nil {
			return func() tea.Msg {
				return errMsg(fmt.Errorf("wait on scan [%s]: %w", path, err))
			}
		}
		return func() tea.Msg {
			return nil
		}
	}
}
