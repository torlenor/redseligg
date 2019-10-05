package ws

func (c *Client) writeMessage(messageType int, data []byte) error {
	c.log.Debug("Sending message")
	return c.ws.WriteMessage(messageType, data)
}

func (c *Client) writeJSON(v interface{}) error {
	c.log.Debug("Sending json")
	return c.ws.WriteJSON(v)
}

func (c *Client) wsWriter() {
	c.log.Debug("Starting wsWriter")
	defer c.workersWG.Done()
	for {
		select {
		case val := <-c.messageChan:
			c.log.Debug("Received message to send...")
			val.errChannel <- c.writeMessage(val.messageType, val.data)
		case val := <-c.jSONMessageChan:
			c.log.Debug("Received JSON to send...")
			val.errChannel <- c.writeJSON(val.v)
		case <-c.stopWorkers:
			return
		}
	}
}
