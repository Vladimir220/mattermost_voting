package tarantool

import (
	"errors"
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) Vote(votingID, optionID uint, votedId string) (err error) {
	voting, err := t.readVoting(votingID)
	if err != nil {
		return
	}
	if !voting.IsActive {
		err = errors.New("голосование уже завершено")
		return
	}

	_, err = t.conn.Do(tarantool.NewInsertRequest(votedSpace).Tuple([]any{votingID, votedId})).Get()
	if err != nil {
		err = errors.New("уже проголосовал")
		return
	}

	option, err := t.readOption(voting.OptionsId, optionID)
	if err != nil {
		return
	}

	_, err = t.conn.Do(
		tarantool.NewUpdateRequest(optionsSpace).
			Index(optionsIndex).
			Key([]any{voting.OptionsId, optionID}).
			Operations(tarantool.NewOperations().Assign(3, option.Count+1)),
	).Get()
	if err != nil {
		err = fmt.Errorf("ошибка добавления голоса: %v", err)
		_, err2 := t.conn.Do(
			tarantool.NewDeleteRequest(votedSpace).
				Index(votedIndex).
				Key([]any{votingID, votedId}),
		).Get()
		if err2 != nil {
			err2 = fmt.Errorf("ошибка отката записи о проголосовавшем: %v", err2)
			err = fmt.Errorf("%v\n%v", err, err2)
			return
		}
		return
	}

	return
}
