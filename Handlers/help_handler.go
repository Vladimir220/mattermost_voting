package Handlers

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
)

func helpHandler(client *model.Client4, post *model.Post, user *model.User) (err error) {
	message := fmt.Sprintf(
		`@%s Доступные команды:
/create <вопрос>+++<вариант1>+++<вариант2> ...
/vote <Id голосования> <Id варианта>
/result <Id голосования>
/stop <Id голосования>
/delete <Id голосования>
Важно: обратите внимание на разделитель "+++" между параметрами команды создания голосования!
`, user.Username)

	err = sendMessage(client, message, post.ChannelId, post.RootId)
	if err != nil {
		return
	}

	return
}
