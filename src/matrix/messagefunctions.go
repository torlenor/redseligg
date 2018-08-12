package matrix

import (
	"log"

	"github.com/pkg/errors"
)

func (b Bot) sendWhisper(userID string, content string) error {
	return nil
}

func (b Bot) sendRoomMessage(roomID string, content string) error {
	response, err := b.apiCall("/client/r0/rooms/"+roomID+"/send/m.room.message?access_token="+b.token, "POST", `{"msgtype":"m.text", "body":"`+content+`"}`)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}
	log.Println(string(response))

	log.Printf(logPrefix+"Sent: MESSAGE to roomID = %s, Content = %s", roomID, content)
	return nil
}
