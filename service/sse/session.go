package sse

// SSESession is a session for SSE request
type SSESession struct {
	writer   *SSEWriter
	messages chan string // event queue
	key      string
}

// NewSSESession will create a new SSE session
func NewSSESession(w *SSEWriter, key string) *SSESession {
	return &SSESession{
		writer:   w,
		messages: make(chan string, 100),
		key:      key,
	}
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
