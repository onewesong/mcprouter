package sse

import (
	"github.com/chatmcp/mcprouter/service/mcpclient"
)

// SSESession is a session for SSE request
type SSESession struct {
	writer   *SSEWriter
	messages chan string // event queue
	command  string
	client   *mcpclient.StdioClient
}

// NewSSESession will create a new SSE session
func NewSSESession(w *SSEWriter, command string) *SSESession {
	return &SSESession{
		writer:   w,
		messages: make(chan string, 100),
		command:  command,
		client:   nil,
	}
}

func (s *SSESession) Command() string {
	return s.command
}

func (s *SSESession) SetClient(client *mcpclient.StdioClient) {
	s.client = client
}

func (s *SSESession) Client() *mcpclient.StdioClient {
	return s.client
}

func (s *SSESession) Messages() chan string {
	return s.messages
}

func (s *SSESession) SendMessage(message string) {
	s.messages <- message
}

func (s *SSESession) Close() {
	close(s.messages)
}
