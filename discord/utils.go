package discord

func extractOwnSnowflakeID(data map[string]interface{}) string {
	userData := data["d"].(map[string]interface{})["user"].(map[string]interface{})
	return userData["id"].(string)
}

func extractSenderSnowflakeID(data map[string]interface{}) string {
	authorData := data["d"].(map[string]interface{})["author"].(map[string]interface{})
	return authorData["id"].(string)
}
