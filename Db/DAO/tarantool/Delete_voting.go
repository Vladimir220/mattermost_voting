package tarantool

import (
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) DeleteVoting(votingID uint) (err error) {
	voting, err := t.readVoting(votingID)
	if err != nil {
		return
	}
	options, err := t.readOptions(voting.OptionsId)
	if err != nil {
		return
	}

	for _, option := range options {
		_, err = t.conn.Do(
			tarantool.NewDeleteRequest(optionsSpace).
				Index("options_id_idx").
				Key([]any{option.Id, option.OptionId}),
		).Get()
		if err != nil {
			err = fmt.Errorf("ошибка удаления вариантов ответа: %v", err)
			return
		}
	}

	_, err = t.conn.Do(
		tarantool.NewDeleteRequest(votingsSpace).
			Key([]any{votingID}),
	).Get()
	if err != nil {
		err = fmt.Errorf("ошибка голосования: %v", err)
		return
	}

	return
}
