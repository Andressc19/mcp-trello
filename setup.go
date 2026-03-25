package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	trelloBlue  = lipgloss.Color("25")
	trelloLight = lipgloss.Color("86")
	gray        = lipgloss.Color("241")
	green       = lipgloss.Color("84")
	red         = lipgloss.Color("204")
	white       = lipgloss.Color("255")

	titleStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(trelloBlue).
			Bold(true).
			Padding(0, 1)

	normalText = lipgloss.NewStyle().
			Foreground(white)

	dimText = lipgloss.NewStyle().
		Foreground(gray)

	selectedText = lipgloss.NewStyle().
			Foreground(trelloBlue).
			Bold(true)

	errorText = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	successText = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	inputLabel = lipgloss.NewStyle().
			Foreground(trelloLight).
			Bold(true)

	helpText = lipgloss.NewStyle().
			Foreground(gray)
)

type model struct {
	step        int
	choices     []string
	cursor      int
	apiKey      string
	token       string
	configPath  string
	opencodeCfg string
	err         error
	// Text inputs
	apiKeyInput *textinput.Model
	tokenInput  *textinput.Model
	focusField  int // 0 = none, 1 = apiKey, 2 = token
}

func initialModel() model {
	apiInput := textinput.New()
	apiInput.Placeholder = "Paste your Trello API Key here"
	apiInput.Focus()
	apiInput.Prompt = ""

	tokenInput := textinput.New()
	tokenInput.Placeholder = "Paste your Trello Token here"
	tokenInput.Prompt = ""

	return model{
		step: 0,
		choices: []string{
			"Configure credentials",
			"Verify current configuration",
			"Uninstall / Remove configuration",
			"Exit",
		},
		cursor:      0,
		apiKeyInput: &apiInput,
		tokenInput:  &tokenInput,
		focusField:  1,
	}
}

// Main entry point for setup command
func RunSetup() error {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	// Clear screen on start and hide cursor
	return tea.Sequence(
		tea.ClearScreen,
		tea.HideCursor,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle text inputs when on input steps
	if m.step == 1 && m.apiKeyInput != nil {
		// Check for Enter key first to move to next step
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			m.apiKey = m.apiKeyInput.Value()
			m.step = 2
			// Focus token input when switching to step 2
			if m.tokenInput != nil {
				m.tokenInput.Focus()
			}
			return m, nil
		}
		var cmd tea.Cmd
		*m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
		return m, cmd
	}
	if m.step == 2 && m.tokenInput != nil {
		// Check for Enter key first to move to next step
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			m.token = m.tokenInput.Value()
			m.step = 3
			return m, nil
		}
		var cmd tea.Cmd
		*m.tokenInput, cmd = m.tokenInput.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if m.step == 0 || m.step == 10 {
				return m, tea.Quit
			}
			m.step = 0
			m.err = nil
			return m, nil
		case "b":
			// Go back
			if m.step > 1 && m.step < 10 {
				m.step = 0
				m.err = nil
			}
			if m.step == 20 || m.step == 21 {
				m.step = 0
				m.err = nil
			}
			return m, nil
		case "o":
			// Open browser to Trello
			if m.step == 1 || m.step == 2 {
				return m, func() tea.Msg {
					openBrowser("https://trello.com/power-ups/admin")
					return nil
				}
			}
		case "up", "k":
			if m.step == 0 && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.step == 0 && m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.step == 0 {
				switch m.cursor {
				case 0:
					m.step = 1 // Start configuration
				case 1:
					m.step = 10 // Verify config
				case 2:
					m.step = 20 // Uninstall
				case 3:
					return m, tea.Quit
				}
			} else if m.step == 20 {
				// Perform uninstall
				if err := uninstall(); err != nil {
					m.err = err
				} else {
					m.step = 21 // Uninstall complete
				}
			} else if m.step == 1 {
				// Save API Key and move to token
				if m.apiKeyInput != nil {
					m.apiKey = m.apiKeyInput.Value()
				}
				m.step = 2
			} else if m.step == 2 {
				// Save Token and continue
				if m.tokenInput != nil {
					m.token = m.tokenInput.Value()
				}
				m.step = 3
			} else if m.step == 3 {
				// Save credentials and continue
				if err := saveCredentials(m.apiKey, m.token); err != nil {
					m.err = err
					return m, nil
				}
				m.step = 4
			} else if m.step == 4 {
				// Add to OpenCode
				if err := addToOpenCodeConfig(); err != nil {
					m.err = err
					return m, nil
				}
				m.step = 5
			} else if m.step == 5 {
				// Done
				return m, tea.Quit
			} else if m.step == 10 {
				// Exit verify view
				m.step = 0
			} else if m.step == 21 {
				// Exit uninstall done
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string

	switch m.step {
	case 0:
		s = mainMenuView(m)
	case 1:
		s = apiKeyInputView(m)
	case 2:
		s = tokenInputView(m)
	case 3:
		s = saveConfigView(m)
	case 4:
		s = addToOpenCodeView(m)
	case 5:
		s = doneView(m)
	case 10:
		s = verifyConfigView(m)
	case 20:
		s = uninstallView(m)
	case 21:
		s = uninstallDoneView(m)
	}

	return s
}

func mainMenuView(m model) string {
	s := titleStyle.Render(" mcp-trello Setup ") + "\n\n"
	s += normalText.Render("Welcome! This wizard will help you configure") + "\n"
	s += normalText.Render("mcp-trello to work with OpenCode.") + "\n\n"
	s += normalText.Render("What would you like to do?") + "\n\n"

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = ">> "
			s += selectedText.Render(cursor+choice) + "\n"
		} else {
			s += cursor + choice + "\n"
		}
	}

	s += "\n" + helpText.Render("up/k down/j: navigate | enter: select | q: quit")
	return s
}

