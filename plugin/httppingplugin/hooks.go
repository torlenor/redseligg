package httppingplugin

import (
	"fmt"
	"strings"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/utils"
)

// OnPost implements the hook from the Bot
func (p *HTTPPingPlugin) OnPost(post model.Post) {
	msg := strings.Trim(post.Content, " ")
	if strings.HasPrefix(msg, "!httpping ") {
		u := utils.StripCmd(msg, "httpping")

		timeMs, err := httpPing(u)

		var response string
		if err != nil {
			response = fmt.Sprintf("FAIL (%s).", err)
		} else {
			response = fmt.Sprintf("SUCCESS. Request took %d ms", timeMs)
		}

		pingReply := post
		pingReply.Content = response
		p.API.CreatePost(pingReply)
	}
}
