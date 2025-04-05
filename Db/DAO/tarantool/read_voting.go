package tarantool

import (
	"errors"
	"fmt"
	"voting_bot/Models"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) readVoting(votingID uint) (voting Models.Voting, err error) {
	inputVoting := []Models.Voting{}
	err = t.conn.Do(tarantool.NewSelectRequest(votingsSpace).
		Limit(1).
		Iterator(tarantool.IterEq).
		Key([]any{uint(votingID)}),
	).GetTyped(&inputVoting)

	if err != nil {
		err = fmt.Errorf("ошибка получения данных о голосовании: %v", err)
		return
	}
	if len(inputVoting) == 0 {
		err = errors.New("empty")
		return
	}

	voting = inputVoting[0]
	return
}