func maskInput(s string) string {
	if len(s) == 0 {
		return "(empty)"
	}
	if len(s) <= 8 {
		return "********"
	}
	return s[:4] + "********" + s[len(s)-4:]
}

func apiKeyInputView(m model) string {
	s := titleStyle.Render(" Step 1: API Key ") + "\n\n"
	s += normalText.Render("Enter your Trello API Key.") + "\n"
	s += dimText.Render("Press 'o' to open Trello if you don't have one.") + "\n\n"

	if m.apiKeyInput != nil {
		currentValue := m.apiKeyInput.Value()
		displayValue := maskInput(currentValue)
		s += inputLabel.Render("API Key: ") + displayValue + "\n"
	} else {
		s += inputLabel.Render("API Key: ") + "(error)\n"
	}

	s += "\n" + helpText.Render("enter: continue | b: back | o: open browser")
	return s
}

func tokenInputView(m model) string {
	s := titleStyle.Render(" Step 2: Token ") + "\n\n"
	s += normalText.Render("Enter your Trello Token.") + "\n"
	s += dimText.Render("Press 'o' to open Trello.") + "\n\n"

	if m.tokenInput != nil {
		currentValue := m.tokenInput.Value()
		displayValue := maskInput(currentValue)
		s += inputLabel.Render("Token: ") + displayValue + "\n"
	} else {
		s += inputLabel.Render("Token: ") + "(error)\n"
	}

	s += "\n" + helpText.Render("enter: continue | b: back | o: open browser")
	return s
}

func saveConfigView(m model) string {
	s := titleStyle.Render(" Step 3: Save ") + "\n\n"
	s += normalText.Render("Ready to save credentials:") + "\n\n"
	s += "  " + inputLabel.Render("API Key: ") + maskString(m.apiKey) + "\n"
	s += "  " + inputLabel.Render("Token: ") + maskString(m.token) + "\n\n"

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "opencode", "mcp-trello.json")
	s += dimText.Render("Save to: ") + configPath + "\n\n"

	s += helpText.Render("enter: save | b: back")
	return s
}

func addToOpenCodeView(m model) string {
	s := titleStyle.Render(" Step 4: Add to OpenCode ") + "\n\n"

	home, _ := os.UserHomeDir()
	opencodePath := filepath.Join(home, ".config", "opencode", "opencode.json")

	s += normalText.Render("Add mcp-trello server to OpenCode config?") + "\n\n"
	s += dimText.Render("Config: ") + opencodePath + "\n\n"

	if m.err != nil {
		s += errorText.Render("Error: "+m.err.Error()) + "\n\n"
	}

	s += helpText.Render("enter: add automatically | b: skip")
	return s
}

