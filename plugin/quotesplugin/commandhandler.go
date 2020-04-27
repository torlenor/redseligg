package quotesplugin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/storagemodels"
)

var now = time.Now

const (
	identFieldList = "list"

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

func (p *QuotesPlugin) addQuoteIdentifierToList(identifier string) (int, error) {
	currentList := p.getQuotesList()
	currentList.UUIDs = append(currentList.UUIDs, identifier)
	s := p.getStorage()
	if s == nil {
		return 0, fmt.Errorf("Not valid storage set")
	}

	err := s.StoreQuotesPluginQuotesList(p.BotID, p.PluginID, identFieldList, currentList)
	return len(currentList.UUIDs), err
}

func (p *QuotesPlugin) getQuotesList() storagemodels.QuotesPluginQuotesList {
	var currentList storagemodels.QuotesPluginQuotesList

	s := p.getStorage()
	if s == nil {
		p.API.LogError("Not valid storage set")
		return currentList
	}
	var err error
	currentList, err = s.GetQuotesPluginQuotesList(p.BotID, p.PluginID, identFieldList)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Could not get QuotesList: %s", err))
	}

	return currentList
}

func (p *QuotesPlugin) getQuote(identifier string) (storagemodels.QuotesPluginQuote, error) {
	s := p.getStorage()
	if s == nil {
		return storagemodels.QuotesPluginQuote{}, fmt.Errorf("Not valid storage set")
	}

	return s.GetQuotesPluginQuote(p.BotID, p.PluginID, identifier)
}

func (p *QuotesPlugin) storeQuote(quote storagemodels.QuotesPluginQuote) int {
	id := generateIdentifier()

	s := p.getStorage()
	if s == nil {
		p.API.LogError("Not valid storage set")
		return 0
	}

	if err := s.StoreQuotesPluginQuote(p.BotID, p.PluginID, id, quote); err == nil {
		if num, err := p.addQuoteIdentifierToList(id); err != nil {
			p.API.LogError(fmt.Sprintf("Error storing quotes list: %s", err))
		} else {
			return num
		}
	} else {
		p.API.LogError(fmt.Sprintf("Error storing quote: %s", err))
	}

	return 0
}

// onCommandAddQuote adds a new quote.
func (p *QuotesPlugin) onCommandQuoteAdd(post model.Post) {
	cont := strings.Split(post.Content, " ")
	quoteText := strings.Join(cont[1:], " ")

	quote := storagemodels.QuotesPluginQuote{
		Author:    post.User.Name,
		Added:     now(),
		AuthorID:  post.User.ID,
		ChannelID: post.ChannelID,
		Text:      quoteText,
	}

	num := p.storeQuote(quote)
	if num > 0 {
		p.returnMessage(post.ChannelID, fmt.Sprintf("Successfully added quote #%d", num))
	} else {
		p.returnMessage(post.ChannelID, "Error storing quote. Try again later!")
	}
}

// onCommandQuoteRemove removes a quote.
func (p *QuotesPlugin) onCommandQuoteRemove(post model.Post) {
	removeID := p.extractRemoveID(post.Content)
	if len(removeID) == 0 {
		p.returnHelpRemove(post.ChannelID)
		return
	}

	// TODO
}

func (p *QuotesPlugin) onCommandQuote(post model.Post) {
	cont := strings.Split(post.Content, " ")

	currentList := p.getQuotesList()
	if len(currentList.UUIDs) == 0 {
		p.returnMessage(post.ChannelID, "No quotes found. Use the command `!quoteadd <your quote>` to add a new one.")
		return
	}

	n := 0
	if len(cont) == 2 {
		var err error
		n, err = strconv.Atoi(cont[1])
		if err == nil && n <= len(currentList.UUIDs) {
			n = n - 1
		} else {
			n = p.randomizer.Intn(len(currentList.UUIDs))
		}
	} else {
		n = p.randomizer.Intn(len(currentList.UUIDs))
	}

	if quote, err := p.getQuote(currentList.UUIDs[n]); err == nil {
		p.returnMessage(post.ChannelID, fmt.Sprintf("%d. %s", n+1, quote))
	} else {
		p.API.LogError(fmt.Sprintf("Could not receive quote with id %s: %s", currentList.UUIDs[n], err))
	}

}
