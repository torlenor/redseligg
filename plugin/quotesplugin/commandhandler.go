package quotesplugin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/storagemodels"
)

var now = time.Now

const (
	identFieldList = "list"

	helpText       = "Type `" + command + " add <your quote>` to add a new quote."
	helpTextRemove = "Type `" + command + " remove <your quote>` or `" + command + " remove (ID)` to remove a quote."
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

func generateIdentifier() string {
	return uuid.New().String()
}

func (p *QuotesPlugin) addQuoteIdentifierToList(identifier string) (int, error) {
	currentList := p.getQuotesList()
	currentList.UUIDs = append(currentList.UUIDs, identifier)
	s := p.getStorage()
	if s == nil {
		return 0, ErrNoValidStorage
	}

	err := s.StoreQuotesPluginQuotesList(p.BotID, p.PluginID, identFieldList, currentList)
	return len(currentList.UUIDs), err
}

func (p *QuotesPlugin) removeQuoteIdentifierToList(identifier string) (int, error) {
	currentList := p.getQuotesList()
	newList := storagemodels.QuotesPluginQuotesList{}
	for _, id := range currentList.UUIDs {
		if id != identifier {
			newList.UUIDs = append(newList.UUIDs, id)
		}
	}
	s := p.getStorage()
	if s == nil {
		return 0, ErrNoValidStorage
	}

	err := s.StoreQuotesPluginQuotesList(p.BotID, p.PluginID, identFieldList, newList)
	return len(currentList.UUIDs), err
}

func (p *QuotesPlugin) removeQuote(identifier string) error {
	s := p.getStorage()
	if s == nil {
		return ErrNoValidStorage
	}
	return s.DeleteQuotesPluginQuote(p.BotID, p.PluginID, identifier)
}

func (p *QuotesPlugin) getQuotesList() storagemodels.QuotesPluginQuotesList {
	var currentList storagemodels.QuotesPluginQuotesList

	s := p.getStorage()
	if s == nil {
		p.API.LogError(ErrNoValidStorage.Error())
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
		return storagemodels.QuotesPluginQuote{}, ErrNoValidStorage
	}

	return s.GetQuotesPluginQuote(p.BotID, p.PluginID, identifier)
}

func (p *QuotesPlugin) storeQuote(quote storagemodels.QuotesPluginQuote) int {
	id := generateIdentifier()

	s := p.getStorage()
	if s == nil {
		p.API.LogError(ErrNoValidStorage.Error())
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
func (p *QuotesPlugin) onCommandQuoteAdd(argument string, post model.Post) {
	quote := storagemodels.QuotesPluginQuote{
		Author:    post.User.Name,
		Added:     now(),
		AuthorID:  post.User.ID,
		ChannelID: post.ChannelID,
		Text:      argument,
	}

	num := p.storeQuote(quote)
	if num > 0 {
		p.returnMessage(post.ChannelID, fmt.Sprintf("Successfully added quote #%d", num-1))
	} else {
		p.returnMessage(post.ChannelID, "Error storing quote. Try again later!")
	}
}

// onCommandQuoteRemove removes a quote.
func (p *QuotesPlugin) onCommandQuoteRemove(argument string, post model.Post) {
	n, err := strconv.Atoi(argument)
	if err != nil {
		p.returnHelpRemove(post.ChannelID)
		return
	}

	currentList := p.getQuotesList()

	indexToRemove := 0
	if n <= len(currentList.UUIDs) {
		indexToRemove = n - 1
	} else {
		p.returnMessage(post.ChannelID, fmt.Sprintf("Quote #%d not found", n))
		return
	}

	p.removeQuoteIdentifierToList(currentList.UUIDs[indexToRemove])
	p.removeQuote(currentList.UUIDs[indexToRemove])

	p.returnMessage(post.ChannelID, fmt.Sprintf("Successfully removed quote #%d", n))
}

func (p *QuotesPlugin) onCommandQuote(argument string, post model.Post) {
	currentList := p.getQuotesList()
	if len(currentList.UUIDs) == 0 {
		p.returnMessage(post.ChannelID, "No quotes found. Use the command `"+command+" add <your quote>` to add a new one.")
		return
	}

	n := 0
	if len(argument) > 0 {
		var err error
		n, err = strconv.Atoi(argument)
		if err == nil && n <= len(currentList.UUIDs) {
			fmt.Printf("\nGetting quote %d\n", n)
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
