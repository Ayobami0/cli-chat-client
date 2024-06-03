package ui

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Ayobami0/cli-chat/pb"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type createModel struct {
	inputs          []textinput.Model
	validationError bool
	credentials     map[string]string
	isCreating      bool
	isCreated       bool
	isLoggedIn      bool
	focusedIdx      int
	spinner         spinner.Model
	width           int
	height          int
	client          pb.ChatServiceClient
	authRes         *pb.UserAuthenticatedResponse
}

func (m createModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case errMsg:
		fmt.Println(fmt.Sprintf("ERROR: %s", msg.Error()))
		return m, tea.Quit

	case statusMsg:
		var cmd tea.Cmd
		switch msg.sType {
		case STATUS_CREATE:
			m.isCreating = false
			m.isCreated = true
			m.spinner, cmd = m.spinner.Update(msg)
			return m, tea.Batch(cmd, login(m.credentials, m.client))
		case STATUS_LOGIN:
			m.isLoggedIn = true
			m.spinner, cmd = m.spinner.Update(msg)
			loginRes := msg.sRes.(*pb.UserAuthenticatedResponse)
			m.authRes = loginRes
			return m, tea.Batch(cmd)
		}
	case spinner.TickMsg: // Only update the spinner when needed
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, tea.Batch(cmd)
	case tea.KeyMsg:
		m.setInputsDefaultPlaceholders()
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			if !m.isCreated && !m.isCreating {
				s := msg.String()
				if s == "up" || s == "shift+tab" {
					if m.focusedIdx == 0 {
						m.focusedIdx = len(m.inputs) - 1
					} else {
						m.focusedIdx--
					}
				} else {
					if m.focusedIdx == len(m.inputs)-1 {
						if s == "enter" {
							cred := make(map[string]string, len(m.inputs)-1) // credentials for account creation
							for i := 0; i <= len(m.inputs)-1; i++ {
								var errTxt string
								val := m.inputs[i].Value()
								key := strings.ToLower(m.inputs[i].Placeholder)
								switch key {
								// Validate inputs
								case "password":
									if val == "" {
										errTxt = "Password cannot be empty"
									} else if len(val) < 7 {
										errTxt = "Password must be at least 7 characters"
									}
								case "username":
									if val == "" {
										errTxt = "Username cannot be empty"
									}
								}
								if errTxt != "" {
									m.inputs[i].Placeholder = errTxt
									m.inputs[i].Reset()
									return m, m.updateInputs(msg)
								}
								cred[key] = val
							}
							m.isCreating = true
							m.inputs[m.focusedIdx].Blur() // Blur current active input
							m.setCredentials(cred)
							return m, tea.Batch(m.spinner.Tick, m.updateInputs(msg), create(cred, m.client))
						}
						m.focusedIdx = 0
					} else {
						m.focusedIdx++
					}
				}
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusedIdx {
						cmds[i] = m.inputs[i].Focus()
						continue
					}
					m.inputs[i].Blur()
				}
				return m, tea.Batch(cmds...)
			} else if m.isLoggedIn {
				return enterChat(m.client, m.width, m.height, m.authRes)
			}
		default:
			if m.isLoggedIn && m.isCreated {
				return enterChat(m.client, m.width, m.height, m.authRes)
			}
		}
	}
	cmds := m.updateInputs(msg)

	return m, cmds
}

func (m *createModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m createModel) View() string {
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	switch {
	case m.isCreating:
		b.WriteString(fmt.Sprintf("\n\n%s creating profile\n", m.spinner.View()))
	case m.isCreated && !m.isLoggedIn:
		b.WriteString(fmt.Sprintf("\n\n%c profile created\n", ICON_DONE))
		b.WriteString(fmt.Sprintf("%s logging into your profile\n", m.spinner.View()))
	case m.isLoggedIn && m.isCreated:
		b.WriteString(fmt.Sprintf("\n\n%c profile created\n", ICON_DONE))
		b.WriteString(fmt.Sprintf("%c logged in\n\n", ICON_DONE))
		b.WriteString("Press any key to proceed\n")
	}
	return b.String()
}

func (m createModel) Init() tea.Cmd {
	return nil
}

func (m *createModel) setCredentials(c map[string]string) {
	m.credentials = c
}

// Sets the default parameter for eachinput placeholder
func (m *createModel) setInputsDefaultPlaceholders() {
	for i := range m.inputs {
		switch {
		case strings.Contains(m.inputs[i].Placeholder, "Password"):
			m.inputs[i].Placeholder = "Password"
		case strings.Contains(m.inputs[i].Placeholder, "Username"):
			m.inputs[i].Placeholder = "Username"
		}
	}
}

func NewCreateModel(client pb.ChatServiceClient) createModel {
	// Spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	// Password Input
	passwordInput := textinput.New()
	passwordInput.CharLimit = 16
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'

	// Username
	usernameInput := textinput.New()
	usernameInput.Placeholder = "Username"
	usernameInput.Focus()
	usernameInput.Validate = func(s string) error {
		if s == "" {
			return errors.New("Username cannot be empty")
		}
		return nil
	}

	model := createModel{
		inputs: []textinput.Model{
			usernameInput,
			passwordInput,
		},
		spinner: sp,
		client:  client,
	}

	return model
}

func create(credential map[string]string, client pb.ChatServiceClient) tea.Cmd {
	return func() tea.Msg {
		res, err := client.CreateNewAccount(context.Background(), &pb.UserRequest{
			Username: credential["username"],
			Password: credential["password"],
		})
		if err != nil {
			return errMsg{err}
		}
		return statusMsg{sType: STATUS_CREATE, sRes: res}
	}
}
