package quotesplugin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/torlenor/abylebotter/model"
)

const (
	LIST_IDENTIFIER = "list"

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

func generateIdentifier() string {
	return uuid.New().String()
}

type quotesList []string

func (p *QuotesPlugin) addQuoteIdentifierToList(identifier string) error {
	var currentList quotesList

	var data interface{}
	var err error
	if data, err = p.Get(LIST_IDENTIFIER); err != nil {
		p.API.LogInfo("QuotesPlugin: No stored Quotes list found, creating a new one")
	} else {
		var ok bool
		if currentList, ok = data.(quotesList); !ok {
			currentList = quotesList{}
		}
	}

	currentList = append(currentList, identifier)

	return p.Store(LIST_IDENTIFIER, currentList)
}

func (p *QuotesPlugin) getQuotesList() quotesList {
	var currentList quotesList

	var data interface{}
	var err error
	if data, err = p.Get(LIST_IDENTIFIER); err != nil {
		return quotesList{}
	}

	var ok bool
	if currentList, ok = data.(quotesList); !ok {
		return quotesList{}
	}

	return currentList
}

func (p *QuotesPlugin) getQuote(identifier string) (string, error) {
	var quoteText string

	var data interface{}
	var err error
	if data, err = p.Get(identifier); err != nil {
		return "", fmt.Errorf("Error receiving quote with id '%s': %s", identifier, err)
	}

	var ok bool
	if quoteText, ok = data.(string); !ok {
		return "", fmt.Errorf("Error receiving quote with id '%s': Not a quote", identifier)
	}

	return quoteText, nil
}

func (p *QuotesPlugin) storeQuote(author model.User, channel string, channelID string, quote string) {
	id := generateIdentifier()

	if err := p.Store(id, quote); err == nil {
		if err := p.addQuoteIdentifierToList(id); err != nil {
			p.API.LogError(fmt.Sprintf("Error storing quotes list: %s", err))
		}
	} else {
		p.API.LogError(fmt.Sprintf("Error storing quote: %s", err))
	}
}

// onCommandAddQuote adds a new quote.
func (p *QuotesPlugin) onCommandQuoteAdd(post model.Post) {
	cont := strings.Split(post.Content, " ")
	quoteText := strings.Join(cont[1:], " ")

	p.storeQuote(post.User, post.Channel, post.ChannelID, quoteText)
}

// onCommandAddQuote adds a new quote.
func (p *QuotesPlugin) onCommandQuoteRemove(post model.Post) {
	removeID := p.extractRemoveID(post.Content)
	if len(removeID) == 0 {
		p.returnHelpRemove(post.ChannelID)
		return
	}
}

func (p *QuotesPlugin) onCommandQuote(post model.Post) {
	currentList := p.getQuotesList()
	if len(currentList) == 0 {
		p.returnMessage(post.ChannelID, "No quotes found. Use the command `!quoteadd <your quote>` to add a new one.")
		return
	}

	id := p.randomizer.Intn(len(currentList))

	if quote, err := p.getQuote(currentList[id]); err == nil {
		p.returnMessage(post.ChannelID, quote)
	}

}
