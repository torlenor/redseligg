package matrix

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func checkEventsOK(repsonse []byte) bool {
	// TODO check if response is containing any errors
	return true
}

func (b *Bot) handleJoinRooms(rooms []room) {
	for _, room := range rooms {
		for _, event := range room.State.Events {
			if event.Type == "m.room.name" {
				b.addKnownRoom(room.RoomID, event.Content.Name)
			}
		}

		for _, event := range room.Timeline.Events {
			if event.Type == "m.room.message" {
				log.Debugf("Received room message from User: %s, Content: %s, MsgType: %s", event.Sender,
					event.Content.Body, event.Content.Msgtype)
			}
		}

	}
}

func (b *Bot) handleLeaveRooms(rooms []room) {
	for _, room := range rooms {
		response, err := b.api.call("/client/r0/rooms/"+room.RoomID+"/forget", "POST", `{}`, true)
		if err != nil {
			log.Errorf("leave room failed, err = %s, response = %s", err, response)
			return
		}
		b.removeKnownRoomFromID(room.RoomID)
	}
}

func (b *Bot) handleInviteRooms(rooms []room) {
	for _, room := range rooms {
		response, err := b.api.call("/client/r0/rooms/"+room.RoomID+"/join", "POST", `{}`, true)
		if err != nil {
			log.Errorln("join room failed:", err)
		}
		log.Println(string(response))
	}
}

func (b *Bot) callSync() error {
	var response []byte
	var err error
	if len(b.nextBatch) == 0 {
		response, err = b.api.call("/client/r0/sync?filter={\"room\":{\"timeline\":{\"limit\":1}}}", "GET", `{}`, true)
		if err != nil {
			log.Println("UNHANDELED ERROR: ", err)
			return err
		}
	} else {
		response, err = b.api.call("/client/r0/sync?since="+b.nextBatch, "GET", `{}`, true)
		if err != nil {
			log.Println("UNHANDELED ERROR: ", err)
			return err
		}
	}

	if !checkEventsOK(response) {
		return errors.New("Failed to get Events")
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response, &data); err != nil {
		log.Println("UNHANDELED ERROR: ", err)
		return err
	}

	sr, err := syncResponseFromMap(data)
	if err != nil {
		log.Println("UNHANDELED ERROR: ", err)
		return err
	}

	b.nextBatch = sr.NextBatch

	b.handleJoinRooms(sr.Rooms.Join)
	b.handleLeaveRooms(sr.Rooms.Leave)
	b.handleInviteRooms(sr.Rooms.Invite)

	return err
}
