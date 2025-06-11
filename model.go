package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appState int

const (
	stateLoading appState = iota
	stateFromCurrency
	stateToCurrency
	stateAmount
	stateConverting
	stateResult
)

type ratesMsg map[string]float64
type namesMsg map[string]string
type errorMsg error

type model struct {
	state         appState
	list          list.Model
	amountInput   textinput.Model
	spinner       spinner.Model
	amount        float64
	fromCurrency  string
	toCurrency    string
	currencyNames map[string]string
	result        string
	err           error
	width, height int
}

// item now holds both code and name for the list.
type item struct {
	code, name string
}

// These methods satisfy the list.Item interface.
func (i item) FilterValue() string { return fmt.Sprintf("%s %s", i.code, i.name) }
func (i item) Title() string       { return i.code }
func (i item) Description() string { return i.name }

func initialModel() model {
	// Spinner for loading states
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Text input for the amount
	ti := textinput.New()
	ti.Placeholder = "100.00"
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = 20

	// List for currency selection
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("229")).BorderLeftForeground(lipgloss.Color("229"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Copy().Foreground(lipgloss.Color("240"))

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Loading currencies..."
	l.Styles.Title = TitleStyle
	l.Styles.PaginationStyle = PaginationStyle
	l.Styles.HelpStyle = HelpStyle

	return model{
		state:       stateLoading,
		spinner:     s,
		amountInput: ti,
		list:        l,
	}
}

func (m model) Init() tea.Cmd {
	// Initial command to fetch currency names
	return func() tea.Msg {
		names, err := getCurrencyNames()
		if err != nil {
			return errorMsg(err)
		}
		return namesMsg(names)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// When the window is resized, update the list's dimensions
		m.list.SetSize(m.width, m.height)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// If on result screen, any key should restart the app
		if m.state == stateResult || m.state == stateLoading && m.err != nil {
			// *** THIS IS THE CRITICAL FIX ***
			// Create a new model, but preserve the window size from the previous model.
			newModel := initialModel()
			newModel.width = m.width
			newModel.height = m.height
			newModel.list.SetSize(m.width, m.height) // Set the size on the new list explicitly.
			return newModel, newModel.Init()
		}

	// This message is received after fetching currency names
	case namesMsg:
		m.currencyNames = msg
		codes := make([]string, 0, len(msg))
		for code := range msg {
			codes = append(codes, code)
		}
		sort.Strings(codes)
		items := make([]list.Item, len(codes))
		for i, code := range codes {
			items[i] = item{code: code, name: msg[code]}
		}
		m.list.SetItems(items)
		m.list.Title = "Convert FROM which currency?"
		m.state = stateFromCurrency
		return m, nil

	// This message is received after fetching conversion rates
	case ratesMsg:
		targetRate, ok := msg[m.toCurrency]
		if !ok {
			m.err = fmt.Errorf("conversion rate for %s not found", m.toCurrency)
		} else {
			converted := m.amount * targetRate
			fromName := m.currencyNames[m.fromCurrency]
			toName := m.currencyNames[m.toCurrency]
			m.result = fmt.Sprintf("%.2f %s (%s) is %.2f %s (%s)", m.amount, m.fromCurrency, fromName, converted, m.toCurrency, toName)
		}
		m.state = stateResult
		return m, nil

	case errorMsg:
		m.err = msg
		m.state = stateResult
		return m, nil
	}

	// Process updates based on the current application state
	switch m.state {
	case stateFromCurrency, stateToCurrency:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)

		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			if i, ok := m.list.SelectedItem().(item); ok {
				if m.state == stateFromCurrency {
					m.fromCurrency = i.code
					m.state = stateToCurrency
					m.list.Title = fmt.Sprintf("Convert from %s TO which currency?", m.fromCurrency)
					m.list.ResetFilter()
				} else {
					m.toCurrency = i.code
					m.state = stateAmount
				}
			}
		}

	case stateAmount:
		m.amountInput, cmd = m.amountInput.Update(msg)
		cmds = append(cmds, cmd)
		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			val, err := strconv.ParseFloat(m.amountInput.Value(), 64)
			if err == nil {
				m.amount = val
				m.state = stateConverting
				// After getting amount, fetch the rates
				return m, func() tea.Msg {
					rates, err := getRates(m.fromCurrency)
					if err != nil {
						return errorMsg(err)
					}
					return ratesMsg(rates)
				}
			}
		}

	case stateLoading, stateConverting:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Handle error view first
	if m.err != nil {
		content := ErrorStyle.Render(fmt.Sprintf("\nAn error occurred: %v", m.err))
		content += "\n\nPress any key to restart."
		// Use Place to center the error message
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, AppStyle.Render(content))
	}

	// If the state is one of the list views, render the list directly.
	if m.state == stateFromCurrency || m.state == stateToCurrency {
		return m.list.View()
	}

	// For all other states, we create the content and then center it.
	var content string
	switch m.state {
	case stateAmount:
		var sb strings.Builder
		fromName := m.currencyNames[m.fromCurrency]
		toName := m.currencyNames[m.toCurrency]
		sb.WriteString(TitleStyle.Render(fmt.Sprintf("From %s (%s) to %s (%s)", m.fromCurrency, fromName, m.toCurrency, toName)))
		sb.WriteString("\n\nEnter amount:\n")
		sb.WriteString(m.amountInput.View())
		content = sb.String()

	case stateResult:
		content = TitleStyle.Render(m.result)
		content += "\n\nPress any key to start over."

	case stateLoading, stateConverting:
		content = fmt.Sprintf("%s Working...", m.spinner.View())
	}

	// Place the simple content views in the center of the screen.
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, AppStyle.Render(content))
}
