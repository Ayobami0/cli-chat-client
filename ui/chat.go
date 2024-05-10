package ui

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	CHATS_PANEL int = iota
	ACTIVE_REQUEST_PANEL
	SEND_REQUEST_PANNEL
	JOIN_ROOM_PANEL
	MESSAGE_PANEL
	MESSSAGE_VIEW_PANEL
	MAX_PANEL_NO
)

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Help        key.Binding
	Quit        key.Binding
	Enter       key.Binding
	Accept      key.Binding
	Reject      key.Binding
	SwitchPanel key.Binding
	Tab         key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.SwitchPanel, k.Tab}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Reject, k.Accept},
		{k.Help, k.Quit},
		{k.Enter, k.SwitchPanel},
	}
}

type chatModel struct {
	chatsLoading       bool
	requestsLoading    bool
	chatsLoaded        bool
	requestsLoaded     bool
	viewport           viewport.Model
	input              textarea.Model
	help               help.Model
	chatList           list.Model
	requestsList       list.Model
	messages           []string
	err                error
	width              int
	height             int
	chats              []string
	focusedPanel       int
	keys               keyMap
	helpHeight         int
	sndRequestHeight   int
	joinRoomHeight     int
	loadingMsg         bool
	addUserInput       textinput.Model
	nameChatInput      textinput.Model
	passkeyChatInput   textinput.Model
	progressIndicator  spinner.Model
	joinRoomFocusIndex int
	sendRequestLoading bool
	joinGroupLoading   bool
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		vCmd        tea.Cmd
		iCmd        tea.Cmd
		lCmd        tea.Cmd
		hCmd        tea.Cmd
		sCmd        tea.Cmd
		rCmd        tea.Cmd
		sndReqCmd   tea.Cmd
		joinNameCmd tea.Cmd
		joinPassCmd tea.Cmd
	)

	if !m.chatsLoading && !m.chatsLoaded {
		m.chatsLoading = true
		return m, tea.Batch(m.chatList.StartSpinner(), getChats())
	}
	if !m.requestsLoading && !m.requestsLoaded {
		m.requestsLoading = true
		return m, tea.Batch(m.requestsList.StartSpinner(), getRequests())
	}

	switch m.focusedPanel {
	case MESSAGE_PANEL:
		m.input.Focus()
		m.nameChatInput.Blur()
		m.passkeyChatInput.Blur()
	case JOIN_ROOM_PANEL:
		if !m.joinGroupLoading {
			if m.joinRoomFocusIndex == 0 {
				m.nameChatInput.Focus()
				m.passkeyChatInput.Blur()
			} else {
				m.nameChatInput.Blur()
				m.passkeyChatInput.Focus()
			}
		} else {
			m.nameChatInput.Blur()
			m.passkeyChatInput.Blur()
		}
		m.addUserInput.Blur()
		m.input.Blur()
	case SEND_REQUEST_PANNEL:
		if !m.sendRequestLoading {
			m.addUserInput.Focus()
		} else {
			m.addUserInput.Blur()
		}
		m.input.Blur()
		m.passkeyChatInput.Blur()
		m.nameChatInput.Blur()
	default:
		m.input.Blur()
		m.addUserInput.Blur()
		m.passkeyChatInput.Blur()
		m.nameChatInput.Blur()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		m.help.Width = m.width
		m.input.SetWidth(m.width - m.chatList.Width() - 4) // 4 is to account for the borders
		m.viewport.Width = m.width - lipgloss.Width(m.chatList.View()) - 4

		m.chatList.SetHeight(m.height - m.helpHeight - m.joinRoomHeight - m.sndRequestHeight - m.requestsList.Height() - 6)
		m.viewport.Height = m.height - lipgloss.Height(m.input.View()) - m.helpHeight - 6
	case spinner.TickMsg:
		m.progressIndicator, sCmd = m.progressIndicator.Update(msg)
		m.chatList, lCmd = m.chatList.Update(msg)
		m.requestsList, rCmd = m.requestsList.Update(msg)

		return m, tea.Batch(sCmd, lCmd, rCmd)
	case statusMsg:
		switch msg.sType {
		case STATUS_MESSAGE_LOAD:
			m.messages = []string{}
			m.loadingMsg = false
			sItem, _ := m.chatList.SelectedItem().(chatItem)
			m.viewport.SetContent(
				lipgloss.PlaceHorizontal(
					lipgloss.Width(
						m.viewport.View(),
					),
					lipgloss.Center,
					notificationTextStyle.Render(
						fmt.Sprintf(" You are now chatting in %s ", sItem.name),
					),
				),
			)
		case STATUS_CHAT_LOAD:
			var chats []list.Item

			for _, v := range msg.sRes.([]chatItem) {
				chats = append(chats, v)
			}
			m.chatsLoading = false
			m.chatsLoaded = true

			m.chatList.StopSpinner()
			m.chatList.SetItems(chats)
		case STATUS_REQUEST_LOAD:
			var requests []list.Item

			for _, v := range msg.sRes.([]requestItem) {
				requests = append(requests, v)
			}

			m.requestsLoading = false
			m.requestsLoaded = true

			m.requestsList.StopSpinner()
			m.requestsList.SetItems(requests)
		}
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			return m, tea.Quit
		}
		if msg.String() == "?" {
			if m.focusedPanel != MESSAGE_PANEL {
				m.help.ShowAll = !m.help.ShowAll
			}
		}
		if !m.joinGroupLoading && !m.sendRequestLoading {
			switch msg.String() {
			case "tab":
				if m.focusedPanel == MAX_PANEL_NO-1 {
					m.focusedPanel = 0
				} else {
					m.focusedPanel++
				}
			case "shift+tab":
				if m.focusedPanel == 0 {
					m.focusedPanel = MAX_PANEL_NO - 1
				} else {
					m.focusedPanel--
				}
			case "alt+1", "alt+2", "alt+3", "alt+4", "alt+5", "alt+6":
				if m.input.Focused() {
					m.input.Blur()
				}
				var focused int
				switch msg.String() {
				case "alt+1":
					focused = CHATS_PANEL
				case "alt+2":
					focused = ACTIVE_REQUEST_PANEL
				case "alt+3":
					focused = SEND_REQUEST_PANNEL
				case "alt+4":
					focused = JOIN_ROOM_PANEL
				case "alt+5":
					focused = MESSAGE_PANEL
				case "alt+6":
					focused = MESSSAGE_VIEW_PANEL
				}
				m.focusedPanel = focused
			default:
				switch msg.String() {
				case "ctrl+a":
					if m.focusedPanel == ACTIVE_REQUEST_PANEL {
						//TODO: Add to friend list
						m.requestsList.RemoveItem(m.requestsList.Index())
					}
				case "ctrl+x":
					if m.focusedPanel == ACTIVE_REQUEST_PANEL {
						//TODO: Cancel request
						m.requestsList.RemoveItem(m.requestsList.Index())
					}
				case "up":
					switch m.focusedPanel {
					case JOIN_ROOM_PANEL:
						if m.joinRoomFocusIndex == 0 {
							m.passkeyChatInput.Focus()
							m.joinRoomFocusIndex++
						} else {
							m.nameChatInput.Focus()
							m.joinRoomFocusIndex--
						}
						m.nameChatInput, joinNameCmd = m.nameChatInput.Update(msg)
						m.passkeyChatInput, joinPassCmd = m.passkeyChatInput.Update(msg)
						return m, tea.Batch(joinNameCmd, joinPassCmd)
					case MESSSAGE_VIEW_PANEL:
						m.viewport.HalfViewUp()
					case CHATS_PANEL:
						m.chatList.CursorUp()
					case ACTIVE_REQUEST_PANEL:
						m.requestsList.CursorUp()
					}
				case "down":
					switch m.focusedPanel {
					case JOIN_ROOM_PANEL:
						if m.joinRoomFocusIndex == 0 {
							m.passkeyChatInput.Focus()
							m.joinRoomFocusIndex++
						} else {
							m.nameChatInput.Focus()
							m.joinRoomFocusIndex--
						}
						m.nameChatInput, joinNameCmd = m.nameChatInput.Update(msg)
						m.passkeyChatInput, joinPassCmd = m.passkeyChatInput.Update(msg)
						return m, tea.Batch(joinNameCmd, joinPassCmd)
					case MESSSAGE_VIEW_PANEL:
						m.viewport.HalfViewDown()
					case CHATS_PANEL:
						m.chatList.CursorDown()
					case ACTIVE_REQUEST_PANEL:
						m.requestsList.CursorDown()
					}
				case "enter":
					switch m.focusedPanel {
					case JOIN_ROOM_PANEL:
						if m.joinRoomFocusIndex == 0 {
							m.joinRoomFocusIndex = 1
						} else {
							// TODO: Add joining room function and inserting to chat
							m.chatList.Select(len(m.chatList.Items()) - 1)
							m.passkeyChatInput.Reset()
							m.nameChatInput.Reset()
							m.joinRoomFocusIndex = 0
							m.joinGroupLoading = true
							m.nameChatInput.Blur()
							m.passkeyChatInput.Blur()
							m.nameChatInput, joinNameCmd = m.nameChatInput.Update(msg)
							m.passkeyChatInput, joinPassCmd = m.passkeyChatInput.Update(msg)
							return m, tea.Batch(joinNameCmd, joinPassCmd, m.progressIndicator.Tick)
						}
					case SEND_REQUEST_PANNEL:
						// TODO: Add sending request funtionality
						m.addUserInput.Reset()
						m.addUserInput, sndReqCmd = m.addUserInput.Update(msg)
						m.sendRequestLoading = true
						return m, tea.Batch(sndReqCmd, m.progressIndicator.Tick)
					case MESSAGE_PANEL:
						m.messages = append(m.messages, "Me: "+m.input.Value())
						m.viewport.SetContent(strings.Join(m.messages, "\n"))
						m.input.Reset()
						m.viewport.GotoBottom()
					case CHATS_PANEL:
						m.viewport.SetContent("Loading...")
						m.viewport, vCmd = m.viewport.Update(msg)
						m.focusedPanel = MESSAGE_PANEL
						m.loadingMsg = true
						return m, tea.Batch(getMessages("test", "test"), vCmd)
					}
				}
			}
		}
	}
	m.viewport, vCmd = m.viewport.Update(msg)
	m.input, iCmd = m.input.Update(msg)
	m.chatList, lCmd = m.chatList.Update(msg)
	m.help, hCmd = m.help.Update(msg)
	m.progressIndicator, sCmd = m.progressIndicator.Update(msg)
	m.addUserInput, sndReqCmd = m.addUserInput.Update(msg)
	m.nameChatInput, joinNameCmd = m.nameChatInput.Update(msg)
	m.passkeyChatInput, joinPassCmd = m.passkeyChatInput.Update(msg)
	m.requestsList, rCmd = m.requestsList.Update(msg)

	return m, tea.Batch(iCmd, vCmd, lCmd, hCmd, joinPassCmd, joinNameCmd, sndReqCmd, sCmd, rCmd)
}

