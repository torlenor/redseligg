package ws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockWSClient struct {
	retError error

	retMessageType int
	retMessage     []byte

	lastMessageType int
	lastData        []byte
	lastJSON        interface{}
}

func (m *mockWSClient) ReadMessage() (messageType int, p []byte, err error) {
	if m.retError != nil {
		return 0, nil, m.retError
	}
	return m.retMessageType, m.retMessage, nil
}

func (m *mockWSClient) WriteMessage(messageType int, data []byte) error {
	m.lastMessageType = messageType
	m.lastData = data

	if m.retError != nil {
		return m.retError
	}
	return nil
}

func (m *mockWSClient) WriteJSON(v interface{}) error {
	m.lastJSON = v

	if m.retError != nil {
		return m.retError
	}
	return nil
}

func (m *mockWSClient) Close() error { return nil }

type mockDialer struct {
	wsClient wsClient

	retError error

	dialCalled    bool
	lastURLCalled string
}

func (m *mockDialer) Dial(wsURL string) (wsClient, error) {
	m.dialCalled = true
	m.lastURLCalled = wsURL

	if m.retError != nil {
		return nil, m.retError
	}

	return m.wsClient, nil
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	client := NewClient()

	assert.NotNil(client.log)
	assert.NotNil(client.messageChan)
	assert.NotNil(client.jSONMessageChan)

	assert.Nil(client.ws)
	assert.Nil(client.stopWorkers)
}

func TestClient_Dial(t *testing.T) {
	// If there is no error, everything should go fine
	{
		mockDialer := mockDialer{
			wsClient: &mockWSClient{},
		}
		dialer = mockDialer.Dial

		assert := assert.New(t)

		client := NewClient()

		assert.Equal(false, mockDialer.dialCalled)
		assert.Equal("", mockDialer.lastURLCalled)

		assert.NoError(client.Dial("SOME_URL"))
		assert.Equal(true, mockDialer.dialCalled)
		assert.Equal("SOME_URL", mockDialer.lastURLCalled)

		assert.NotNil(client.ws)
		assert.NotNil(client.stopWorkers)

		assert.NotPanics(func() { client.Close() })

		assert.Nil(client.ws)
		assert.Nil(client.stopWorkers)
	}

	// If there is an error in dialing Dial should fail
	{
		mockDialer := mockDialer{
			wsClient: &mockWSClient{},
			retError: fmt.Errorf("Some Error"),
		}
		dialer = mockDialer.Dial

		assert := assert.New(t)

		client := NewClient()

		assert.Equal(false, mockDialer.dialCalled)
		assert.Equal("", mockDialer.lastURLCalled)

		assert.Error(client.Dial("SOME_URL"))
		assert.Equal(true, mockDialer.dialCalled)
		assert.Equal("SOME_URL", mockDialer.lastURLCalled)

		assert.Nil(client.ws)
		assert.Nil(client.stopWorkers)

		assert.NotPanics(func() { client.Close() })
	}

}

func TestClient_ReadMessage(t *testing.T) {
	assert := assert.New(t)

	mockWsClient := &mockWSClient{}
	mockDialer := mockDialer{
		wsClient: mockWsClient,
	}
	dialer = mockDialer.Dial
	client := NewClient()
	assert.NoError(client.Dial("SOME_URL"))

	mockWsClient.retMessageType = 2
	mockWsClient.retMessage = []byte("SOME MESSAGE")

	actualMessageType, actualMessage, err := client.ReadMessage()
	assert.NoError(err)
	assert.Equal(mockWsClient.retMessageType, actualMessageType)
	assert.Equal(mockWsClient.retMessage, actualMessage)

	mockWsClient.retError = fmt.Errorf("SOME ERROR")
	_, _, err = client.ReadMessage()
	assert.Error(err)

	assert.NotPanics(func() { client.Close() })
}

func TestClient_SendMessage(t *testing.T) {
	assert := assert.New(t)

	mockWsClient := &mockWSClient{}
	mockDialer := mockDialer{
		wsClient: mockWsClient,
	}
	dialer = mockDialer.Dial
	client := NewClient()

	err := client.SendMessage(4, []byte("MESSAGE"))
	assert.Errorf(err, "WebSocket client not connected. Use Dial first")

	assert.NoError(client.Dial("SOME_URL"))

	err = client.SendMessage(4, []byte("MESSAGE"))
	assert.NoError(err)
	assert.Equal(4, mockWsClient.lastMessageType)
	assert.Equal([]byte("MESSAGE"), mockWsClient.lastData)

	mockWsClient.retError = fmt.Errorf("ERROR")
	err = client.SendMessage(5, []byte("OTHER MESSAGE"))
	assert.Error(err)
	assert.Equal(5, mockWsClient.lastMessageType)
	assert.Equal([]byte("OTHER MESSAGE"), mockWsClient.lastData)

	assert.NotPanics(func() { client.Close() })
}

func TestClient_SendJSONMessage(t *testing.T) {
	assert := assert.New(t)

	mockWsClient := &mockWSClient{}
	mockDialer := mockDialer{
		wsClient: mockWsClient,
	}
	dialer = mockDialer.Dial
	client := NewClient()

	err := client.SendJSONMessage([]byte("MESSAGE"))
	assert.Errorf(err, "WebSocket client not connected. Use Dial first")

	assert.NoError(client.Dial("SOME_URL"))

	err = client.SendJSONMessage([]byte("MESSAGE"))
	assert.NoError(err)
	assert.Equal([]byte("MESSAGE"), mockWsClient.lastJSON)

	mockWsClient.retError = fmt.Errorf("ERROR")
	err = client.SendJSONMessage([]byte("OTHER MESSAGE"))
	assert.Error(err)
	assert.Equal([]byte("OTHER MESSAGE"), mockWsClient.lastJSON)

	assert.NotPanics(func() { client.Close() })
}
