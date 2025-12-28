.PHONY: build install uninstall clean test run help

BINARY_NAME=fs
INSTALL_PATH=$(HOME)/.local/bin
GO=go

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	$(GO) build -o $(BINARY_NAME) cmd/fs/main.go
	@echo "Built $(BINARY_NAME)"

install: build ## Build and install to ~/.local/bin
	@mkdir -p $(INSTALL_PATH)
	@cp $(BINARY_NAME) $(INSTALL_PATH)/
	@echo ""
	@echo "fs installed to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo ""
	@echo "Make sure $(INSTALL_PATH) is in your PATH."
	@echo "If not, add this to your ~/.bashrc or ~/.zshrc:"
	@echo "    export PATH=\"\$$HOME/.local/bin:\$$PATH\""
	@echo ""
	@echo "Next steps:"
	@echo "  1. Run: source ~/.bashrc  (or restart your terminal)"
	@echo "  2. Run: fs init           (to set up the shell function)"
	@echo "  3. Try: fs add ~/projects/myproject proj"
	@echo ""

uninstall: ## Remove the installed binary
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"

clean: ## Remove built binary
	@rm -f $(BINARY_NAME)
	@echo "Cleaned build artifacts"

test: ## Run tests
	$(GO) test -v ./...

run: ## Run without installing
	$(GO) run cmd/fs/main.go

.DEFAULT_GOAL := help