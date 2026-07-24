package main

import (
	"fmt"
	"os"
	"strings"
	"golang.design/x/clipboard"

	"github.com/Danyiyk/FunCodex/internal/decoder"
	"github.com/Danyiyk/FunCodex/internal/encoder"
	"github.com/Danyiyk/FunCodex/internal/utils"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func PrintHelp() {
	fmt.Println("Registered charmaps", len(utils.RegisteredCharmaps))
	fmt.Println("Available commands:\n -l(ength) [text]\n -c(rypt) [text] [hidden_text]\n -d(ecrypt) [crypted_text]")
}

func copyToClipboard(text string) {
	data := []byte(text)
	
	clipboard.Write(clipboard.FmtText, data)
	
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	mode        string // active mode: "default", "encrypt", "decrypt", "length"
	err         error
}

// initialize model
func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "type something..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	// ascii art banner

	vp := viewport.New(30, 5)
vp.SetContent("Welcome to FunCodex!\nType a message and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		mode:        "default",
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// update state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - 3

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

		// switch mode
		case "ctrl+d":
			m.mode = "decrypt"
			m.messages = append(m.messages, m.senderStyle.Render("funcodex: ")+"mode set to decrypt")
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
			return m, nil

		case "ctrl+e":
			m.mode = "encrypt"
			m.messages = append(m.messages, m.senderStyle.Render("funcodex: ")+"mode set to encrypt")
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
			return m, nil

		case "ctrl+l":
			m.mode = "length"
			m.messages = append(m.messages, m.senderStyle.Render("funcodex: ")+"mode set to length")
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
			return m, nil

		// process input based on mode
		case "enter":
			val := strings.TrimSpace(m.textarea.Value())
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			if val == "" {
				return m, nil
			}

			var response string

			switch m.mode {
			case "encrypt":
				parts := strings.SplitN(val, "|", 2)
				if len(parts) < 2 {
					response = m.senderStyle.Render("funcodex: ") + "error: use format 'text | hidden_text'"
				} else {
					text := strings.TrimSpace(parts[0])
					hiddenText := strings.TrimSpace(parts[1])
					encodedResult := encoder.Encode(text, hiddenText)

					// copy result to clipboard (primary and standard)
					copyToClipboard(encodedResult)

					response = fmt.Sprintf("encrypted: %s", encodedResult)
				}
			case "decrypt":
				decodedResult := decoder.Decode(val)
				response = fmt.Sprintf("decrypted: %s", decodedResult)

			case "length":
				encodableChars, maxLen := utils.GetAvailableSpace(val)
				response = fmt.Sprintf("encodable chars: %d | max hidden text: %d chars", encodableChars, maxLen)

			default:
				response = fmt.Sprintf("Select a mode", val)
			}

			m.messages = append(m.messages, response)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
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

// render view
func (m model) View() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	footer := helpStyle.Render("ctrl+e: encrypt • ctrl+d: decrypt • ctrl+l: length • esc: quit")

	return fmt.Sprintf("%s\n\n%s\n%s", m.viewport.View(), m.textarea.View(), footer)
}

func main() {
	if err := clipboard.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to init clipboard: %v\n", err)
	}
	// run tui if no flags passed
	if len(os.Args) == 1 {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error starting tui: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// cli execution
	if len(os.Args) < 3 {
		fmt.Println("Mismatching argument count!")
		PrintHelp()
		return
	}

	mode := os.Args[1]
	text := os.Args[2]

	switch mode {
	case "-c":
		if len(os.Args) < 4 {
			fmt.Println("-c [text] [text to encrypt]")
			return
		}
		hiddenText := os.Args[3]
		fmt.Println(encoder.Encode(text, hiddenText))

	case "-d":
		fmt.Println(decoder.Decode(text))

	case "-l":
		encodableCharacters, hiddenTextMaxLength := utils.GetAvailableSpace(text)
		fmt.Printf("Analysis:\n Text: \"%s\"\n Encodable characters: %d\n Hidden text max length: %d characters\n", text, encodableCharacters, hiddenTextMaxLength)

	default:
		PrintHelp()
	}
}
