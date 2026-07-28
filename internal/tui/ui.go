package main

import (
	"fmt"
	"os"
	"strings"
	"golang.design/x/clipboard"

	"github.com/Danyiyk/FunCodex/internal/config"
	"github.com/Danyiyk/FunCodex/internal/decoder"
	"github.com/Danyiyk/FunCodex/internal/encoder"
	"github.com/Danyiyk/FunCodex/internal/utils"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const asciiBanner = `
██████ ▄▄ ▄▄ ▄▄  ▄▄ ▄█████  ▄▄▄  ▄▄▄▄  ▄▄▄▄▄ ▄▄ ▄▄ 
██▄▄   ██ ██ ███▄██ ██     ██▀██ ██▀██ ██▄▄  ▀█▄█▀ 
██     ▀███▀ ██ ▀██ ▀█████ ▀███▀ ████▀ ██▄▄▄ ██ ██ 
`

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
	cfg         *config.Config
}

// initialize model
func initialModel(cfg *config.Config) model {
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

	vp := viewport.New(30, 10)

	accent := lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.AccentColor))
	banner := accent.Render(asciiBanner)

	welcomeMsg := fmt.Sprintf(
		"%s\n\nWelcome to %s!!\nSelect a mode, Type a message and press %s to send.",
		banner,
		accent.Render("FunCodex"),
		accent.Render("[Enter]"),
	)
	vp.SetContent(welcomeMsg)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.AccentColor)),
		mode:        "default",
		err:         nil,
		cfg:         cfg,
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
		case m.cfg.Keybinds.DecryptMode:
			m.mode = "decrypt"
			m.messages = append(m.messages, m.senderStyle.Render("FunCodex: ")+"mode set to decrypt")
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
			return m, nil

		case m.cfg.Keybinds.EncryptMode:
			m.mode = "encrypt"
			m.messages = append(m.messages, m.senderStyle.Render("FunCodex: ")+"mode set to encrypt")
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
			return m, nil

		case m.cfg.Keybinds.Length:
			m.mode = "length"
			m.messages = append(m.messages, m.senderStyle.Render("FunCodex: ")+"mode set to length")
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
    sep := m.cfg.SeparatorCharacter
    if sep == "" {
        sep = "|"
    }
    parts := strings.SplitN(val, sep, 2)
    if len(parts) < 2 {
        response = m.senderStyle.Render("FunCodex: ") + fmt.Sprintf("error: use format 'text %s hidden_text'", sep)
    } else {
        text := strings.TrimSpace(parts[0])
        hiddenText := strings.TrimSpace(parts[1])
        encodedResult := encoder.Encode(text, hiddenText)

        if encodedResult == "" {
            response = m.senderStyle.Render("FunCodex: ") + "Error: [ERROR]"
        } else {
            copyToClipboard(encodedResult)
            response = m.senderStyle.Render("encrypted: ") + fmt.Sprintf(encodedResult)
        }
    }
			case "decrypt":
				decodedResult := decoder.Decode(val)
				response = m.senderStyle.Render("FunCodex: ") + fmt.Sprintf("decrypted: %s", decodedResult)

			case "length":
				encodableChars, maxLen := utils.GetAvailableSpace(val)
				response = m.senderStyle.Render("FunCodex: ") + fmt.Sprintf("encodable chars: %d | max hidden text: %d chars", encodableChars, maxLen)

			default:
				response = m.senderStyle.Render("FunCodex: ") + "Select a mode"
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
	defaultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.cfg.AccentColor)).Bold(true)

	encryptLabel := fmt.Sprintf("%s: encrypt", m.cfg.Keybinds.EncryptMode)
	decryptLabel := fmt.Sprintf("%s: decrypt", m.cfg.Keybinds.DecryptMode)
	lengthLabel := fmt.Sprintf("%s: length", m.cfg.Keybinds.Length)
	quitLabel := "esc: quit"

	if m.mode == "encrypt" {
		encryptLabel = activeStyle.Render(encryptLabel)
	} else {
		encryptLabel = defaultStyle.Render(encryptLabel)
	}

	if m.mode == "decrypt" {
		decryptLabel = activeStyle.Render(decryptLabel)
	} else {
		decryptLabel = defaultStyle.Render(decryptLabel)
	}

	if m.mode == "length" {
		lengthLabel = activeStyle.Render(lengthLabel)
	} else {
		lengthLabel = defaultStyle.Render(lengthLabel)
	}

	footer := fmt.Sprintf(
		"%s • %s • %s • %s",
		encryptLabel,
		decryptLabel,
		lengthLabel,
		defaultStyle.Render(quitLabel),
	)

	return fmt.Sprintf("%s\n\n%s\n%s", m.viewport.View(), m.textarea.View(), footer)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to load config: %v\n", err)
		defaultCfg := config.Default()
		cfg = &defaultCfg
	}

	if err := clipboard.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to init clipboard: %v\n", err)
	}
	// run tui if no flags passed
	if len(os.Args) == 1 {
		p := tea.NewProgram(initialModel(cfg), tea.WithAltScreen())
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
