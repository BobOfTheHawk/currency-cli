# Makefile for the Currency Converter CLI

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOTIDY=$(GOMOD) tidy

# Binary name
BINARY_NAME=currency-cli
INSTALL_PATH=/usr/local/bin

# Default target executed when you just run `make`
all: build

# Builds the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOTIDY)
	$(GOBUILD) -o $(BINARY_NAME) .

# Installs the binary to a system-wide location
# Requires sudo privileges to write to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully."
	@echo "You can now run 'currency-cli' from anywhere."

# Uninstalls the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_PATH)..."
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled."

# Cleans up the build artifacts
clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# A phony target is one that is not a file. It's a way to tell `make`
# that these commands do not produce a file with the same name.
.PHONY: all build install uninstall clean
