package Handlers

import (
	"fmt"
	"strings"
	"voting_bot/Db/DAO"

	"github.com/mattermost/mattermost/server/public/model"
)

func createHandler(input []string, client *model.Client4, post *model.Post, user *model.User, dao DAO.DAO) (err error) {
	params, err := checkInputParameters("/create", 2, input, client, post)
	if err != nil {
		return
	}
	cleanParams(&params)

	question := params[0]

	votingID, resOptions, err := dao.CreateVoting(post.UserId, question, params[1:])
	if err != nil {
		message := fmt.Sprintf("@%s Внутренняя ошибка команды /create.", user.Username)
		err = sendMessage(client, message, post.ChannelId, post.RootId)
		if err != nil {
			return
		}
		return
	}

	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("@%s Id голосования: %d\nВарианты ответа:\n", user.Username, votingID))
	for _, option := range resOptions {
		message.WriteString(fmt.Sprintf("(ID: %d) %s -- %d голосов\n", option.OptionId, option.Text, option.Count))
	}

	err = sendMessage(client, message.String(), post.ChannelId, post.RootId)
	if err != nil {
		return
	}

	return
}
