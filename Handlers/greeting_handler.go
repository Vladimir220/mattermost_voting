package Handlers

/*
func RespondWithGreeting(client *model.Client4, userID, rootId string, channelID string) {
	user, _, err := client.GetUser(context.Background(), userID, "")
	if err != nil {
		log.Printf("Не удалось получить информацию о пользователе %s: %v", userID, err)
		return
	}

	message := fmt.Sprintf("Привет, @%s! Я бот!", user.Username) // Или user.Username, user.Nickname, в зависимости от ваших нужд.

	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
		RootId:    rootId,
	}

	_, _, err = client.CreatePost(context.Background(), post)
	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}
}
*/
