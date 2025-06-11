# Currency Converter CLI

A simple, fast, and elegant command-line tool for converting currencies using the latest exchange rates.

This tool has two modes:
1.  **Interactive TUI Mode**: A beautiful, full-screen interface for selecting currencies and amounts.
2.  **Direct Command Mode**: A quick way to get a conversion without leaving your current command line.

![screenshot](https://user-images.githubusercontent.com/4254531/182820427-3558f4f6-b1b4-4786-8843-779836125028.png)

## Features

* Fetches the latest exchange rates from [ExchangeRate-API.com](https://www.exchangerate-api.com).
* Fuzzy search for currencies by code (e.g., "USD") or name (e.g., "dollar").
* Supports all major world currencies.
* Polished and user-friendly terminal interface.

## Prerequisites

* Go 1.18 or later.
* An API key from [ExchangeRate-API.com](https://www.exchangerate-api.com). The free tier is sufficient.

## Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/BobOfTheHawk/currency-cli.git
    cd currency-cli
    ```

2.  **Set your API Key:**
    Open the `api.go` file and replace the placeholder `YOUR_API_KEY` with your actual key from ExchangeRate-API.
    ```go
    // In api.go
    const apiKey = "YOUR_API_KEY" // <-- EDIT THIS LINE
    ```

3.  **Build and Install:**
    Use the provided `Makefile` to build the binary and install it to `/usr/local/bin`. You may be prompted for your password as this requires administrator privileges.
    ```bash
    make install
    ```
    This command compiles the Go source code and moves the executable to a directory that is in your system's `PATH`, so you can run it from anywhere.

## Usage

Once installed, you can use the `currency-cli` command.

#### Interactive Mode

For the full-screen, user-friendly interface, simply run the command without any arguments:

```bash
currency-cli
```
Follow the on-screen prompts to select your 'from' and 'to' currencies and enter the amount.

#### Direct Command Mode

For quick, one-off conversions, you can provide the arguments directly on the command line.

**Syntax:**
```bash
currency-cli <amount> <from_currency> <to_currency>
```

**Example:**
```bash
$ currency-cli 10.50 USD UZS
Converting 10.50 USD to UZS...

Result: 10.50 USD is equal to 132483.75 UZS
```

#### Getting Help

To see the help message with usage instructions, use the `--help` flag:
```bash
currency-cli --help
```

## Uninstallation

To remove the `currency-cli` command from your system, you can use the `uninstall` command from the project directory:

```bash
make uninstall
```

This will remove the binary from `/usr/local/bin`.
