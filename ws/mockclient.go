package ws

// MockClient is a mock abstraction of a WebSocket client
type MockClient struct {
	WasDialCalled bool
	LastDialURL   string

	LastSendMessageType int
	LastSendMessageData []byte
	LastSendJSONMessage interface{}

	WasStopCalled bool

	ReturnError error
}

// Reset the MockClient
func (c *MockClient) Reset() {
	c.WasDialCalled = false
	c.LastDialURL = ""
	c.LastSendMessageType = 0
	c.LastSendMessageData = nil
	c.LastSendJSONMessage = nil
	c.WasStopCalled = false
}

// Dial is used to connect the ws Client to a WebSocket server
func (c *MockClient) Dial(wsURL string) error {
	c.LastDialURL = wsURL
	c.WasDialCalled = true

	return c.ReturnError
}

// Stop sending data to the WebSocket
func (c *MockClient) Stop() {
	c.WasStopCalled = true
}

// Close the websocket
func (c *MockClient) Close() error { return nil }

// ReadMessage can be used to read the next message from WebSocket.
// it blocks until somehing is received or the ws is closed.
func (c *MockClient) ReadMessage() (int, []byte, error) {
	return 0, nil, c.ReturnError
}

// SendMessage is used to send a message via the connected WebSocket to the server
func (c *MockClient) SendMessage(messageType int, data []byte) error {
	c.LastSendMessageType = messageType
	c.LastSendMessageData = data

	return c.ReturnError
}

// SendJSONMessage is used to send an arbitrary struct as JSON to the server
func (c *MockClient) SendJSONMessage(v interface{}) error {
	c.LastSendJSONMessage = v

	return c.ReturnError
}