func verifyConfigView(m model) string {
	s := titleStyle.Render(" Verify Configuration ") + "\n\n"

	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")

	s += dimText.Render("Environment variables:") + "\n"
	s += "  TRELLO_API_KEY: " + getEnvStatus(apiKey) + "\n"
	s += "  TRELLO_TOKEN: " + getEnvStatus(token) + "\n\n"

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "opencode", "mcp-trello.json")
	if data, err := os.ReadFile(configPath); err == nil {
		var creds Credentials
		if json.Unmarshal(data, &creds) == nil {
			s += dimText.Render("Config file:") + "\n"
			s += "  Found: " + successText.Render("yes") + "\n"
			s += "  API Key: " + getEnvStatus(creds.APIKey) + "\n"
			s += "  Token: " + getEnvStatus(creds.Token) + "\n"
		}
	} else {
		s += dimText.Render("Config file: ") + "not found\n"
	}

	s += "\n" + helpText.Render("press any key to go back...")
	return s
}

func doneView(m model) string {
	s := titleStyle.Render(" Setup Complete! ") + "\n\n"
	s += successText.Render("+ Configuration saved") + "\n"
	s += successText.Render("+ Added to OpenCode") + "\n\n"
	s += normalText.Render("You can now use mcp-trello!") + "\n"
	s += "Run: " + inputLabel.Render("opencode mcp list") + "\n\n"
	s += helpText.Render("press enter to exit")
	return s
}

func maskString(s string) string {
	if len(s) < 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

func getEnvStatus(val string) string {
	if val != "" {
		return successText.Render("set")
	}
	return dimText.Render("not set")
}

func uninstallView(m model) string {
	s := titleStyle.Render(" Uninstall ") + "\n\n"
	s += normalText.Render("This will remove:") + "\n"
	s += "  - Credentials file: ~/.config/opencode/mcp-trello.json\n"
	s += "  - MCP server from OpenCode config\n\n"

	if m.err != nil {
		s += errorText.Render("Error: "+m.err.Error()) + "\n\n"
	}

	s += dimText.Render("Press ") + normalText.Render("enter") + dimText.Render(" to confirm uninstall")
	s += dimText.Render("\nor ") + normalText.Render("b") + dimText.Render(" to go back")
	return s
}

func uninstallDoneView(m model) string {
	s := titleStyle.Render(" Uninstall Complete ") + "\n\n"
	s += successText.Render("- Credentials file removed") + "\n"
	s += successText.Render("- OpenCode config updated") + "\n\n"
	s += normalText.Render("mcp-trello has been uninstalled.")
	s += helpText.Render("\npress enter to exit")
	return s
}

func uninstall() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	// Remove credentials file
	credsPath := filepath.Join(home, ".config", "opencode", "mcp-trello.json")
	if err := os.Remove(credsPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not remove credentials file: %w", err)
	}

	// Remove from OpenCode config
	opencodePath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if data, err := os.ReadFile(opencodePath); err == nil {
		var config map[string]interface{}
		if json.Unmarshal(data, &config) == nil {
			if mcp, ok := config["mcp"].(map[string]interface{}); ok {
				delete(mcp, "trello")
				if len(mcp) == 0 {
					delete(config, "mcp")
				}
				if newData, err := json.MarshalIndent(config, "", "  "); err == nil {
					os.WriteFile(opencodePath, newData, 0644)
				}
			}
		}
	}

	return nil
}

// Helper functions for file operations
func saveCredentials(apiKey, token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "opencode")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "mcp-trello.json")
	creds := Credentials{
		APIKey: apiKey,
		Token:  token,
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal credentials: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

func addToOpenCodeConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")

	var config map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &config)
	} else {
		config = make(map[string]interface{})
	}

	// Add mcp section if not exists
	if _, ok := config["mcp"]; !ok {
		config["mcp"] = make(map[string]interface{})
	}

	mcp := config["mcp"].(map[string]interface{})
	mcp["trello"] = map[string]interface{}{
		"command": []string{"mcp-trello"},
		"enabled": true,
		"type":    "local",
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
