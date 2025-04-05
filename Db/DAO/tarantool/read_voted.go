package tarantool

import (
	"errors"
	"fmt"
	"voting_bot/Models"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) readVoted(votingId uint, votedId string) (voted Models.Voted, err error) {
	var buf []Models.Voted
	err = t.conn.Do(tarantool.NewSelectRequest(votedSpace).
		Index(votedIndex).
		Iterator(tarantool.IterEq).
		Key([]any{votingId, votedId}),
	).GetTyped(&buf)

	if err != nil {
		err = fmt.Errorf("ошибка получения проголосовавших: %v", err)
		return
	}
	if len(buf) == 0 {
		err = errors.New("ошибка получения проголосовавших: получен пустой список")
		return
	}

	voted = buf[0]

	return
}

func (t TarantoolDAO) readAllVoted(votingId uint) (voted []Models.Voted, err error) {
	err = t.conn.Do(tarantool.NewSelectRequest(votedSpace).
		Index(votedIndex).
		Iterator(tarantool.IterEq).
		Key([]any{votingId}),
	).GetTyped(&voted)

	if err != nil {
		err = fmt.Errorf("ошибка получения проголосовавших: %v", err)
		return
	}
	if len(voted) == 0 {
		err = errors.New("empty")
		return
	}

	return
}
