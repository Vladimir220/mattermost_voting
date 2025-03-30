package tarantool

import (
	"errors"
	"fmt"
	"voting_bot/Models"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) readOptions(optionsId uint) (options []Models.Option, err error) {
	err = t.conn.Do(tarantool.NewSelectRequest(optionsSpace).
		Index("options_id_idx").
		Iterator(tarantool.IterEq).
		Key([]any{uint(optionsId)}),
	).GetTyped(&options)

	if err != nil {
		err = fmt.Errorf("ошибка получения результатов голосования: %v", err)
		return
	}
	if len(options) == 0 {
		err = errors.New("ошибка получения результатов голосования: получен пустой список")
		return
	}

	return
}
