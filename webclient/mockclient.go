package webclient

// MockClient implements mocked features of a web client
type MockClient struct {
	LastCallPath   string
	LastCallMethod string
	LastCallBody   string

	ReturnOnCall      APIResponse
	ReturnOnCallError error
}

// NewMock creates a new web client mock
func NewMock() *MockClient {
	return &MockClient{}
}

// Reset the MockClient
func (c *MockClient) Reset() {
	c.LastCallPath = ""
	c.LastCallMethod = ""
	c.LastCallBody = ""
}

// Call executes a web request (mock)
func (c *MockClient) Call(path string, method string, body string) (APIResponse, error) {
	c.LastCallPath = path
	c.LastCallMethod = method
	c.LastCallBody = body
	return c.ReturnOnCall, c.ReturnOnCallError
}
