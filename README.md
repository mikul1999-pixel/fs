# fs - Filesystem Toolkit

A lightweight CLI tool for managing filesystem shortcuts, tags, and quick navigation.

## Features

- **Shortcuts**: Save frequently-used paths with memorable names
- **Tags**: Organize shortcuts by project, category, or context
- **Search**: Find shortcuts by name, path, or tags
- **Quick Navigation**: Jump to any saved location instantly
- **Peek**: Preview directory contents before jumping
- **Local Storage**: SQLite database stored in `~/.config/fs/`


## Requirements
- Go 1.21 or higher
- Bash or Zsh shell

## Installation

### Setup Project
```bash
git clone https://github.com/mikul1999-pixel/fs.git
cd fs
make install
```

This will:
1. Build the binary
2. Install it to `~/.local/bin/fs`
3. Show setup instructions

### Verify installation
```bash
fs --help
```

If you get "command not found", make sure `~/.local/bin` is in your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

Add that line to your `~/.bashrc` or `~/.zshrc` to make it permanent.

### Initialize commands

After installation, run the init command to set up quick navigation:
```bash
fs init
```

This will show you the shell functions to add to your shell config. Add to `~/.bashrc` or `~/.zshrc`:
```bash
# Add functions to your shell config
echo 'f() { cd "$(fs go "$1")"; }' >> ~/.bashrc
echo 'ff() { local path=$(fs find "$@" </dev/tty); [ $? -eq 0 ] && [ -n "$path" ] && cd "$path"; }' >> ~/.bashrc

# Reload your shell
source ~/.bashrc
```
For zsh, replace `.bashrc` with `.zshrc`.

## Usage

### CLI Commands

```bash
# Add a shortcut
fs add <path> <shortcut>

# List all shortcuts
fs list

# Edit a shortcut
fs edit-path <shortcut> <new-path>
fs edit-name <shortcut> <new-name>

# Remove a shortcut
fs rm <shortcut>

# Preview directory contents
fs peek <shortcut>

# Add tags to a shortcut
fs tag <shortcut> <tag1> <tag2> ...

# Delete a tag
fs untag <shortcut>
fs untag <shortcut> <tag1> <tag2> ...

# Jump to a shortcut
f <shortcut>

# Search and jump
ff <like:shortcut-or-path>
ff <like:shortcut-or-path> --tag <tag1> <tag2> ...
ff --tag <tag1> <tag2> ...
```

## Appendix

### Project Structure
```
fs/
├── cmd/fs/           # Main CLI application
├── internal/
│   ├── storage/      # SQLite Database layer    
│   └── ui/           # Bubbletea TUI components
├── pkg/config/       # Config
└── Makefile          # Build automation
```

### Building
```bash
make build    # Build the binary
make test     # Run tests
make clean    # Remove build artifacts
make run      # Run without installing
```

### Database Management

- **Database location**: `~/.config/fs/shortcuts.db`
- **Data format**: SQLite

To reset everything:
```bash
rm -rf ~/.config/fs/
```


<br>

---

**Note:** This is a personal project built to learn Go and scratch a personal itch. Inspired by tools like `z`, `autojump`, and `ranger`. And built using `Cobra (CLI framework)` and `Bubbletea (TUI)`