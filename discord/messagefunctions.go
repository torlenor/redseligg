package discord

import (
	"encoding/json"
	"log"
)

func (b Bot) sendWhisper(snowflakeID string, content string) {
	response, _ := b.apiCall("/users/@me/channels", "POST", `{"recipient_id": "`+snowflakeID+`"}`)

	var channelResponseData map[string]interface{}
	if err := json.Unmarshal(response, &channelResponseData); err != nil {
		log.Printf("DiscordBot: Error: Response was: %s\n", response)
		return
	}

	response, err := b.apiCall("/channels/"+channelResponseData["id"].(string)+"/messages", "POST", `{"content": "`+content+`"}`)
	if err != nil {
		log.Printf("DiscordBot: Error: Response was: %s\n", response)
		return
	}
	log.Printf("DiscordBot: Sent: WHISPER to UserID = %s, Content = %s", snowflakeID, content)
}

func (b Bot) sendMessage(channelID string, content string) {
	response, err := b.apiCall("/channels/"+channelID+"/messages", "POST", `{"content": "`+content+`"}`)
	if err != nil {
		log.Printf("DiscordBot: Error: Response was: %s\n", response)
		return
	}
	log.Printf("DiscordBot: Sent: MESSAGE to ChannelID = %s, Content = %s", channelID, content)
}
