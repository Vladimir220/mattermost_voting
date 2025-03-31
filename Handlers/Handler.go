package Handlers

import (
	"context"
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

			if len(parts) < 2 {
				message := "Кажется вы забыли про параметры команды."
				err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}
				return
			}

			params := strings.Split(parts[1], "\" \"")
			user, _, err := clientHTTP.GetUser(context.TODO(), post.UserId, "")
			if err != nil {
				errLog.Printf("Не удалось получить информацию о пользователе %s: %v\n", post.UserId, err)
				return
			}

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

				votingID, resOptions, err := h.Dao.CreateVoting(post.UserId, question, params[1:])
				if err != nil {
					errLog.Println(err)
					message := "Внутренняя ошибка команды /create"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}
				message := strings.Builder{}
				message.WriteString(fmt.Sprintf("@%s Id голосования: %d\nВарианты ответа:\n", user.Username, votingID))
				for _, option := range resOptions {
					message.WriteString(fmt.Sprintf("(ID: %d) %s -- %d голосов\n", option.OptionId, option.Text, option.Count))
				}

				err = sendMessage(clientHTTP, message.String(), post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}
			case "/vote":
				params = strings.Split(params[0], " ")
				if len(params) < 2 {
					message := "Вы не ввели Id голосования и/или Id варианта ответа для команды /vote"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}
				votingID, err := strconv.ParseUint(params[0], 10, 64)
				if err != nil {
					errLog.Println(err)
					return
				}
				optionID, err := strconv.ParseUint(params[1], 10, 64)
				if err != nil {
					errLog.Println(err)
					return
				}

				err = h.Dao.Vote(uint(votingID), uint(optionID))
				if err != nil {
					errLog.Println(err)
					message := "Голосования с таким ID не существует, либо оно закрыто владельцем"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}

				err = sendMessage(clientHTTP, fmt.Sprintf("@%s Принято!", user.Username), post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}

			case "/result":
				if len(params) < 1 {
					message := "Вы не ввели Id голосования для команды /result"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}
				votingID, err := strconv.ParseUint(params[0], 10, 64)
				if err != nil {
					errLog.Println(err)
					return
				}

				question, resOptions, err := h.Dao.ReadResults(uint(votingID))
				if err != nil {
					errLog.Println(err)
					message := "Голосования с таким ID не существует"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}

				message := strings.Builder{}
				message.WriteString(fmt.Sprintf("@%s Id голосования: %d\nВопрос: %s\nВарианты ответа:\n", user.Username, votingID, question))
				for _, option := range resOptions {
					message.WriteString(fmt.Sprintf("(ID: %d) %s -- %d голосов\n", option.OptionId, option.Text, option.Count))
				}

				err = sendMessage(clientHTTP, message.String(), post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}

			case "/stop":
				if len(params) < 1 {
					message := "Вы не ввели Id голосования для команды /stop"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}
				votingID, err := strconv.ParseUint(params[0], 10, 64)
				if err != nil {
					errLog.Println(err)
					return
				}

				err = h.Dao.StopVoting(uint(votingID))
				if err != nil {
					errLog.Println(err)
					message := "Внутреняя ошибка команды /stop"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}
				err = sendMessage(clientHTTP, fmt.Sprintf("@%s Голосование (Id:%d) закрыто!", user.Username, votingID), post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}

			case "/delete":
				if len(params) < 1 {
					message := "Вы не ввели Id голосования для команды /delete"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}
				votingID, err := strconv.ParseUint(params[0], 10, 64)
				if err != nil {
					errLog.Println(err)
					return
				}

				err = h.Dao.DeleteVoting(uint(votingID))
				if err != nil {
					errLog.Println(err)
					message := "Внутреняя ошибка команды /delete"
					err := sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
					if err != nil {
						errLog.Println(err)
						return
					}
					return
				}
				err = sendMessage(clientHTTP, fmt.Sprintf("@%s Голосование (Id:%d) удалено!", user.Username, votingID), post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}

			default:
				message := fmt.Sprintf("@%s Такой команды не существует. Доступные команды:\n/create\n/vote\n/result\n/stop\n/delete\n", user.Username)
				err = sendMessage(clientHTTP, message, post.ChannelId, post.RootId)
				if err != nil {
					errLog.Println(err)
					return
				}
			}

		}
	}
}
