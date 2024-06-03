package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ayobami0/cli-chat/pb"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginModel struct {
	username        string
	password        textinput.Model
	isLoading       bool
	isLoggedIn      bool
	spinner         spinner.Model
	width           int
	height          int
	validationError bool
	client          pb.ChatServiceClient
	authRes         *pb.UserAuthenticatedResponse
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case errMsg:
		fmt.Println("ERROR: " + msg.Error())
		return m, tea.Quit
	case statusMsg:
		var cmd tea.Cmd
		switch msg.sType {
		case STATUS_LOGIN:
			m.isLoading = false
			m.isLoggedIn = true

		}

		m.spinner, cmd = m.spinner.Update(msg)
		loginRes := msg.sRes.(*pb.UserAuthenticatedResponse)

		m.authRes = loginRes
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
				return m, tea.Batch(m.spinner.Tick, cmd, login(
					map[string]string{"password": password, "username": m.username},
					m.client,
				),
				)
			} else if m.isLoggedIn {
				return enterChat(m.client, m.width, m.height, m.authRes)
			}
		default:
			if m.isLoggedIn {
				return enterChat(m.client, m.width, m.height, m.authRes)
			}
		}
	}
	m.password, cmd = m.password.Update(msg)
	return m, cmd
}

func (m loginModel) View() string {
	var b strings.Builder

	b.WriteString(m.password.View())

	if !m.validationError {
		if m.isLoading {
			b.WriteString(fmt.Sprintf("\n\n%s logging into your profile\n", m.spinner.View()))
		} else if !m.isLoading && m.isLoggedIn {
			b.WriteString((fmt.Sprintf("\n\n%c logged in\n", ICON_DONE)))
			b.WriteString("Press any key to proceed\n")
		}
	}
	return b.String()
}

func (m loginModel) Init() tea.Cmd {
	return nil
}
func NewLoginModel(username string, client pb.ChatServiceClient) loginModel {
	// Spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	// Password Input
	passwordInput := textinput.New()
	passwordInput.CharLimit = 16
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'

	passwordInput.Focus()

	model := loginModel{
		username: username,
		password: passwordInput,
		spinner:  sp,
		client:   client,
	}

	return model
}

func login(credential map[string]string, client pb.ChatServiceClient) tea.Cmd {
	return func() tea.Msg {
		res, err := client.LogIntoAccount(
			context.Background(),
			&pb.UserRequest{
				Username: credential["username"],
				Password: credential["password"],
			},
		)
		if err != nil {
			return errMsg{err}
		}
		return statusMsg{sType: STATUS_LOGIN, sRes: res}
	}
}
