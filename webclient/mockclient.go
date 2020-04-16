package webclient

import "fmt"

// MockClient implements mocked features of a web client
type MockClient struct{}

// NewMock creates a new web client mock
func NewMock() *MockClient {
	return &MockClient{}
}

// Call executes a web request (mock)
func (c *MockClient) Call(path string, method string, body string) (APIResponse, error) {
	return APIResponse{}, fmt.Errorf("not implemented")
}
