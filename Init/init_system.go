package Init

import (
	"log"
	"os"
	"strings"
	"voting_bot/Models"

	"github.com/joho/godotenv"
)

// Читает файлы env, открывает файлы log.
func InitSystem() (loginInfo Models.LoginInfo, errLog *log.Logger) {

	err := godotenv.Load("./ENV/.env")
	if err != nil {
		panic(err.Error())
	}

	errLogFile, err := os.OpenFile(os.Getenv("ERR_LOG_FILE_PATH"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	errLog = log.New(errLogFile, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)

	loginInfo.Url = os.Getenv("MATTERMOST_URL")
	loginInfo.BotToken = os.Getenv("BOT_TOKEN")
	loginInfo.BotName = os.Getenv("BOT_NAME")

	isEmpty := false
	emptyVars := make([]string, 0, 3)
	if loginInfo.Url == "" {
		emptyVars = append(emptyVars, "MATTERMOST_URL")
		isEmpty = true
	}

	if loginInfo.BotToken == "" {
		emptyVars = append(emptyVars, "BOT_TOKEN")
		isEmpty = true
	}
	if loginInfo.BotName == "" {
		emptyVars = append(emptyVars, "BOT_NAME")
		isEmpty = true
	}

	if isEmpty {
		errLog.Panicln("Проверьте переменные среды:", strings.Join(emptyVars, ", "))
	}

	return
}
