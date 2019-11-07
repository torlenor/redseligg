package pool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/torlenor/abylebotter/utils"
)

type getBotsResponse struct {
	Bots []string `json:"bots"`
}

func (b *BotPool) getBotsEndpoint(w http.ResponseWriter, r *http.Request) {
	response := getBotsResponse{
		Bots: b.GetBotIDs(),
	}

	out, err := json.Marshal(response)
	if err != nil {
		b.log.Errorln(err)
		http.Error(w, utils.GenerateErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Server error, try again later")), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(out))
}
