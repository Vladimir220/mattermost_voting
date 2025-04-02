package Handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
				message := "Введите команду и параметры для работы с ботом голосования. Для уточнения существующих команд введите команду /help."
				err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}
				return
			} else {
				parts = strings.SplitN(parts[1], " ", 2)
			}
			comm := strings.TrimSpace(parts[0])

			user, _, err := clientHTTP.GetUser(context.TODO(), post.UserId, "")
			if err != nil {
				errLog.Printf("Не удалось получить информацию о пользователе %s: %v\n", post.UserId, err)
				return
			}

			switch comm {
			case "/help":
				err = helpHandler(clientHTTP, post, user)
				if err != nil {
					errLog.Println(err)
					return
				}

			case "/create":
				err = createHandler(parts, clientHTTP, post, user, h.Dao)
				if err != nil {
					errLog.Println(err)
					return
				}

			case "/vote":
				err = voteHandler(parts, clientHTTP, post, user, h.Dao)
				if err != nil {
					errLog.Println(err)
					return
				}

			case "/result":
				err = resultsHandler(parts, clientHTTP, post, user, h.Dao)
				if err != nil {
					errLog.Println(err)
					return
				}

			case "/stop":
				err = stopHandler(parts, clientHTTP, post, user, h.Dao)
				if err != nil {
					errLog.Println(err)
					return
				}

			case "/delete":
				err = deleteHandler(parts, clientHTTP, post, user, h.Dao)
				if err != nil {
					errLog.Println(err)
					return
				}

			default:
				message := fmt.Sprintf("@%s Такой команды не существует. Для уточнения существующих команд введите команду /help.", user.Username)
				err = sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}
			}

		}
	}
}