func (m chatModel) View() string {
	chatView := unfocusedBorderStyle
	inputView := unfocusedBorderStyle
	listView := unfocusedBorderStyle
	requestView := unfocusedBorderStyle
	sendRequestView := unfocusedBorderStyle
	joinRoomView := unfocusedBorderStyle

	var joinPlaceholder string
	var sendRequestPlaceholder string

	if m.sendRequestLoading {
		sendRequestPlaceholder = m.progressIndicator.View()
	} else {
		sendRequestPlaceholder = ""
	}
	if m.joinGroupLoading {
		joinPlaceholder = m.progressIndicator.View()
	} else {
		joinPlaceholder = ""
	}

	switch m.focusedPanel {
	case CHATS_PANEL:
		listView = focusedBorderStyle
	case MESSAGE_PANEL:
		inputView = focusedBorderStyle
	case MESSSAGE_VIEW_PANEL:
		chatView = focusedBorderStyle
	case ACTIVE_REQUEST_PANEL:
		requestView = focusedBorderStyle
	case SEND_REQUEST_PANNEL:
		sendRequestView = focusedBorderStyle
	case JOIN_ROOM_PANEL:
		joinRoomView = focusedBorderStyle
	}

	return lipgloss.JoinVertical(
		lipgloss.Top, lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Top,
				listView.Render(m.chatList.View()),
				requestView.Render(m.requestsList.View()),
				sendRequestView.Render(
					fmt.Sprintf(
						"Send a chat request\n%s\n%s",
						m.addUserInput.View(),
						sendRequestPlaceholder,
					),
				),
				joinRoomView.Render(
					fmt.Sprintf(
						"Join a room\n%s\n%s\n%s",
						m.nameChatInput.View(),
						m.passkeyChatInput.View(),
						joinPlaceholder,
					),
				),
			),
			lipgloss.JoinVertical(
				lipgloss.Top,
				chatView.Render(m.viewport.View()),
				inputView.Render(m.input.View()),
			),
		),
		m.help.View(m.keys),
	)
}

