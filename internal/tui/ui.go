package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	mode        string // current mode: "default", "encrypt", "decrypt", "length"
	err         error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Type something..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent("Welcome to FunCodex!\nType a message and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		mode:        "default", // set initial mode to default
		err:         nil,
	}
}


func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - 2

		if len(m.messages) > 0 {
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")),
			)
		}
		m.viewport.GotoBottom()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+d":
			m.messages = append(
				m.messages,
				m.senderStyle.Render("FunCodex: ")+"decrypt",
			)
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")),
			)
			m.viewport.GotoBottom()
			return m, nil

		case "ctrl+e":
			m.messages = append(
				m.messages,
				m.senderStyle.Render("FunCodex: ")+"encrypt",
			)
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")),
			)
			m.viewport.GotoBottom()
			return m, nil

		case "ctrl+l":
			m.messages = append(
				m.messages,
				m.senderStyle.Render("FunCodex: ")+"length",
			)
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")),
			)
			m.viewport.GotoBottom()
			return m, nil

		case "enter":
			if strings.TrimSpace(m.textarea.Value()) == "" {
				return m, nil
			}
			m.messages = append(
				m.messages,
				m.senderStyle.Render("Tu: ")+m.textarea.Value(),
			)
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")),
			)
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, nil
		}

	case cursor.BlinkMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	var taCmd, vpCmd tea.Cmd
	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, taCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s", m.viewport.View(), m.textarea.View())
}
