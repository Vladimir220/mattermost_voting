package Handlers

import (
	"fmt"
	"strconv"
	"voting_bot/Db/DAO"

	"github.com/mattermost/mattermost/server/public/model"
)

func voteHandler(input []string, client *model.Client4, post *model.Post, user *model.User, dao DAO.DAO) (err error) {
	params, err := checkInputParameters("/vote", " ", 2, input, client, post)
	if err != nil {
		return
	}
	cleanParams(&params)

	votingID, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return
	}
	optionID, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return
	}

	err = dao.Vote(uint(votingID), uint(optionID), user.Id)
	if err != nil {
		var message string
		if err.Error() == "уже проголосовал" {
			message = fmt.Sprintf("@%s Вы уже проголосовали.", user.Username)
		} else {
			message = fmt.Sprintf("@%s Голосования с таким ID не существует, либо оно закрыто владельцем.", user.Username)
		}
		err = sendMessage(client, message, post.ChannelId, post.RootId)
		if err != nil {
			return
		}
		return
	}

	err = sendMessage(client, fmt.Sprintf("@%s Принято!", user.Username), post.ChannelId, post.RootId)
	if err != nil {
		return
	}

	return
}
