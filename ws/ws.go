package ws

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/torlenor/abylebotter/logging"
)

// message is used to send a []byte message of messageType over the WebSocket
type message struct {
	messageType int
	data        []byte

	errChannel chan error
}

// jSONMessage is used to send an arbitrary struct V via JSON message over the WebSocket
type jSONMessage struct {
	v interface{}

	errChannel chan error
}

// Client is an abstraction of a WebSocket client
type Client struct {
	log *logrus.Entry

	messageChan     chan message
	jSONMessageChan chan jSONMessage

	ws *websocket.Conn

	startStopMutex sync.Mutex

	workersWG   sync.WaitGroup
	stopWorkers chan struct{}
}

// NewClient prepares a new ws Client
func NewClient() *Client {
	log := logging.Get("WSClient")
	log.Debug("WS Client is CREATING itself")

	c := Client{
		log: log,

		messageChan:     make(chan message),
		jSONMessageChan: make(chan jSONMessage),

		stopWorkers: make(chan struct{}),
	}
	return &c
}

// Dial is used to connect the ws Client to a WebSocket server
func (c *Client) Dial(wsURL string) error {
	c.startStopMutex.Lock()
	defer c.startStopMutex.Unlock()

	var err error
	c.ws, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}

	c.workersWG.Add(1)
	go c.wsWriter()

	return nil
}

// Stop sending data to the WebSocket
func (c *Client) Stop() {
	c.startStopMutex.Lock()

	close(c.stopWorkers)
	c.workersWG.Wait()
	c.ws = nil

	c.startStopMutex.Unlock()
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
