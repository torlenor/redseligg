package matrix

import (
	"github.com/pkg/errors"
)

// The sendWhisper function sends a whisper to the user with userID
//
// Note: sending a whisper is the same as sending a message
// just with a room id which belongs to just the two
// participants
func (b Bot) sendWhisper(userID string, content string) error {
	return b.sendRoomMessage(userID, content)
}

// The sendRoomMessage function sends a message to the room with
// the ID roomID.
func (b Bot) sendRoomMessage(roomID string, content string) error {
	response, err := b.apiCall("/client/r0/rooms/"+roomID+"/send/m.room.message?access_token="+b.token, "POST", `{"msgtype":"m.text", "body":"`+content+`"}`, false)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}
	log.Println(string(response))

	log.Printf("Sent: MESSAGE to roomID = %s, Content = %s", roomID, content)
	return nil
}
