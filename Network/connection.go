package Network

import (
	"context"
	"log"
	"strings"
	"time"
	"voting_bot/Models"

	"github.com/mattermost/mattermost/server/public/model"
)

func (n Network) ConnectionHTTP(loginInfo Models.LoginInfo, errLog *log.Logger) (botUser *model.User, client *model.Client4, err error) {
	client = model.NewAPIv4Client(loginInfo.Url)
	client.SetToken(loginInfo.BotToken)

	ctx, ctxFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxFunc()
	botUser, _, err = client.GetUserByUsername(ctx, loginInfo.BotName, "")

	return
}

func (n Network) ConnectionWS(loginInfo Models.LoginInfo, errLog *log.Logger) (wsClient *model.WebSocketClient, err error) {
	wsURL := strings.Replace(loginInfo.Url, "http", "ws", 1)
	wsURL = strings.Replace(wsURL, "https", "wss", 1)

	wsClient, err = model.NewWebSocketClient4(wsURL, loginInfo.BotToken)

	return
}
