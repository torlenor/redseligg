package rollplugin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/utils"
)

// OnPost implements the hook from the Bot
func (p *RollPlugin) OnPost(post model.Post) {
	msg := strings.Trim(post.Content, " ")
	if strings.HasPrefix(msg, "!roll") {
		u := utils.StripCmd(msg, "roll")
		if len(msg) == len("!roll") && u == "!roll" {
			u = "100"
		}
		var response string
		num, err := strconv.Atoi(u)
		if err != nil {
			response = fmt.Sprintf("Not a number")
		} else if num <= 0 {
			response = fmt.Sprintf("Number must be > 0")
		} else {
			response = "<@" + post.UserID + "> rolled *" + strconv.Itoa(p.randomizer.random(num)) + "* in [0," + strconv.Itoa(num) + "]"
		}
		echo := post
		echo.Content = response
		p.API.CreatePost(echo)
	}
}
