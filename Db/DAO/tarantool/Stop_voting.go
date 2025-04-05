package tarantool

import (
	"errors"
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) StopVoting(votingID uint, initiatorId string) (err error) {
	voting, err := t.readVoting(votingID)
	if err != nil {
		return
	}
	if voting.CreatorId != initiatorId {
		err = errors.New("403")
		return
	}

	_, err = t.conn.Do(
		tarantool.NewUpdateRequest(votingsSpace).
			Key([]any{votingID}).
			Operations(tarantool.NewOperations().Assign(4, false)),
	).Get()
	if err != nil {
		err = fmt.Errorf("ошибка деактивации голосования: %v", err)
		return
	}

	return

}
