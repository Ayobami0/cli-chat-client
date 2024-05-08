package ui

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func login(credential map[string]string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{
			Timeout: 10 * time.Second,
		}
		res, err := c.Get("https://google.com")
		if err != nil {
			log.Println(err)
			return errMsg{err}
		}
		defer res.Body.Close()
		log.Println(res.Status, credential["password"], credential["username"])
		return statusMsg{sType: STATUS_LOGIN}
	}
}

type loginModel struct {
	username        string
	storedPass      string
	password        textinput.Model
	isLoading       bool
	isLoggedIn      bool
	spinner         spinner.Model
	width           int
	height          int
	validationError bool
}

func NewLoginModel(username, storedPass string) loginModel {
	// Spinner
	sp := spinner.New()
	sp.Spinner = spinner.Points
	// Password Input
	passwordInput := textinput.New()
	passwordInput.CharLimit = 16
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'

	if storedPass == "" {
		passwordInput.Focus()
	}

	model := loginModel{
		username:   username,
		password:   passwordInput,
		spinner:    sp,
		storedPass: storedPass,
	}

	return model
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case statusMsg:
		var cmd tea.Cmd
		switch msg.sType {
		case STATUS_LOGIN:
			m.isLoading = false
			m.isLoggedIn = true
		}

		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case spinner.TickMsg: // Only update the spinner when needed
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	case tea.KeyMsg:
		m.password.Placeholder = "Password"
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if !m.isLoggedIn {
				var password string
				password = m.password.Value()
				if password == "" {
					m.password.Placeholder = "Password cannot be empty"
					m.validationError = true
				} else if len(password) < 7 {
					m.password.Placeholder = "Password must be at least 7 characters"
					m.validationError = true
				} else {
					m.validationError = false
				}
				m.password, cmd = m.password.Update(msg)
				if m.validationError {
					m.password.Reset()
					return m, cmd
				}
				m.isLoading = true
				m.password.Blur() // Blur current active input
				return m, tea.Batch(m.spinner.Tick, cmd, login(map[string]string{"password": password, "username": m.username}))
			} else if m.isLoggedIn {
				log.Println(m.height, m.width)
				return enterChat(m.width, m.height)
			}
		default:
			if m.isLoggedIn {
				return enterChat(m.width, m.height)
			}
		}
	}
	if m.storedPass != "" {
		m.isLoading = true
		return m, tea.Batch(m.spinner.Tick, login(map[string]string{"password": m.storedPass, "username": m.username}))
	}
	m.password, cmd = m.password.Update(msg)
	return m, cmd
}

func (m loginModel) View() string {
	var b strings.Builder

	if m.storedPass == "" {
		b.WriteString(m.password.View())
	}

	if !m.validationError {
		if m.isLoading {
			if m.storedPass == "" {
				b.WriteRune('\n')
				b.WriteRune('\n')
			}
			b.WriteString(fmt.Sprintf("Logging into your profile %s\n", m.spinner.View()))
		} else if !m.isLoading && m.isLoggedIn {
			if m.storedPass == "" {
				b.WriteRune('\n')
				b.WriteRune('\n')
			}
			b.WriteString(fmt.Sprintf("Logged in %c\n\n", ICON_DONE))
			b.WriteString("press any key to proceed\n")
		}
	}
	return b.String()
}

func (m loginModel) Init() tea.Cmd {
	return nil
}
