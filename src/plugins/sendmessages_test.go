package plugins

import (
	"events"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSendMessages(t *testing.T) {
	sendMessages, err := CreateSendMessagesPlugin()

	if err != nil {
		t.Fatalf("Could not create SendMessagesPlugin")
	}

	var bot mockBot
	bot.receiveMessageChan = make(chan events.ReceiveMessage)
	bot.sendMessageChan = make(chan events.SendMessage)
	bot.commandChan = make(chan events.Command)

	go bot.startSendChannelReceiver()
	go bot.startCommandChannelReceiver()

	bot.reset()

	err = sendMessages.ConnectChannels(bot.receiveMessageChan, bot.sendMessageChan, bot.commandChan)
	if err != nil {
		t.Fatalf("Error connecting channels")
	}

	if sendMessages.IsStarted() != false {
		t.Fatalf("SendMessagesPlugin should not have reported started")
	}

	sendMessages.Start()

	if sendMessages.IsStarted() != true {
		t.Fatalf("SendMessagesPlugin should report started now")
	}

	bot.reset()

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/plugins/sendmessages", strings.NewReader(`{"Ident":"ROOM_NAME","Content":"TEST_MESSAGE"}`))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sendMessages.sendMessage)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	time.Sleep(time.Millisecond * 100)

	if bot.lastSendMessage.Ident != "ROOM_NAME" ||
		bot.lastSendMessage.Content != "TEST_MESSAGE" ||
		bot.lastSendMessage.Type != events.MESSAGE {
		t.Fatalf("SendMessagesPlugin did not relay the send request to the bot")
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"sent":true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

	bot.reset()

	// Test faulty request
	req, err = http.NewRequest("POST", "/plugins/sendmessages", strings.NewReader(`{"Idsdfsfent":"ROOM_NAME","Condsfsdftent":"TEST_MESSAGE"}`))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(sendMessages.sendMessage)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	time.Sleep(time.Millisecond * 100)

	if bot.lastSendMessage.Ident != "" ||
		bot.lastSendMessage.Content != "" {
		t.Fatalf("SendMessagesPlugin relayed something to the bot even though the request was invalid")
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = `{"sent":false, "error":"Invalid Request"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

	sendMessages.Stop()

	if sendMessages.IsStarted() != false {
		t.Fatalf("SendMessagesPlugin should be stopped now")
	}
}
