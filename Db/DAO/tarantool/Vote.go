package tarantool

import (
	"errors"
	"fmt"
	"voting_bot/Models"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) Vote(votingID, optionID uint) (err error) {
	voting, err := t.readVoting(votingID)
	if err != nil {
		return
	}
	if !voting.IsActive {
		err = errors.New("голосование уже завершено")
		return
	}

	options, err := t.readOptions(voting.OptionsId)
	if err != nil {
		return
	}

	var option Models.Option
	var isFind bool = false
	for _, option = range options {
		if option.OptionId == optionID {
			isFind = true
			break
		}
	}
	if !isFind {
		err = errors.New("такого варианта не существует")
		return
	}

	_, err = t.conn.Do(
		tarantool.NewUpdateRequest(optionsSpace).
			Index("options_id_idx").
			Key([]any{voting.OptionsId, optionID}).
			Operations(tarantool.NewOperations().Assign(3, option.Count+1)),
	).Get()
	if err != nil {
		err = fmt.Errorf("ошибка добавления голоса: %v", err)
		return
	}

	return
}
