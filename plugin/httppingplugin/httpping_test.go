package httppingplugin

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/plugin"

	"github.com/stretchr/testify/assert"
)

type httpGetter struct {
	WasGetCalled      bool
	LastURL           string
	ShouldReturnError bool
}

func (h *httpGetter) reset() {
	h.WasGetCalled = false
	h.LastURL = ""
	h.ShouldReturnError = false
}

func (h *httpGetter) get(url string) (resp *http.Response, err error) {
	if h.ShouldReturnError {
		return nil, fmt.Errorf("Some error")
	}

	return &http.Response{
		StatusCode: 200,
	}, nil
}

func TestHTTPPingPlugin_OnCommand(t *testing.T) {
	assert := assert.New(t)

	mockHTTPGetter := httpGetter{}
	httpGet = mockHTTPGetter.get

	p := HTTPPingPlugin{}
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	command := "httpping"

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	mockHTTPGetter.reset()
	api.Reset()
	postToPlugin.Content = "!httpping not a valid url"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "FAIL (Not a valid url).",
		IsPrivate: false,
	}
	p.OnCommand(command, "not a valid url", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	mockHTTPGetter.reset()
	api.Reset()
	postToPlugin.Content = "!" + command + "http://validurl.com"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "SUCCESS. Request took ",
		IsPrivate: false,
	}
	p.OnCommand(command, "http://validurl.com", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	// we do not want to check the time it took, because this is not mocked
	api.LastCreatePostPost.Content = api.LastCreatePostPost.Content[:len(expectedPostFromPlugin.Content)]
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	mockHTTPGetter.reset()
	mockHTTPGetter.ShouldReturnError = true
	api.Reset()
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "FAIL (Error pinging the url: Some error).",
		IsPrivate: false,
	}
	p.OnCommand(command, "http://validurl.com", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
