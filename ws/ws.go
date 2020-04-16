package ws

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/torlenor/abylebotter/logging"
)

type wsClient interface {
	ReadMessage() (messageType int, p []byte, err error)

	WriteMessage(messageType int, data []byte) error
	WriteJSON(v interface{}) error

	Close() error
}

// Client is an abstraction of a WebSocket client
type Client struct {
	log *logrus.Entry

	messageChan     chan message
	jSONMessageChan chan jSONMessage

	ws             wsClient
	startStopMutex sync.Mutex

	workersWG   sync.WaitGroup
	stopWorkers chan struct{}
}

var dialer = func(wsURL string) (wsClient, error) {
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	return ws, err
}

// NewClient prepares a new ws Client
func NewClient() *Client {
	log := logging.Get("WSClient")
	log.Debug("WS Client is CREATING itself")

	c := Client{
		log: log,

		messageChan:     make(chan message),
		jSONMessageChan: make(chan jSONMessage),
	}
	return &c
}

// Dial is used to connect the ws Client to a WebSocket server
func (c *Client) Dial(wsURL string) error {
	c.startStopMutex.Lock()
	defer c.startStopMutex.Unlock()

	var err error
	c.ws, err = dialer(wsURL)
	if err != nil {
		return err
	}

	c.stopWorkers = make(chan struct{})

	c.workersWG.Add(1)
	go c.wsWriter()

	return nil
}

// Stop sending data to the WebSocket
func (c *Client) Stop() {
	c.startStopMutex.Lock()

	if c.stopWorkers != nil {
		close(c.stopWorkers)
		c.workersWG.Wait()
	}

	c.ws = nil
	c.stopWorkers = nil

	c.startStopMutex.Unlock()
}

// Close the websocket without sending a message to the server.
// It stops the websocket sender first
func (c *Client) Close() error {
	c.Stop()
	return c.ws.Close()
}

// ReadMessage can be used to read the next message from WebSocket.
// it blocks until somehing is received or the ws is closed.
func (c *Client) ReadMessage() (int, []byte, error) {
	return c.ws.ReadMessage()
}

// SendMessage is used to send a message via the connected WebSocket to the server
func (c *Client) SendMessage(messageType int, data []byte) error {
	c.startStopMutex.Lock()
	defer c.startStopMutex.Unlock()

	if c.ws == nil {
		return fmt.Errorf("WebSocket client not connected. Use Dial first")
	}

	err := make(chan error)
	c.messageChan <- message{
		messageType: messageType,
		data:        data,

		errChannel: err,
	}

	return <-err
}

// SendJSONMessage is used to send an arbitrary struct as JSON to the server
func (c *Client) SendJSONMessage(v interface{}) error {
	c.startStopMutex.Lock()
	defer c.startStopMutex.Unlock()

	if c.ws == nil {
		return fmt.Errorf("WebSocket client not connected. Use Dial first")
	}

	err := make(chan error)
	c.jSONMessageChan <- jSONMessage{
		v: v,

		errChannel: err,
	}

	return <-err
}
