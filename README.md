# fs - Filesystem Toolkit

A lightweight CLI tool for managing path shortcuts, tags, and quick navigation.

**Note:** This is a personal project built to learn Go and scratch a personal itch. Inspired by tools like `zoxide` and `autojump`

## Features

- **Shortcuts**: Save frequently-used paths with memorable names. Then jump with ```f <shortcut>```
- **Tags**: Organize shortcuts by project, category, or context
- **Search**: Find shortcuts by name, path, or tags
- **Peek**: Preview directory contents before jumping
- **Local Storage**: SQLite database stored in `~/.config/fs/`

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
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

### Initialize commands

After installation, add the init command to your shell config. Add to `~/.bashrc` or `~/.zshrc`:
```bash
# Add function to your shell config
eval "$(fs init)" 

# Or run a cmd to add it
echo 'eval "$(fs init)"' >> ~/.bashrc

# Reload your shell
source ~/.bashrc
```
For zsh, replace `.bashrc` with `.zshrc`. <br> <br>
This inititalizes cd shortcuts `f()` and `ff()` *(see CLI usage below)*. Or you can create your own aliases
```bash
# customize function names:
eval "$(fs init go)"           # creates go() and ff()
eval "$(fs init go search)"    # creates go() and search()
```

## Usage

### CLI Commands

```bash
# Add a shortcut (path defaults to cwd)
fs add <name>
fs add <name> <path>

# List all shortcuts
fs list

# Edit a shortcut
fs edit-path <name> <new-path>
fs edit-name <name> <new-name>

# Remove a shortcut
fs rm <name>

# Preview directory contents
fs peek <name>

# Add tags to a shortcut
fs tag <name> <tag1> <tag2> ...

# Remove tags from a shortcut (no tag defaults to remove all)
fs untag <name>
fs untag <name> <tag1> <tag2> ...

# Jump to a shortcut
f <name>

# Search and jump
ff <like:name-or-path>
ff <like:name-or-path> -t <tag1> -t <tag2> ....
ff --tag <tag1> --tag <tag2> -o and


# Example workflow
fs add cli
fs tag cli proj
fs tag cli go
ff --tag proj # Show all personal projects
ff --tag go   # Show all Go repos
```

## Appendix

**Note:** This is a personal project. Features and functionality may change.

<br>

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
