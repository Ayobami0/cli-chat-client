package ui

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Ayobami0/cli-chat/pb"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	client             pb.ChatServiceClient
	sessionToken       string
	user               *pb.User
	chatStream         pb.ChatService_ChatStreamClient
	msgChan            chan *pb.MessageStream
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
		return m, tea.Batch(m.chatList.StartSpinner(), m.getChats())
	}
	if !m.requestsLoading && !m.requestsLoaded {
		m.requestsLoading = true
		return m, tea.Batch(m.requestsList.StartSpinner(), m.getRequests())
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
		case STATUS_REQUEST_ACTION_SEND:
			m.chatsLoading = false
			m.chatsLoaded = false
		case STATUS_GROUP_REQUEST_SEND:
		case STATUS_CHATS_LOAD:
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
		case STATUS_MESSAGE_RECV:
			message := msg.sRes.(*pb.MessageStream)

			if message == nil {
				return m, tea.Batch(vCmd, iCmd, wait(m.msgChan))
			}
			m.messages = []string{}

			for _, v := range m.chatList.SelectedItem().(chatItem).messages {
				m.messages = append(m.messages, m.formatMessage(v))
			}
			m.messages = append(m.messages, m.formatMessage(message.Message))
			vStr := strings.Join(m.messages, "\n")
			m.viewport.SetContent(vStr)
			m.viewport.GotoBottom()

			m.viewport, vCmd = m.viewport.Update(msg)
			m.input, iCmd = m.input.Update(msg)

			return m, tea.Batch(vCmd, iCmd, wait(m.msgChan), m.getChats())
		}
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			m.chatStream.CloseSend()
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
						req := m.requestsList.SelectedItem().(requestItem)
						m.requestsList.RemoveItem(m.requestsList.Index())
						m.requestsList, rCmd = m.requestsList.Update(msg)
						return m, tea.Batch(rCmd, m.sendRequestAction(req.id, pb.DirectChatAction_ACTION_ACCEPT))
					}
				case "ctrl+x":
					if m.focusedPanel == ACTIVE_REQUEST_PANEL {
						req := m.requestsList.SelectedItem().(requestItem)
						m.requestsList.RemoveItem(m.requestsList.Index())
						m.requestsList, rCmd = m.requestsList.Update(msg)
						return m, tea.Batch(rCmd, m.sendRequestAction(req.id, pb.DirectChatAction_ACTION_REJECT))
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
							name, passkey := m.nameChatInput.Value(), m.passkeyChatInput.Value()
							m.chatList.Select(len(m.chatList.Items()) - 1)
							m.passkeyChatInput.Reset()
							m.nameChatInput.Reset()
							m.joinRoomFocusIndex = 0
							m.joinGroupLoading = true
							m.nameChatInput.Blur()
							m.passkeyChatInput.Blur()
							m.nameChatInput, joinNameCmd = m.nameChatInput.Update(msg)
							m.passkeyChatInput, joinPassCmd = m.passkeyChatInput.Update(msg)
							return m, tea.Batch(joinNameCmd, joinPassCmd, m.progressIndicator.Tick, m.sendGroupChatJoinRequest(name, passkey))
						}
					case SEND_REQUEST_PANNEL:
						receiver := m.addUserInput.Value()
						m.addUserInput.Reset()
						m.addUserInput, sndReqCmd = m.addUserInput.Update(msg)
						m.sendRequestLoading = true
						return m, tea.Batch(sndReqCmd, m.progressIndicator.Tick, m.sendDirectChatJoinRequest(receiver))
					case MESSAGE_PANEL:
						if m.chatStream != nil {

							chat := m.chatList.SelectedItem().(chatItem)

							m.viewport.SetContent(strings.Join(m.messages, "\n"))
							msgContent := m.input.Value()
							m.input.Reset()
							m.input, iCmd = m.input.Update(msg)

							return m, tea.Batch(iCmd,
								m.send(&pb.MessageStream{ChatId: chat.id, Message: &pb.Message{
									Sender:  m.user,
									Content: msgContent,
									SentAt:  timestamppb.Now(),
									Type:    pb.Message_MESSAGE_TYPE_REGULAR,
								}}))
						}
					case CHATS_PANEL:
						if m.chatStream != nil {
							m.chatStream.CloseSend()
						}

						m.msgChan = make(chan *pb.MessageStream)

						chat := m.chatList.SelectedItem().(chatItem)

						m.viewport.GotoBottom()
						m.viewport, vCmd = m.viewport.Update(msg)
						m.input, iCmd = m.input.Update(msg)
						m.focusedPanel = MESSAGE_PANEL

						err := m.initializeStream(chat.id)

						if err != nil {
							return m, tea.Quit
						}

						return m, tea.Batch(vCmd, iCmd, lCmd, m.recv(), wait(m.msgChan))
					}
				}
			}
		}
	case errMsg:
		if msg.err == io.EOF {
			return m, m.getChats()
		}
		return m, tea.Quit
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
	log.Println(m)
	return nil
}

func NewChatModel(client pb.ChatServiceClient, w, h int, auth *pb.UserAuthenticatedResponse) chatModel {
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
	// disable movements
	vp.KeyMap.Down.SetEnabled(false)
	vp.KeyMap.Up.SetEnabled(false)
	vp.KeyMap.PageUp.SetEnabled(false)
	vp.KeyMap.PageDown.SetEnabled(false)
	vp.KeyMap.HalfPageUp.SetEnabled(false)
	vp.KeyMap.HalfPageDown.SetEnabled(false)

	m := chatModel{
		user:              auth.User,
		input:             ta,
		viewport:          vp,
		messages:          []string{},
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
		sessionToken:      auth.Token,
		client:            client,
	}

	return m
}

