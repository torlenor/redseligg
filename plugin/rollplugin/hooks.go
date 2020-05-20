package rollplugin

import (
	"fmt"
	"strconv"

	"github.com/torlenor/redseligg/model"
)

// OnRun implements the hool from the Boot
func (p *RollPlugin) OnRun() {
	p.API.RegisterCommand(p, "roll")
}

// OnCommand implements the hook from the Bot
func (p *RollPlugin) OnCommand(cmd string, content string, post model.Post) {
	u := content
	if len(u) == 0 {
		u = "100"
	}
	var response string
	num, err := strconv.Atoi(u)
	if err != nil {
		response = fmt.Sprintf("Not a number")
	} else if num <= 0 {
		response = fmt.Sprintf("Number must be > 0")
	} else {
		response = "<@" + post.User.ID + "> rolled *" + strconv.Itoa(p.randomizer.random(num)) + "* in [0," + strconv.Itoa(num) + "]"
	}
	echo := post
	echo.Content = response
	p.API.CreatePost(echo)
}
