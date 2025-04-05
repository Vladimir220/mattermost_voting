package tarantool

import (
	"errors"
	"fmt"
	"voting_bot/Models"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) readAllOptions(optionsId uint) (options []Models.Option, err error) {
	err = t.conn.Do(tarantool.NewSelectRequest(optionsSpace).
		Index(optionsIndex).
		Iterator(tarantool.IterEq).
		Key([]any{optionsId}),
	).GetTyped(&options)

	if err != nil {
		err = fmt.Errorf("ошибка получения результатов голосования: %v", err)
		return
	}
	if len(options) == 0 {
		err = errors.New("empty")
		return
	}

	return
}

func (t TarantoolDAO) readOption(optionsId, optionId uint) (option Models.Option, err error) {
	var buf []Models.Option
	err = t.conn.Do(tarantool.NewSelectRequest(optionsSpace).
		Index(optionsIndex).
		Iterator(tarantool.IterEq).
		Key([]any{optionsId, optionId}),
	).GetTyped(&buf)

	if err != nil {
		err = fmt.Errorf("ошибка получения результата голосования: %v", err)
		return
	}
	if len(buf) == 0 {
		err = errors.New("empty")
		return
	}
	option = buf[0]

	return
}
