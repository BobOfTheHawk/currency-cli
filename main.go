package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

const helpText = `
Currency Converter CLI

This tool converts an amount from one currency to another using the latest exchange rates.

USAGE (Interactive TUI):
  currency-converter

USAGE (Direct Command):
  currency-converter <amount> <from_currency> <to_currency>
  Example: currency-converter 10.50 USD UZS

FLAGS:
  --help    Show this help message
`

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// Using WithAltScreen gives our app its own screen buffer, which is
		// the standard and most stable way to run full-screen TUIs.
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("An error occurred: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if len(args) == 1 && args[0] == "--help" {
		fmt.Println(helpText)
		return
	}

	if len(args) == 3 {
		amountStr, fromCurrency, toCurrency := args[0], args[1], args[2]

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			fmt.Printf("Error: Invalid amount '%s'. Please provide a number.\n", amountStr)
			os.Exit(1)
		}

		fmt.Printf("Converting %.2f %s to %s...\n", amount, fromCurrency, toCurrency)

		rates, err := getRates(fromCurrency)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		targetRate, ok := rates[toCurrency]
		if !ok {
			fmt.Printf("Error: The target currency '%s' was not found.\n", toCurrency)
			os.Exit(1)
		}

		convertedAmount := amount * targetRate
		fmt.Printf("\nResult: %.2f %s is equal to %.2f %s\n", amount, fromCurrency, convertedAmount, toCurrency)
		return
	}

	fmt.Println(ErrorStyle.Render("Error: Invalid arguments."))
	fmt.Println(helpText)
	os.Exit(1)
}
