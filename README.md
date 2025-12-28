# fs - Filesystem Toolkit

A lightweight CLI tool for managing filesystem shortcuts, tags, and quick navigation.

## Features

- **Shortcuts**: Save frequently-used paths with memorable names
- **Tags**: Organize shortcuts by project, category, or context
- **Search**: Find shortcuts by name, path, or tags
- **Quick Navigation**: Jump to any saved location instantly
- **Local Storage**: SQLite database stored in `~/.config/fs/`


## Requirements
- Go 1.21 or higher

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

This will show you the shell function to add to your shell config. Add it to `~/.bashrc` or `~/.zshrc`:
```bash
# sh function
fsg() { cd "$(fs go "$1")"; }

# cmd to add it to ~/.bashrc
echo 'fsg() { cd "$(fs go "$1")"; }' >> ~/.bashrc
```

Then reload your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

## Usage

### CLI Commands

```bash
# Add a shortcut
fs add ~/projects/homelab homelab
fs add /var/log/nginx nginx-logs

# List all shortcuts
fs list

# Jump to a shortcut
fsg homelab           # Using the shell function
cd $(fs go homelab)   # Alternative syntax

# Remove a shortcut
fs rm homelab

# Tag shortcuts
fs tag homelab work docker
fs tag nginx-logs work debugging

# Search shortcuts
fs search docker
fs search --tag work

# Find files in shortcuts
fs find "docker-compose.yml"
fs find "*.log" --shortcut nginx-logs
```

## Appendix

### Project Structure
```
fs/
├── cmd/fs/           # Main entry point
├── internal/
│   ├── storage/      # SQLite Database layer
│   ├── shortcut/     # Business logic
│   └── fileops/      # File operations
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

**Note:** This is a personal project built to learn Go and scratch a personal itch. Inspired by tools like `z`, `autojump`, and `ranger`.