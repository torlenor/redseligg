package pool

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/torlenor/abylebotter/utils"
)

// GetBotsResponse is the response for bots endpoint
type GetBotsResponse struct {
	Bots []string `json:"bots"`
}

func (b *BotPool) getBotsEndpoint(w http.ResponseWriter, r *http.Request) {
	response := GetBotsResponse{
		Bots: b.GetBotIDs(),
	}

	out, err := json.Marshal(response)
	if err != nil {
		b.log.Errorln(err)
		http.Error(w, utils.GenerateErrorResponse(fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))
}

func (b *BotPool) postBotsEndpoint(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, utils.GenerateErrorResponse(err.Error()), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, utils.GenerateErrorResponse(fmt.Sprintf("Invalid body received: JSON invalid")), http.StatusBadRequest)
		return
	}

	var botID string
	if val, ok := data["botId"].(string); ok {
		botID = val
	} else {
		http.Error(w, utils.GenerateErrorResponse(fmt.Sprintf("Invalid body received: 'botId' required")), http.StatusBadRequest)
		return
	}

	err = b.AddViaID(botID)
	if err != nil {
		http.Error(w, utils.GenerateErrorResponse(fmt.Sprintf("Not able to add Bot with ID %s: %s", botID, err)), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (b *BotPool) deleteBotEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	botID := vars["botId"]

	b.RemoveViaID(botID)

	w.WriteHeader(http.StatusOK)
}

func (b *BotPool) getBotEndPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	botID := vars["botId"]

	if _, ok := b.bots[botID]; !ok {
		http.Error(w, utils.GenerateErrorResponse(fmt.Sprintf("Bot ID %s unknown", botID)), http.StatusBadRequest)
		return
	}

	info := b.bots[botID].GetInfo()
	info.BotID = botID

	out, err := json.Marshal(info)
	if err != nil {
		b.log.Errorln(err)
		http.Error(w, utils.GenerateErrorResponse(fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))
}
