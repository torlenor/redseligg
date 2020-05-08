package echoplugin

import (
	"fmt"
	"strings"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/utils"
)

// OnPost implements the hook from the Bot
func (p *EchoPlugin) OnPost(post model.Post) {
	msg := strings.Trim(post.Content, " ")
	if (!p.onlyOnWhisper || post.IsPrivate) && strings.HasPrefix(msg, "!echo ") {
		p.API.LogTrace(fmt.Sprintf("Echoing message back to Channel = %s, content = %s", post.Channel, utils.StripCmd(msg, "echo")))
		echo := post
		echo.Content = utils.StripCmd(msg, "echo")
		p.API.CreatePost(echo)
	}
}
