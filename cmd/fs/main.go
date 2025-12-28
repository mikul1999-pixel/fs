package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/mikul1999-pixel/fs/internal/storage"
	"github.com/mikul1999-pixel/fs/pkg/config"
)

var store storage.Storage

var rootCmd = &cobra.Command{
	Use:   "fs",
	Short: "Filesystem shortcut toolkit",
	Long:  `A CLI tool for managing filesystem shortcuts, tags, and quick navigation`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Show setup instructions for shell integration",
	Run: func(cmd *cobra.Command, args []string) {
		shell := detectShell()
		fmt.Printf("To enable quick directory switching, run this command:\n\n")
		
		if shell == ".bashrc" || shell == ".zshrc" {
			fmt.Printf("    echo 'f() { cd \"$(fs go \"$1\")\"; }' >> ~/%s\n", shell)
		} else {
			fmt.Printf("    echo 'function f; cd (fs go $argv[1]); end' >> ~/%s\n", shell)
		}
		
		fmt.Printf("\nThen reload your shell:\n\n")
		fmt.Printf("    source ~/%s\n", shell)
		fmt.Println("\nAfter that, you can use: f <shortcut-name>")
		fmt.Println("\nExample:")
		fmt.Println("  fs add ~/projects/homelab homelab")
		fmt.Println("  f homelab")
	},
}

func detectShell() string {
    shell := os.Getenv("SHELL")
    if strings.Contains(shell, "zsh") {
        return ".zshrc"
    }
    return ".bashrc"
}

var addCmd = &cobra.Command{
	Use:   "add <path> <name>",
	Short: "Add a new shortcut",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		name := args[1]

		// Expand path (handle ~ and relative paths)
		absPath, err := expandPath(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid path: %v\n", err)
			os.Exit(1)
		}

		// Validate path exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: path does not exist: %s\n", absPath)
			os.Exit(1)
		}

		// Add to database
		if err := store.AddShortcut(name, absPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added shortcut: %s -> %s\n", name, absPath)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all shortcuts",
	Run: func(cmd *cobra.Command, args []string) {
		shortcuts, err := store.ListShortcuts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(shortcuts) == 0 {
			fmt.Println("No shortcuts found. Add one with: fs add <path> <name>")
			return
		}

		fmt.Println("Shortcuts:")
		for _, sc := range shortcuts {
			tagStr := ""
			if len(sc.Tags) > 0 {
				tagStr = fmt.Sprintf(" [%s]", strings.Join(sc.Tags, ", "))
			}
			fmt.Printf("  %s -> %s%s\n", sc.Name, sc.Path, tagStr)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:     "delete <name>",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete a shortcut",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := store.DeleteShortcut(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deleted shortcut: %s\n", name)
	},
}

var goCmd = &cobra.Command{
	Use:   "go <name>",
	Short: "Get path for a shortcut",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		sc, err := store.GetShortcut(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Print the path. Use with cd $(fs go <name>)
		fmt.Println(sc.Path)
	},
}

func expandPath(path string) (string, error) {
	// Handle ~ for home directory
	if path[:1] == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[1:])
	}

	// Convert to absolute path
	return filepath.Abs(path)
}

var peekCmd = &cobra.Command{
	Use:   "peek <name>",
	Aliases: []string{"ls"},
	Short: "Preview the contents of a shortcut location",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		sc, err := store.GetShortcut(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Contents of %s (%s):\n\n", name, sc.Path)

		// Execute ls -lah on the shortcut path
		lsCmd := exec.Command("ls", "-lah", sc.Path)
		lsCmd.Stdout = os.Stdout
		lsCmd.Stderr = os.Stderr

		if err := lsCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running ls: %v\n", err)
			os.Exit(1)
		}
	},
}

var tagCmd = &cobra.Command{
	Use:   "tag <shortcut> <tags...>",
	Short: "Add tags to a shortcut",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		shortcutName := args[0]
		tags := args[1:]

		if err := store.AddTags(shortcutName, tags); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added tags to %s: %s\n", shortcutName, strings.Join(tags, ", "))
	},
}

var untagCmd = &cobra.Command{
	Use:   "untag <shortcut> <tags...>",
	Short: "Remove tags from a shortcut",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		shortcutName := args[0]
		tags := args[1:]

		if err := store.RemoveTags(shortcutName, tags); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Removed tags from %s: %s\n", shortcutName, strings.Join(tags, ", "))
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Aliases: []string{"find"},
	Short: "Search shortcuts by name, path, or tags",
	Run: func(cmd *cobra.Command, args []string) {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		tags, _ := cmd.Flags().GetStringSlice("tag")

		shortcuts, err := store.SearchShortcuts(query, tags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(shortcuts) == 0 {
			fmt.Println("No shortcuts found.")
			return
		}

		fmt.Println("Shortcuts:")
		for _, sc := range shortcuts {
			tagStr := ""
			if len(sc.Tags) > 0 {
				tagStr = fmt.Sprintf(" [%s]", strings.Join(sc.Tags, ", "))
			}
			fmt.Printf("  %s -> %s%s\n", sc.Name, sc.Path, tagStr)
		}
	},
}

func init() {
	searchCmd.Flags().StringSliceP("tag", "t", []string{}, "Filter by tags") // Add flags to search before adding it to root

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(goCmd)
	rootCmd.AddCommand(peekCmd)
	rootCmd.AddCommand(tagCmd)
	rootCmd.AddCommand(untagCmd)
	rootCmd.AddCommand(searchCmd)
}

func main() {
	// Initialize storage
	dbPath := config.GetDBPath()
	var err error
	store, err = storage.NewSQLiteStorage(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}