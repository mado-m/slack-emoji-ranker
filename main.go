package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	println("begin.")

	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(slackApiToken)

	channels, err := getAllChannels(api)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Get [%d] channels.\n", len(channels))

	oldest := time.Now().AddDate(0, 0, -30).Unix()
	for _, channel := range channels {
		fmt.Printf("channelId:[%s], channelName:[%s] .\n", channel.ID, channel.Name)
		messages, err := getAllMessages(api, channel.ID, oldest)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Get [%d] messages.\n", len(messages))
	}

	println("end.")
}

func getAllChannels(api *slack.Client) ([]slack.Channel, error) {
	channels, _, err := api.GetConversations(&slack.GetConversationsParameters{
		ExcludeArchived: true,
		Limit:           10,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	return channels, nil
}

func getAllMessages(api *slack.Client, channelId string, oldest int64) ([]slack.Message, error) {
	res, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: channelId,
		Limit:     1000,
		Oldest:    strconv.FormatInt(oldest, 10),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	return res.Messages, nil
}
