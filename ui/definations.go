package ui

import "time"

type statusMsg struct {
	sType string
	sRes  interface{}
	sCode int
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type chatItem struct {
	chatType    string
	name        string
	id          int
	maxMember   int
	lastMessage string
}

func (c chatItem) Title() string       { return c.name }
func (c chatItem) Description() string { return c.lastMessage }
func (c chatItem) FilterValue() string { return c.name }

type requestItem struct {
	name   string
	id     int
	sentAt time.Time
}

func (r requestItem) Title() string       { return r.name }
func (r requestItem) Description() string { return r.sentAt.String() }
func (r requestItem) FilterValue() string { return r.name }
