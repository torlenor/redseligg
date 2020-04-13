package discord

import "fmt"

func getAbyleBotterEmojiFromDiscordEmoji(discordEmoji string) (string, error) {
	switch discordEmoji {

	case "0️⃣":
		return "zero", nil
	case "1️⃣":
		return "one", nil
	case "2️⃣":
		return "two", nil
	case "3️⃣":
		return "three", nil
	case "4️⃣":
		return "four", nil
	case "5️⃣":
		return "five", nil
	case "6️⃣":
		return "six", nil
	case "7️⃣":
		return "seven", nil
	case "8️⃣":
		return "eight", nil
	case "9️⃣":
		return "nine", nil
	case "🔟":
		return "ten", nil
	default:
		return discordEmoji, fmt.Errorf("Emoji not known")
	}
}
