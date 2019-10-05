package ws

func (c *Client) writeMessage(messageType int, data []byte) error {
	return c.ws.WriteMessage(messageType, data)
}

func (c *Client) writeJSON(v interface{}) error {
	return c.ws.WriteJSON(v)
}

func (c *Client) wsWriter() {
	defer c.workersWG.Done()
	for {
		select {
		case val := <-c.messageChan:
			val.errChannel <- c.writeMessage(val.messageType, val.data)
		case val := <-c.jSONMessageChan:
			val.errChannel <- c.writeJSON(val.v)
		case <-c.stopWorkers:
			return
		}
	}
}
