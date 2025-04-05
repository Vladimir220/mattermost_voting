package Handlers

import (
	"fmt"
	"strconv"
	"voting_bot/Db/DAO"

	"github.com/mattermost/mattermost/server/public/model"
)

func deleteHandler(input []string, client *model.Client4, post *model.Post, user *model.User, dao DAO.DAO) (err error) {
	params, err := checkInputParameters("/delete", " ", 1, input, client, post)
	if err != nil {
		return
	}
	cleanParams(&params)

	votingID, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return
	}

	err = dao.DeleteVoting(uint(votingID), user.Id)
	var message string
	if err != nil {
		if err.Error() == "403" {
			message = fmt.Sprintf("@%s У вас нет прав удалять это голосование!", user.Username)
		} else {
			message = fmt.Sprintf("@%s Внутреняя ошибка команды /delete.", user.Username)
		}
		err = sendMessage(client, message, post.ChannelId, post.RootId)
		if err != nil {
			return
		}
		return
	}
	err = sendMessage(client, fmt.Sprintf("@%s Голосование (Id:%d) удалено!", user.Username, votingID), post.ChannelId, post.RootId)
	if err != nil {
		return
	}

	return
}