func enterChat(client pb.ChatServiceClient, w, h int, auth *pb.UserAuthenticatedResponse) (chatModel, tea.Cmd) {
	altScrCmd := tea.EnterAltScreen
	return NewChatModel(client, w, h, auth), altScrCmd
}

func (m chatModel) formatMessage(msg *pb.Message) string {
	var fmtMsg string
	switch msg.Type {
	case pb.Message_MESSAGE_TYPE_REGULAR:
		snder := "Me"
		if msg.Sender.Username != m.user.Username {
			snder = msg.Sender.Username
		}
		fmtMsg = fmt.Sprintf("%s: %s", snder, msg.Content)
	case pb.Message_MESSAGE_TYPE_NOTIFICATION:
		fmtMsg = lipgloss.PlaceHorizontal(
			lipgloss.Width(
				m.viewport.View(),
			),
			lipgloss.Center,
			notificationTextStyle.Render(
				msg.Content,
			),
		)

	}

	return fmtMsg
}

func (c chatModel) send(msgStream *pb.MessageStream) tea.Cmd {
	return func() tea.Msg {
		err := c.chatStream.Send(msgStream)

		if err != nil {
			return errMsg{err}
		}
		return statusMsg{sType: STATUS_MESSAGE_SEND}
	}
}

func (c *chatModel) recv() tea.Cmd {
	return func() tea.Msg {
		for {
			msg, err := c.chatStream.Recv()
			if err == io.EOF {
				close(c.msgChan)
				return errMsg{err}
			}
			if err != nil {
				return errMsg{err}
			}
			c.msgChan <- msg
		}

	}
}

func wait(chat chan *pb.MessageStream) tea.Cmd {
	return func() tea.Msg {
		return statusMsg{sRes: <-chat, sType: STATUS_MESSAGE_RECV}
	}
}

func (c *chatModel) initializeStream(chatID string) error {
	meta := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", c.sessionToken), "stream_chat_id", chatID, "stream_username", c.user.Username)
	ctx := metadata.NewOutgoingContext(context.Background(), meta)

	stream, err := c.client.ChatStream(ctx)

	if err != nil {
		return err
	}

	c.chatStream = stream

	return nil
}

func (c chatModel) sendGroupChatJoinRequest(groupName, groupPasskey string) tea.Cmd {
	return func() tea.Msg {
		meta := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", c.sessionToken))

		ctx := metadata.NewOutgoingContext(context.Background(), meta)

		res, err := c.client.JoinGroupChat(ctx, &pb.GroupChatRequest{GroupName: groupName, GroupPasskey: groupPasskey})
		if err != nil {
			return errMsg{err}
		}

		return statusMsg{sType: STATUS_GROUP_REQUEST_SEND, sRes: res}
	}
}

func (c chatModel) getChats() tea.Cmd {
	return func() tea.Msg {
		meta := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", c.sessionToken))

		log.Println(meta)
		ctx := metadata.NewOutgoingContext(context.Background(), meta)

		res, err := c.client.GetChats(ctx, &emptypb.Empty{})
		if err != nil {
			return errMsg{err}
		}

		var chatItems []chatItem
		for _, v := range res.Chats {
			chatItems = append(chatItems, chatItem{
				id:       v.Id,
				name:     *v.Name,
				messages: v.Messages,
				chatType: v.Type,
				members:  v.Members,
			})
		}
		return statusMsg{sType: STATUS_CHATS_LOAD, sRes: chatItems}
	}
}

func (c chatModel) sendDirectChatJoinRequest(receiver string) tea.Cmd {
	return func() tea.Msg {
		meta := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", c.sessionToken))

		log.Println(meta)
		ctx := metadata.NewOutgoingContext(context.Background(), meta)

		res, err := c.client.JoinDirectChat(ctx, &pb.JoinDirectChatRequest{SentAt: timestamppb.Now(), Receiver: &pb.User{Username: receiver}})
		if err != nil {
			return errMsg{err}
		}

		return statusMsg{sType: STATUS_DIRECT_REQUEST_SEND, sRes: res}
	}
}

func (c chatModel) getRequests() tea.Cmd {
	return func() tea.Msg {
		meta := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", c.sessionToken))

		ctx := metadata.NewOutgoingContext(context.Background(), meta)

		res, err := c.client.GetDirectChatRequests(ctx, &emptypb.Empty{})
		if err != nil {
			return errMsg{err}
		}

		var requestItems []requestItem
		for _, v := range res.Requests {
			requestItems = append(requestItems, requestItem{
				name: *&v.Sender.Username,
				id:   v.Id,
			})
		}
		return statusMsg{sType: STATUS_REQUEST_LOAD, sRes: requestItems}
	}
}

func (c chatModel) sendRequestAction(chatRequestId string, action pb.DirectChatAction_Action) tea.Cmd {
	return func() tea.Msg {
		meta := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", c.sessionToken))

		log.Println(chatRequestId)
		ctx := metadata.NewOutgoingContext(context.Background(), meta)

		res, err := c.client.DirectChatRequestAction(ctx, &pb.DirectChatAction{Action: action, Id: chatRequestId})
		if err != nil {
			return errMsg{err}
		}

		return statusMsg{sType: STATUS_REQUEST_ACTION_SEND, sRes: res}
	}
}
