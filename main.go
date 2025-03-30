/*package main


import (
	"fmt"
	"log"
	"net/http"
	"strings"



	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Plugin структура основного плагина
type Plugin struct {
	plugin.MattermostPlugin
}

// OnActivate вызывается при активации плагина
func (p *Plugin) OnActivate() error {
	log.Println("Activating bot plugin")
	return nil
}

// MessageHasBeenPosted вызывается при публикации нового сообщения
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	// Получаем информацию о пользователе, отправившем сообщение
	user, appErr := p.API.GetUser(post.UserId)
	if appErr != nil {
		log.Printf("Failed to get user: %s\n", appErr.Error())
		return
	}

	// Получаем информацию о боте
	botUser, appErr := p.API.GetUserByUsername("votint" os.Getenv("BOT_USERNAME"))
	if appErr != nil {
		log.Printf("Failed to get bot user: %s\n", appErr.Error())
		return
	}

	if post.UserId == botUser.Id {
		return
	}

	// Проверяем, что сообщение адресовано боту (можно добавить логику упоминания)
	if !p.isMessageForBot(post.Message, botUser.Username) {
		return
	}

	// Формируем ответ
	response := fmt.Sprintf("Привет, %s! Рад тебя видеть.", user.FirstName)

	// Создаем новый пост с ответом
	newPost := &model.Post{
		ChannelId: post.ChannelId,
		Message:   response,
		RootId:    post.RootId, // Отвечаем в треде, если это ответ на другое сообщение
		UserId:    botUser.Id,  // Указываем, что пост отправлен ботом
	}

	// Отправляем ответ
	_, appErr = p.API.CreatePost(newPost)
	if appErr != nil {
		log.Printf("Failed to create post: %s\n", appErr.Error())
	}
}

// isMessageForBot проверяет, адресовано ли сообщение боту (простой пример)
func (p *Plugin) isMessageForBot(message string, botUsername string) bool {
	//return true // Отвечаем на любое сообщение
	// Расширенный пример:
	return strings.Contains(strings.ToLower(message), "@"+botUsername)
}

// ServeHTTP позволяет плагину обслуживать HTTP-запросы (не используется в этом примере)
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// main Функция main, необходимая для плагина
func main() {
	plugin.ClientMain(&Plugin{})
}
*/

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

	dao, err := tarantool.CreateTarantoolDAO("localhost:3301", "user", "qwerty") // заменить!
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

	/*fmt.Println(dao, err)
	data, err2 := dao.GetNextOptionsId()
	fmt.Println(data, err2)*/

	/*votingID, _, err := dao.CreateVoting(13, "Ху?", []string{"я", "ты"})
	fmt.Println(votingID, err)

	res, ress, err := dao.ReadResults(votingID)
	fmt.Println(votingID, res, err, ress)

	err = dao.Vote(votingID, 1)
	fmt.Println(err)

	res, ress, err = dao.ReadResults(votingID)
	fmt.Println(votingID, res, err, ress)

	err = dao.StopVoting(votingID)
	fmt.Println(err)

	err = dao.Vote(votingID, 1)
	fmt.Println(err)

	res, ress, err = dao.ReadResults(votingID)
	fmt.Println(votingID, res, err, ress)

	err = dao.DeleteVoting(votingID)
	fmt.Println(err)

	res, ress, err = dao.ReadResults(votingID)
	fmt.Println(votingID, res, err, ress)

	//fmt.Printf("%T\n", res[0].CreatorId)*/
}
