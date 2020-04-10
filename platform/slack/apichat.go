package slack

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type messagePost struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type genericChatResponse struct {
	Ok      bool   `json:"ok"`
	Channel string `json:"channel"`
	TS      string `json:"ts"`
}

func (b *Bot) chatPostMessage(channel, content string) (genericChatResponse, error) {
	body, err := json.Marshal(
		messagePost{
			Channel: channel,
			Text:    content,
		},
	)
	if err != nil {
		errors.Wrap(err, "apiCall failed")
	}

	rawResponse, err := b.apiCallJON("/api/chat.postMessage", "POST", string(body))
	if err != nil {
		return genericChatResponse{}, errors.Wrap(err, "apiCall failed")
	}

	response := genericChatResponse{}
	err = json.Unmarshal(rawResponse.body, &response)

	if err == nil && !response.Ok {
		return genericChatResponse{}, fmt.Errorf("Error in sending message: Received not OK. Response was: %s", rawResponse.body)
	} else if err != nil {
		return genericChatResponse{}, err
	}

	return response, nil
}

type deleteBody struct {
	Channel string `json:"channel"`
	TS      string `json:"ts"`
}

func (b *Bot) chatDelete(channel, ts string) (genericChatResponse, error) {
	body, err := json.Marshal(
		deleteBody{
			Channel: channel,
			TS:      ts,
		},
	)
	if err != nil {
		errors.Wrap(err, "apiCall failed")
	}

	rawResponse, err := b.apiCallJON("/api/chat.delete", "POST", string(body))
	if err != nil {
		return genericChatResponse{}, errors.Wrap(err, "apiCall failed")
	}

	response := genericChatResponse{}
	err = json.Unmarshal(rawResponse.body, &response)

	if err == nil && !response.Ok {
		return genericChatResponse{}, fmt.Errorf("Error in deleting message: Received not OK. Response was: %s", rawResponse.body)
	} else if err != nil {
		return genericChatResponse{}, err
	}

	return response, nil
}

type updateBody struct {
	Channel string `json:"channel"`
	TS      string `json:"ts"`
	AsUser  bool   `json:"as_user"`

	Text string `json:"text"`
}

func (b *Bot) chatUpdate(channel, ts, content string) (genericChatResponse, error) {
	body, err := json.Marshal(
		updateBody{
			Channel: channel,
			TS:      ts,
			AsUser:  true,

			Text: content,
		},
	)
	if err != nil {
		errors.Wrap(err, "apiCall failed")
	}

	rawResponse, err := b.apiCallJON("/api/chat.update", "POST", string(body))
	if err != nil {
		return genericChatResponse{}, errors.Wrap(err, "apiCall failed")
	}

	response := genericChatResponse{}
	err = json.Unmarshal(rawResponse.body, &response)

	if err == nil && !response.Ok {
		return genericChatResponse{}, fmt.Errorf("Error in updating message: Received not OK. Response was: %s", rawResponse.body)
	} else if err != nil {
		return genericChatResponse{}, err
	}

	return response, nil
}
