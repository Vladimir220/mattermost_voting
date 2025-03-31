package Models

import "github.com/mattermost/mattermost/server/public/model"

type LoginInfo struct {
	Url               string
	BotToken          string
	BotName           string
	TarantoolLogin    string
	TarantoolPassword string
	TarantoolUrl      string
}

type Client struct {
	Http *model.Client4
	Ws   *model.WebSocketClient
}

type Voting struct {
	Id        uint   `msgpack:"id"`
	CreatorId string `msgpack:"creator_id"`
	Question  string `msgpack:"question"`
	OptionsId uint   `msgpack:"options_id"`
	IsActive  bool   `msgpack:"is_active"`
}

type Option struct {
	Id       uint   `msgpack:"options_id"`
	OptionId uint   `msgpack:"option_id"`
	Text     string `msgpack:"text"`
	Count    int    `msgpack:"count"`
}
