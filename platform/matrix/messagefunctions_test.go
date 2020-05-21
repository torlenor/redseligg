package matrix

import (
	"testing"

	"github.com/torlenor/redseligg/commanddispatcher"
)

func TestSendRoomMessage(t *testing.T) {
	api := &mockAPI{server: "TEST_SERVER", authToken: "TEST_TOKEN"}
	api.reset()

	dispatcher := commanddispatcher.New("")

	bot, err := createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", dispatcher)
	if bot == nil || err != nil {
		t.Fatalf("Could not create Matrix Bot")
	}

	// successfully sending message with assumption of roomID
	err = bot.sendRoomMessage("_ROOMID_", "_MSGCONTENT_")
	if api.lastAPICallPath != `/client/r0/rooms/_ROOMID_/send/m.room.message` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `POST` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{"msgtype":"m.text", "body":"_MSGCONTENT_"}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err != nil {
		t.Fatalf("sending message failed even though it shouldn't")
	}

	// successfully sending message where room is mapped via NAME
	bot.addKnownRoom("_REALROOMID_", "_REALROOMNAME_")
	err = bot.sendRoomMessage("_REALROOMNAME_", "_MSGCONTENT_")
	if api.lastAPICallPath != `/client/r0/rooms/_REALROOMID_/send/m.room.message` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `POST` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{"msgtype":"m.text", "body":"_MSGCONTENT_"}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err != nil {
		t.Fatalf("sending message failed even though it shouldn't")
	}

	// successfully sending message where room is mapped via ID
	bot.addKnownRoom("_REALROOMID_", "_REALROOMNAME_")
	err = bot.sendRoomMessage("_REALROOMID_", "_MSGCONTENT_")
	if api.lastAPICallPath != `/client/r0/rooms/_REALROOMID_/send/m.room.message` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `POST` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{"msgtype":"m.text", "body":"_MSGCONTENT_"}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err != nil {
		t.Fatalf("sending message failed even though it shouldn't")
	}

	// return error when api call fails
	api.reset()
	api.letAPICallFail = true
	err = bot.sendRoomMessage("_ROOMID_", "_MSGCONTENT_")
	if api.lastAPICallPath != `/client/r0/rooms/_ROOMID_/send/m.room.message` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `POST` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{"msgtype":"m.text", "body":"_MSGCONTENT_"}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err == nil {
		t.Fatalf("sending message not failed even though mock api call failed")
	}
}

func TestSendWhisper(t *testing.T) {
	// For now whisper is the same as room message, at least as we know
	// from Matrix API docs. If it ever changes these tests will fail
	// to remind us
	api := &mockAPI{server: "TEST_SERVER", authToken: "TEST_TOKEN"}
	api.reset()

	dispatcher := commanddispatcher.New("")

	bot, err := createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", dispatcher)
	if bot == nil || err != nil {
		t.Fatalf("Could not create Matrix Bot")
	}

	// successfully sending message with assumption of roomID
	err = bot.sendWhisper("_ROOMID_", "_MSGCONTENT_")
	if api.lastAPICallPath != `/client/r0/rooms/_ROOMID_/send/m.room.message` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `POST` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{"msgtype":"m.text", "body":"_MSGCONTENT_"}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err != nil {
		t.Fatalf("sending message failed even though it shouldn't")
	}

	// successfully sending message where room is mapped via NAME
	bot.addKnownRoom("_REALROOMID_", "_REALROOMNAME_")
	err = bot.sendWhisper("_REALROOMNAME_", "_MSGCONTENT_")
	if api.lastAPICallPath != `/client/r0/rooms/_REALROOMID_/send/m.room.message` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `POST` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{"msgtype":"m.text", "body":"_MSGCONTENT_"}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err != nil {
		t.Fatalf("sending message failed even though it shouldn't")
	}

	// successfully sending message where room is mapped via ID
	bot.addKnownRoom("_REALROOMID_", "_REALROOMNAME_")
	err = bot.sendWhisper("_REALROOMID_", "_MSGCONTENT_")
	if api.lastAPICallPath != `/client/r0/rooms/_REALROOMID_/send/m.room.message` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `POST` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{"msgtype":"m.text", "body":"_MSGCONTENT_"}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err != nil {
		t.Fatalf("sending message failed even though it shouldn't")
	}

	// return error when api call fails
	api.reset()
	api.letAPICallFail = true
	err = bot.sendWhisper("_ROOMID_", "_MSGCONTENT_")
	if api.lastAPICallPath != `/client/r0/rooms/_ROOMID_/send/m.room.message` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `POST` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{"msgtype":"m.text", "body":"_MSGCONTENT_"}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err == nil {
		t.Fatalf("sending message not failed even though mock api call failed")
	}
}
