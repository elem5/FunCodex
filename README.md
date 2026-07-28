# FunCodex
Simple text encoder made in Go *(a.k.a. best scripting language on Earth)*

## Main concepts
This program is able to hide a small binary code inside a normal text alternating similar characters (like the latin a and the cyrillic a) between them.

## Disclaimer
I'm making this small project just to try NeoVim and to see if I like its workflow, I know that this program is impractical for many reasons and should not be used in a critical context where messages should be truly hidden, also, its biggest flaw is the **limited space** of binary code you can store in an actual sentence, which makes it work *badly* in small text messages.

This doesn't mean that it's *useless* though, I mean look at it, if you want to tell someone something but you're not brave enough to do so, just send them a FunCodex message and pretend that you did something ;)

Apart from that, I don't even know if there's already something like this, if so, let me know, I want to specify that I'm not trying to copy anyone, I had this idea on my own but I also acknowledge the fact that the web is vast and it's easy to have similar ideas.

---

## Installation

### Prerequisites

- **Go** version 1.20 or later installed on your system.

### Building from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/Danyiyk/FunCodex.git
   cd FunCodex
   ```

2. Build the executable:

   ```bash
   go build -o funcodex ./main.go
   ```

---

## Configuration

On the first run, **FunCodex** automatically creates a default configuration file (`config.json`) in your system's user configuration directory (for example, `~/.config/funcodex/config.json` on Linux/macOS or `%APPDATA%\funcodex\config.json` on Windows).

### Example `config.json`

```json
{
    "accent_color": "#5865F2",
    "show_ascii": true,
    "separator_character": "|",
    "keybinds": {
        "encrypt_mode": "ctrl+e",
        "decrypt_mode": "ctrl+d",
        "length_mode": "ctrl+l"
    }
}
```

---

## Usage

### Terminal User Interface (TUI)

To launch the interactive interface, run the executable without any arguments:

```bash
./funcodex
```

#### Default Keybindings

- **Ctrl + E**: Switch to **Encrypt** mode
- **Ctrl + D**: Switch to **Decrypt** mode
- **Ctrl + L**: Switch to **Length Analysis** mode
- **Esc** or **Ctrl + C**: Quit the application

#### Using the TUI

1. Select a mode using the keyboard shortcuts.
2. Enter your text into the input field.

- In **Encrypt** mode, separate the cover text and secret message using the configured separator (default: `|`):

  ```text
  Visible cover text | Hidden secret message
  ```

  The encoded output is automatically copied to your clipboard.

- In **Decrypt** mode, paste the encoded text.

- In **Length Analysis** mode, enter text to calculate the maximum number of hidden characters it can contain.

---

### Command Line Interface (CLI)

You can also use FunCodex directly from the terminal by passing command-line flags.

#### Encrypt a message (`-c`)

Hide a secret message inside visible cover text:

```bash
./funcodex -c "This is the cover text" "This is the secret message"
```

#### Decrypt a message (`-d`)

Extract the hidden message from encoded text:

```bash
./funcodex -d "ReceivedEncodedText..."
```

#### Analyze available capacity (`-l`)

Calculate how many characters can be hidden within a given text:

```bash
./funcodex -l "Text to analyze"
```

#### Help

Display the available commands:

```bash
./funcodex -h
```
