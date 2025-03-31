package tarantool

import (
	"errors"
	"fmt"
	"voting_bot/Models"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) CreateVoting(creatorID string, question string, options []string) (votingID uint, resOptions []Models.Option, err error) {

	// Получаем индекс вариантов ответа
	optionsID, err := t.getNextOptionsId()
	if err != nil {
		return
	}

	// Создаём голосование
	voting := []any{
		nil,
		creatorID,
		question,
		optionsID,
		true,
	}

	res := []Models.Voting{}
	err = t.conn.Do(tarantool.NewInsertRequest(votingsSpace).Tuple(voting)).GetTyped(&res)
	if err != nil {
		err = fmt.Errorf("ошибка создания нового голосования: %v", err)
		return
	}
	if len(res) == 0 {
		err = errors.New("ошибка создания нового голосования: получен пустой список")
		return
	}
	votingID = res[0].Id

	// Создаём варианты ответа
	resOptions = make([]Models.Option, 0, len(options))
	for i, optionText := range options {
		option := []any{
			optionsID,
			uint(i),
			optionText,
			0,
		}

		_, err = t.conn.Do(tarantool.NewInsertRequest(optionsSpace).Tuple(option)).Get()
		if err != nil {
			err = fmt.Errorf("ошибка создания вариантов ответа: %v", err)
			return
		}

		resOp := Models.Option{
			Id:       optionsID,
			OptionId: uint(i),
			Text:     optionText,
			Count:    0,
		}
		resOptions = append(resOptions, resOp)
	}

	return
}
