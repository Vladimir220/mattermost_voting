package tarantool

import (
	"errors"
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) DeleteVoting(votingID uint, initiatorId string) (err error) {
	voting, err := t.readVoting(votingID)
	if err != nil {
		return
	}
	if voting.CreatorId != initiatorId {
		err = errors.New("403")
		return
	}
	options, err := t.readAllOptions(voting.OptionsId)
	if err != nil {
		return
	}
	voted, err := t.readAllVoted(votingID)
	if err != nil && err.Error() != "empty" {
		return
	}

	for _, option := range options {
		_, err = t.conn.Do(
			tarantool.NewDeleteRequest(optionsSpace).
				Index(optionsIndex).
				Key([]any{option.Id, option.OptionId}),
		).Get()
		if err != nil {
			err = fmt.Errorf("ошибка удаления вариантов ответа: %v", err)
			return
		}
	}

	for _, v := range voted {
		_, err = t.conn.Do(
			tarantool.NewDeleteRequest(votedSpace).
				Index(votedIndex).
				Key([]any{v.VotingId, v.Id}),
		).Get()
		if err != nil {
			err = fmt.Errorf("ошибка удаления записей о проголосовавших: %v", err)
			return
		}
	}

	_, err = t.conn.Do(
		tarantool.NewDeleteRequest(votingsSpace).
			Key([]any{votingID}),
	).Get()
	if err != nil {
		err = fmt.Errorf("ошибка удаления голосования: %v", err)
		return
	}

	return
}
