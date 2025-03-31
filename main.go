package main

import (
	"voting_bot/Db/DAO/tarantool"
	"voting_bot/Handlers"
	"voting_bot/Init"
	"voting_bot/Models"
	"voting_bot/Network"
)

func main() {
	loginInfo, errLog := Init.InitSystem()

	dao, err := tarantool.CreateTarantoolDAO(loginInfo.TarantoolUrl, loginInfo.TarantoolLogin, loginInfo.TarantoolPassword)
	if err != nil {
		errLog.Panicln(err)
	}
	defer dao.Close()

	handler := Handlers.Handler{Dao: dao}
	network := Network.Network{Handler: handler}
	botUser, clientHTTP, err := network.ConnectionHTTP(loginInfo, errLog)
	if err != nil {
		errLog.Panicln("Не удалось получить информацию о боте:", err)
	}

	clientWS, err := network.ConnectionWS(loginInfo, errLog)
	if err != nil {
		errLog.Panicln("Не удалось подключиться к WebSocket:", err)
	}

	client := Models.Client{Http: clientHTTP, Ws: clientWS}

	network.BotListening(loginInfo, botUser, client, errLog)

}
