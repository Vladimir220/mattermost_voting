package Network

import (
	"log"
	"time"
	"voting_bot/Models"

	"github.com/mattermost/mattermost/server/public/model"
)

func (n Network) BotListening(loginInfo Models.LoginInfo, botUser *model.User, client Models.Client, errLog *log.Logger) {
	stopSignal := make(chan struct{})

	go func() {
		for {
			select {
			case <-stopSignal:
				return
			default:
				select {
				case event := <-client.Ws.EventChannel:
					n.Handler.HandleEvent(event, client.Http, botUser, errLog)
				default:
					time.Sleep(300 * time.Millisecond)
				}

				if client.Ws.ListenError != nil {
					errLog.Println("Ошибка WebSocket:", client.Ws.ListenError)
					var err error
					client.Ws, err = n.ConnectionWS(loginInfo, errLog)
					if err != nil {
						errLog.Println("Не удалось переподключиться к WebSocket:", err)
					}
					client.Ws.ListenError = nil
				}
			}
		}
	}()

	defer client.Ws.Close()
	defer close(stopSignal)
	client.Ws.Listen()
	select {}
}
