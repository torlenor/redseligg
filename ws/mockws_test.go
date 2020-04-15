package ws

// MockClient is a mock abstraction of a WebSocket client
type MockClient struct {
	WasDialCalled bool
	LastDialURL   string

	WasStopCalled bool
}

// Dial is used to connect the ws Client to a WebSocket server
func (c *MockClient) Dial(wsURL string) error {
	c.LastDialURL = wsURL
	c.WasDialCalled = true

	return nil
}

// Stop sending data to the WebSocket
func (c *MockClient) Stop() {
	c.WasStopCalled = true
}

// ReadMessage can be used to read the next message from WebSocket.
// it blocks until somehing is received or the ws is closed.
func (c *MockClient) ReadMessage() (int, []byte, error) {
	return 0, nil, nil
}

// SendMessage is used to send a message via the connected WebSocket to the server
func (c *MockClient) SendMessage(messageType int, data []byte) error {
	return nil
}

// SendJSONMessage is used to send an arbitrary struct as JSON to the server
func (c *MockClient) SendJSONMessage(v interface{}) error {
	return nil
}
