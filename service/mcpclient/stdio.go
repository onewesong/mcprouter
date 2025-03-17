package mcpclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
	"github.com/spf13/cast"
)

// StdioClient is a client that uses stdin and stdout to communicate with the backend mcp server.
type StdioClient struct {
	cmd      *exec.Cmd
	stdin    io.WriteCloser
	stdout   *bufio.Reader
	stderr   *bufio.Reader
	done     chan struct{}
	messages map[int64]chan *jsonrpc.Response // store stdout messages
	err      chan error                       // store stderr messages
	mu       sync.RWMutex
}

// NewStdioClient creates a new StdioClient.
func NewStdioClient(command string) (*StdioClient, error) {
	cmd := exec.Command(
		"sh",
		"-c",
		command,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	client := &StdioClient{
		cmd:      cmd,
		stdin:    stdin,
		stdout:   bufio.NewReader(stdout),
		stderr:   bufio.NewReader(stderr),
		done:     make(chan struct{}),
		messages: make(map[int64]chan *jsonrpc.Response),
		err:      make(chan error, 1),
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			errmsg := scanner.Text()
			fmt.Printf("stderr: %s\n", errmsg)
			// client.err <- fmt.Errorf("mcp server run failed: %s", errmsg)
			// client.Close()
			// return
		}

		if err := scanner.Err(); err != nil {
			client.err <- fmt.Errorf("stderr scanner error: %w", err)
			client.Close()
		}
	}()

	fmt.Printf("mcp server running with command: %s\n", command)

	ready := make(chan struct{})
	go func() {
		close(ready)
		client.listen()
	}()
	<-ready

	return client, nil
}

// listen for messages from the backend mcp server.
func (c *StdioClient) listen() {
	for {
		select {
		case <-c.done:
			return

		default:
			message, err := c.stdout.ReadBytes('\n')
			fmt.Printf("stdout read message: %s\n", string(message))
			if err != nil {
				if err != io.EOF {
					fmt.Printf("failed to read message: %v\n", err)
				}
				return
			}

			var response *jsonrpc.Response
			if err := json.Unmarshal(message, &response); err != nil {
				fmt.Printf("failed to parse response: %v\n", err)
				continue
			}

			if response.ID == nil {
				// handle notification message
				continue
			}

			id := cast.ToInt64(response.ID)

			// result or error message
			c.mu.RLock()
			messages, ok := c.messages[id]
			c.mu.RUnlock()

			if !ok {
				fmt.Printf("invalid message with id: %d\n", id)
				continue
			}

			messages <- response
		}
	}
}

// ForwardRequest forwards a JSON-RPC request to the MCP server and returns the response
func (c *StdioClient) ForwardRequest(request *jsonrpc.Request) *jsonrpc.Response {
	// fmt.Printf("forward request: %+v\n", request)

	response, err := c.SendRequest(request)
	if err != nil {
		fmt.Printf("failed to forward request: %v\n", err)
		return jsonrpc.NewErrorResponse(jsonrpc.ErrorProxyError, request.ID)
	}

	// fmt.Printf("forward response: %+v\n", response)

	return response
}

func (c *StdioClient) SendRequest(request *jsonrpc.Request) (*jsonrpc.Response, error) {
	message, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request message: %w", err)
	}
	message = append(message, '\n')

	if request.ID == nil {
		// notification message
		if _, err := c.stdin.Write(message); err != nil {
			return nil, fmt.Errorf("failed to write request message: %w", err)
		}

		fmt.Printf("stdin write message: %s\n", string(message))

		return nil, nil
	}

	id := cast.ToInt64(request.ID)
	messages := make(chan *jsonrpc.Response, 1)

	c.mu.Lock()
	c.messages[id] = messages
	c.mu.Unlock()

	if _, err := c.stdin.Write(message); err != nil {
		return nil, fmt.Errorf("failed to write request message: %w", err)
	}

	fmt.Printf("stdin write message: %s\n", string(message))

	for {
		select {
		case <-c.done:
			return nil, fmt.Errorf("client closed")
		case err := <-c.err:
			return nil, err
		case response := <-messages:
			return response, nil
		}
	}
}

// Error returns the error message from stderr
func (c *StdioClient) Error() error {
	select {
	case err := <-c.err:
		return err
	default:
		return nil
	}
}

// Close client
func (c *StdioClient) Close() error {
	close(c.done)
	if err := c.stdin.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	return c.cmd.Wait()
}
