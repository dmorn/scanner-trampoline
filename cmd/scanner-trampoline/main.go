package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type model struct {
	input    textinput.Model
	lastScan string
	err      error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.lastScan = m.input.Value()
			m.input.Reset()
			return m, func() tea.Msg {
				return open(m.lastScan)
			}
			return m, nil

		case tea.KeySpace:
			m.input.Reset()
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
			return clearMsg{}
		})
	case clearMsg:
		m.err = nil
		m.lastScan = ""
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var body string
	header := fmt.Sprintf("%s\n", m.input.View())

	body += header + "\n"
	body += m.lastScan + "\n"

	if m.err != nil {
		body += fmt.Sprintf("error! %v\n", m.err)
	} else {
		body += "\n"
	}

	footer := "(ctrl-c to quit)"

	body += footer
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

func main() {
	input := textinput.New()
	input.Placeholder = "scan it!"
	input.Focus()
	input.SetCursorMode(textinput.CursorBlink)

	model := model{
		input: input,
	}
	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		fmt.Printf("%s: error: %v", os.Args[0], err)
		os.Exit(1)
	}
}
