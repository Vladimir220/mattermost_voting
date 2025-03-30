package tarantool

import (
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) StopVoting(votingID uint) (err error) {
	_, err = t.conn.Do(
		tarantool.NewUpdateRequest(votingsSpace).
			Index("primary").
			Key([]any{votingID}).
			Operations(tarantool.NewOperations().Assign(4, false)),
	).Get()
	if err != nil {
		err = fmt.Errorf("ошибка деактивации голосования: %v", err)
		return
	}

	return

}
