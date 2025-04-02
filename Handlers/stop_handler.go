package Handlers

import (
	"fmt"
	"strconv"
	"voting_bot/Db/DAO"

	"github.com/mattermost/mattermost/server/public/model"
)

func stopHandler(input []string, client *model.Client4, post *model.Post, user *model.User, dao DAO.DAO) (err error) {
	params, err := checkInputParameters("/stop", 1, input, client, post)
	if err != nil {
		return
	}
	cleanParams(&params)

	votingID, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return
	}

	err = dao.StopVoting(uint(votingID))
	if err != nil {
		message := fmt.Sprintf("@%s Внутреняя ошибка команды /stop.", user.Username)
		err = sendMessage(client, message, post.ChannelId, post.RootId)
		if err != nil {
			return
		}
		return
	}
	err = sendMessage(client, fmt.Sprintf("@%s Голосование (Id:%d) закрыто!", user.Username, votingID), post.ChannelId, post.RootId)
	if err != nil {
		return
	}

	return
}
