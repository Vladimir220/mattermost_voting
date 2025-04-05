package Handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

func isMessageForBot(message, botUsername string) bool {
	mention := fmt.Sprintf("@%s ", botUsername)
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
		err = fmt.Errorf("не удалось отправить сообщение: %v", err)
	}
	return
}

func checkInputParameters(handlerName, sep string, minCountParams int, input []string, client *model.Client4, post *model.Post) (params []string, err error) {
	if len(input) == 0 {
		message := fmt.Sprintf("Кажется вы забыли про параметры команды %s. Для дополнительной информации используйте команду /help.", handlerName)
		err = sendMessage(client, message, post.ChannelId, post.RootId)
		if err != nil {
			return
		}
		return
	}

	params = strings.Split(input[1], sep)

	if len(params) < minCountParams {
		message := fmt.Sprintf("Вы ввели не все требуемые параметры для команды %s. Для дополнительной информации используйте команду /help.", handlerName)
		err = sendMessage(client, message, post.ChannelId, post.RootId)
		if err != nil {
			return
		}
		return
	}
	return
}

func cleanParams(params *[]string) {
	for i := range *params {
		(*params)[i] = strings.TrimSpace((*params)[i])
	}
}
