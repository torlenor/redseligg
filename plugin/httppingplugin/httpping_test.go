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

func TestHTTPPingPlugin_OnPost(t *testing.T) {
	assert := assert.New(t)

	mockHTTPGetter := httpGetter{}
	httpGet = mockHTTPGetter.get

	p := HTTPPingPlugin{}
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)
	assert.Equal(false, mockHTTPGetter.WasGetCalled)

	mockHTTPGetter.reset()
	api.Reset()
	postToPlugin.Content = "!httpping"
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

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
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	mockHTTPGetter.reset()
	api.Reset()
	postToPlugin.Content = "!httpping http://validurl.com"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "SUCCESS. Request took ",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
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
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
