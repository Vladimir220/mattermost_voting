package Handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"voting_bot/Db/DAO"

	"github.com/mattermost/mattermost/server/public/model"
)

type Handler struct {
	Dao DAO.DAO
}

func (h Handler) HandleEvent(event *model.WebSocketEvent, clientHTTP *model.Client4, botUser *model.User, errLog *log.Logger) {
	if event == nil {
		return
	}

	if event.EventType() == model.WebsocketEventPosted {
		postJSON, ok := event.GetData()["post"].(string)
		if !ok {
			errLog.Println("Не удалось получить данные сообщения из события.")
			return
		}

		var post *model.Post
		if err := json.Unmarshal([]byte(postJSON), &post); err != nil {
			errLog.Println("Не удалось unmarshal сообщение:", err)
			return
		}

		if post.UserId == botUser.Id {
			return
		}

		if isMessageForBot(post.Message, botUser.Username) {
			parts := strings.SplitN(post.Message, " ", 2)

			if len(parts) == 1 {
				message := "Введите команду и параметры для работы с ботом голосования."
				err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}
				return
			} else {
				parts = strings.SplitN(parts[1], " ", 2)
			}
			comm := parts[0]

			params := strings.Split(parts[1], "\" \"")

			switch comm {
			case "/create":
				if len(params) < 2 {
					message := "Вы не ввели вопрос или варинтов ответа для команды /create"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}

				question := params[0]
				userId, err := strconv.ParseUint(post.UserId, 10, 64)
				if err != nil {
					errLog.Println(err)
					return
				}

				votingID, resOptions, err := h.Dao.CreateVoting(uint(userId), question, params[1:])
				if err != nil {
					errLog.Println(err)
					message := "Внутреняя ошибка команды /create"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}
				message := fmt.Sprintf("Id голосования: %d\nВарианты ответа: %v", votingID, resOptions)
				err = sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}
			case "/vote":
			case "/result":
			case "/strop":
			case "/delete":
			default:
			}

			//RespondWithGreeting(clientHTTP, post.UserId, post.RootId, post.ChannelId)
		}
	}
}
