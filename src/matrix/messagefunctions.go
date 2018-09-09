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
func (b Bot) sendRoomMessage(roomIdent string, content string) error {
	var roomID string
	if _, ok := b.knownRoomIDs[roomIdent]; ok {
		roomID = roomIdent
	} else if val, ok := b.knownRooms[roomIdent]; ok {
		roomID = val
	} else {
		log.Warnf("Unknown roomIdent %s. We will try to use it as a roomID", roomIdent)
		roomID = roomIdent
	}

	response, err := b.apiCall("/client/r0/rooms/"+roomID+"/send/m.room.message?access_token="+b.token, "POST", `{"msgtype":"m.text", "body":"`+content+`"}`, false)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}

	log.Debugln("send api response:", string(response))
	log.Debugf("Sent: MESSAGE to roomID = %s, Content = %s", roomID, content)

	return nil
}