func (m chatModel) Init() tea.Cmd {
	return nil
}

func NewChatModel(w, h int) chatModel {
	var keys = keyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up  "),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down  "),
		),
		Tab: key.NewBinding(
			key.WithKeys("󰌥/󰌒"),
			key.WithHelp("󰌥/󰌒", "cycle panels  "),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help  "),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit  "),
		),
		Enter: key.NewBinding(
			key.WithKeys("󰌑"),
			key.WithHelp("󰌑", "enter room/send message  "),
		),
		Accept: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "accept request  "),
		),
		Reject: key.NewBinding(
			key.WithKeys("ctrl+x"),
			key.WithHelp("ctrl+x", "reject request  "),
		),
		SwitchPanel: key.NewBinding(
			key.WithKeys("alt+[n]"),
			key.WithHelp("alt+[n]", "switch panel (1|chats 2|requests 3|send request 4|join room 5|chat input 6|chat view)  "),
		),
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	joinNameInput := textinput.New()
	joinNameInput.Placeholder = "Chatroom Name"
	joinNameInput.Blur()
	joinNameInput.Width = 28

	joinPasskeyInput := textinput.New()
	joinPasskeyInput.Placeholder = "Chatroom Passkey"
	joinNameInput.Blur()
	joinPasskeyInput.Width = 28

	sendRequestInput := textinput.New()
	sendRequestInput.Placeholder = "Username"
	sendRequestInput.Blur()
	sendRequestInput.Width = 28

	sndReqPanelHeight := 5
	joinRoomPanelHeight := 6

	hp := help.New()
	hp.Styles = helpStyle
	hpHeight := strings.Count(hp.View(keys), "\n")

	reqLt := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	reqLt.InfiniteScrolling = true
	reqLt.Title = "Requests"
	reqLt.SetShowPagination(false)
	reqLt.SetShowStatusBar(false)
	reqLt.SetShowHelp(false)
	reqLt.SetFilteringEnabled(false)
	reqLt.SetSize(31, 8)
	reqLt.KeyMap = list.KeyMap{}
	reqLt.SetSpinner(spinner.Dot)

	lt := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	lt.InfiniteScrolling = true
	lt.Title = "Chat List"
	lt.SetShowPagination(false)
	lt.SetShowStatusBar(false)
	lt.SetShowHelp(false)
	lt.SetFilteringEnabled(false)
	lt.SetSize(31, h-hpHeight-joinRoomPanelHeight-sndReqPanelHeight-reqLt.Height()-6)
	lt.KeyMap = list.KeyMap{}
	lt.SetSpinner(spinner.Dot)

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Reset()
	ta.Prompt = "┃ "
	ta.SetHeight(1)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.Blur()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.SetWidth(w - lt.Width() - 4)

	vp := viewport.New(w-lt.Width()-4, h-7-hpHeight)

	return chatModel{
		input:             ta,
		viewport:          vp,
		messages:          []string{},
		err:               nil,
		width:             w,
		height:            h,
		chatList:          lt,
		keys:              keys,
		requestsList:      reqLt,
		helpHeight:        hpHeight,
		loadingMsg:        true,
		nameChatInput:     joinNameInput,
		passkeyChatInput:  joinPasskeyInput,
		addUserInput:      sendRequestInput,
		sndRequestHeight:  sndReqPanelHeight,
		joinRoomHeight:    joinRoomPanelHeight,
		progressIndicator: sp,
	}
}

func enterChat(w, h int) (chatModel, tea.Cmd) {
	altScrCmd := tea.EnterAltScreen
	return NewChatModel(w, h), altScrCmd
}

func getChats() tea.Cmd {
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
		return statusMsg{sType: STATUS_CHAT_LOAD, sRes: []chatItem{
			{name: "User I", lastMessage: "Hello"},
			{name: "Group I", lastMessage: "Jide: A long ass text to see the width of the list"},
		}}
	}
}

func getRequests() tea.Cmd {
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
		return statusMsg{sType: STATUS_REQUEST_LOAD, sRes: []requestItem{
			{name: "Ayobami", sentAt: time.Now()},
		}}
	}
}

func getMessages(chat, username string) tea.Cmd {
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
		log.Println(chat, username)
		return statusMsg{sType: STATUS_MESSAGE_LOAD}
	}
}
