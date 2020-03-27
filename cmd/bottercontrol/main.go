package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/torlenor/abylebotter/pool"
)

/**
 * Version should be set while build using ldflags (see Makefile)
 */
var version string

func getBots(url string) {
	data, code, err := realAPICall(url, "/bots", "GET", "")
	if err != nil {
		fmt.Printf("Error in getBots: %s, %s", url+"/bots", err)
		return
	}

	if code != 200 {
		fmt.Printf("Error in getBots: %s, unknown StatusCode: %d", url+"/bots", code)
		return
	}

	bots := pool.GetBotsResponse{}
	err = json.Unmarshal(data, &bots)
	if err != nil {
		fmt.Printf("Error in getBots: %s, %s", url+"/bots", err)
		return
	}

	fmt.Printf("Running bots @ %s:\n", url)
	for _, bot := range bots.Bots {
		fmt.Printf("\t%s\n", bot)
	}
}

type botAddRemoveRequest struct {
	BotID string `json:"botId"`
}

func startBot(url string, botID string) {
	body, err := json.Marshal(botAddRemoveRequest{
		BotID: botID,
	})
	if err != nil {
		fmt.Printf("Error starting bot with ID %s: Preparation of requestBody failed with error %s", botID, err)
		return
	}

	data, code, err := realAPICall(url, "/bots", "POST", string(body))
	if err != nil {
		fmt.Printf("Error starting bot with ID %s: %s", botID, err)
		return
	}

	if code == 400 {
		fmt.Printf("Error starting bot with ID %s: Bad Request %s", botID, data)
		return
	} else if code != 200 {
		fmt.Printf("Error starting bot with ID %s: Unknown StatusCode: %d", botID, code)
		return
	}

	fmt.Printf("Bot with ID %s started @ %s.\n", botID, url)
}

func stopBot(url string, botID string) {
	data, code, err := realAPICall(url, "/bots/"+botID, "DELETE", "")
	if err != nil {
		fmt.Printf("Error stopping bot with ID %s: %s", botID, err)
		return
	}

	if code == 404 {
		fmt.Printf("Error stopping bot with ID %s: Does not exist %s", botID, data)
		return
	} else if code == 400 {
		fmt.Printf("Error stopping bot with ID %s: Bad Request %s", botID, data)
		return
	} else if code != 200 {
		fmt.Printf("Error stopping bot with ID %s: Unknown StatusCode: %d", botID, code)
		return
	}

	fmt.Printf("Bot with ID %s stopped @ %s.\n", botID, url)
}

func main() {

	fmt.Printf("BotterControl Version %s\n\n", version)

	var (
		url          = flag.String("u", "", "URL to the Botter API")
		command      = flag.String("c", "", "Command to send to the Botter")
		argument     = flag.String("a", "", "Argument to send to the Botter")
		v            = flag.Bool("v", false, "prints current version and exits")
		listCommands = flag.Bool("list", false, "Lists all known commands")
	)

	flag.Parse()

	if *v {
		fmt.Printf("Version %s\n", version)
		os.Exit(0)
	}

	if *listCommands {
		fmt.Printf("Commands (case insensitive):\n")
		fmt.Printf("GetBots: \nReturns all running bots, argument: NONE\n")
		fmt.Printf("StartBot: \nStarts a bot, argument: botID\n")
		fmt.Printf("StopBot: \nStops a bot, argument: botID\n")
		os.Exit(0)
	}

	if len(*url) == 0 {
		fmt.Printf("Must specify an URL!\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if len(*command) == 0 {
		fmt.Printf("Must specify a command!\n\n")
		flag.Usage()
		os.Exit(1)
	}

	lowerCaseCommand := strings.ToLower(*command)

	cleanURL := strings.TrimRight(*url, "/")

	switch lowerCaseCommand {
	case "getbots":
		getBots(cleanURL)
	case "startbot":
		if len(*argument) == 0 {
			fmt.Printf("Must specify a bot id as argument")
			os.Exit(1)
		}
		startBot(cleanURL, *argument)
	case "stopbot":
		if len(*argument) == 0 {
			fmt.Printf("Must specify a bot id as argument")
			os.Exit(1)
		}
		stopBot(cleanURL, *argument)
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		os.Exit(1)
	}

	os.Exit(0)
}
