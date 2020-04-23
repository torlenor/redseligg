package quotesplugin

import (
	"regexp"
	"strings"

	"github.com/torlenor/abylebotter/model"
)

const (
	helpText       = "Type !quoteadd <your quote> to add a new quote."
	helpTextRemove = "Type `!quoteremove <your quote>` or !quoteremove (ID) to remove a quote."
)

func (p *QuotesPlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, helpText)
}

func (p *QuotesPlugin) returnHelpRemove(channelID string) {
	p.returnMessage(channelID, helpTextRemove)
}

func (p *QuotesPlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}

func (p *QuotesPlugin) extractRemoveID(fullText string) string {
	re := regexp.MustCompile(`!quoteremove (\(.*\))?`)
	const captureGroup = 1

	matches := re.FindAllStringSubmatch(fullText, -1)

	if matches == nil || len(matches) < 1 {
		return ""
	} else if len(matches) > 1 {
		p.API.LogWarn("QuotesPlugin: extractRemoveID matched more than one occurrence")
	}

	return strings.Trim(matches[0][captureGroup], " ")
}

// onCommandAddQuote adds a new quote.
func (p *QuotesPlugin) onCommandQuoteAdd(post model.Post) {
	// cont := strings.Split(post.Content, " ")
	// quoteText := strings.Join(cont[:1], " ")
	// TODO add quote to storage
}

// onCommandAddQuote adds a new quote.
func (p *QuotesPlugin) onCommandQuoteRemove(post model.Post) {
	removeID := p.extractRemoveID(post.Content)
	if len(removeID) == 0 {
		p.returnHelpRemove(post.ChannelID)
		return
	}
}

func (p *QuotesPlugin) onCommandRandomQuote(post model.Post) {
	// TODO return a random quote
	// p.returnMessage(post.ChannelID, "No vote running with that description in this channel. Use the !vote command to start a new one.")
}
