package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	println("begin.")

	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(slackApiToken)

	// 全チャンネル取得
	channels, err := getAllChannels(api)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Get [%d] channels.\n", len(channels))

	// 対象期間の全emojiを取得
	emojiMap := map[string]int{}
	oldest := time.Now().AddDate(0, 0, -30).Unix()
	for _, channel := range channels {
		fmt.Printf("channelId:[%s], channelName:[%s] .\n", channel.ID, channel.Name)
		messages, err := getAllMessages(api, channel.ID, oldest)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Get [%d] messages.\n", len(messages))

		for _, message := range messages {
			reactions := message.Reactions
			for _, rereaction := range reactions {
				emojiMap[rereaction.Name]++
			}
		}
	}

	// mapのvalue順にソート
	type kv struct {
		Key   string
		Value int
	}
	var emojiCounter []kv
	for k, v := range emojiMap {
		emojiCounter = append(emojiCounter, kv{k, v})
	}
	sort.Slice(emojiCounter, func(i, j int) bool {
		return emojiCounter[i].Value > emojiCounter[j].Value
	})
	for i, kv := range emojiCounter {
		// 上位だけ取得
		if i >= 30 {
			break
		}
		fmt.Printf("%s: %d\n", kv.Key, kv.Value)
	}

	println("end.")
}

func getAllChannels(api *slack.Client) ([]slack.Channel, error) {
	channels, _, err := api.GetConversations(&slack.GetConversationsParameters{
		ExcludeArchived: true,
		Limit:           1000, // 全チャンネルが取得できる件数にする
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	return channels, nil
}

func getAllMessages(api *slack.Client, channelId string, oldest int64) ([]slack.Message, error) {
	messages := []slack.Message{}
	cursor := ""
	for {
		res, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
			ChannelID: channelId,
			Limit:     1000,
			Oldest:    strconv.FormatInt(oldest, 10),
			Cursor:    cursor,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		if res.Ok {
			messages = append(messages, res.Messages...)
			if res.HasMore {
				cursor = res.ResponseMetaData.NextCursor
			} else {
				break
			}
		} else {
			return nil, fmt.Errorf("message res is error: %s", res.Error)
		}
	}
	return messages, nil
}
