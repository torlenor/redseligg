package ws

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
