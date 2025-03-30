package Handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

func isMessageForBot(message, botUsername string) bool {
	mention := fmt.Sprintf("@%s", botUsername)
	return strings.Contains(message, mention)
}

func sendMessage(client *model.Client4, message, channelId, rootId string) (err error) {

	post := &model.Post{
		ChannelId: channelId,
		Message:   message,
		RootId:    rootId,
	}

	_, _, err = client.CreatePost(context.Background(), post)
	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}
	return
}
