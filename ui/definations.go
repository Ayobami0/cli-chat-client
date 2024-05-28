package ui

import (
	"fmt"
	"time"

	"github.com/Ayobami0/cli-chat/pb"
)

type statusType int

type statusMsg struct {
	sType statusType
	sRes  interface{}
	sCode int
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type chatItem struct {
	chatType  pb.ChatType
	name      string
	id        string
	maxMember int
	members   []*pb.User
	messages  []*pb.Message
}

func (c chatItem) Title() string {
	if c.chatType == pb.ChatType_CHAT_TYPE_DIRECT {
		return fmt.Sprintf("%s + %s", c.members[0].Username, c.members[1].Username)
	}
	return c.name
}
func (c chatItem) Description() string {
	if len(c.messages) == 0 {
		return ""
	}
	return c.messages[len(c.messages)-1].Content
}
func (c chatItem) FilterValue() string { return c.name }

type requestItem struct {
	name   string
	id     string
	sentAt time.Time
}

func (r requestItem) Title() string       { return r.name }
func (r requestItem) Description() string { return r.sentAt.String() }
func (r requestItem) FilterValue() string { return r.name }
