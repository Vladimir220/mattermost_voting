package Handlers

import (
	"fmt"
	"strconv"
	"strings"
	"voting_bot/Db/DAO"

	"github.com/mattermost/mattermost/server/public/model"
)

func resultsHandler(input []string, client *model.Client4, post *model.Post, user *model.User, dao DAO.DAO) (err error) {
	params, err := checkInputParameters("/results", 1, input, client, post)
	if err != nil {
		return
	}
	cleanParams(&params)

	votingID, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return
	}

	question, resOptions, err := dao.ReadResults(uint(votingID))
	if err != nil {
		message := fmt.Sprintf("@%s Голосования с таким ID не существует", user.Username)
		err = sendMessage(client, message, post.ChannelId, post.RootId)
		if err != nil {
			return
		}
		return
	}

	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("@%s Id голосования: %d\nВопрос: %s\nВарианты ответа:\n", user.Username, votingID, question))
	for _, option := range resOptions {
		message.WriteString(fmt.Sprintf("(ID: %d) %s -- %d голосов\n", option.OptionId, option.Text, option.Count))
	}

	err = sendMessage(client, message.String(), post.ChannelId, post.RootId)
	if err != nil {
		return
	}

	return
}
