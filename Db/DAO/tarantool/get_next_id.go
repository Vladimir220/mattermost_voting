package tarantool

import (
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

func (t TarantoolDAO) getNextId(funcName string) (id uint, err error) {
	var res []uint
	err = t.conn.Do(
		tarantool.NewCallRequest(funcName).Args([]any{}),
	).GetTyped(&res)

	if err != nil {
		err = fmt.Errorf("ошибка получения следующего индекса: %v", err)
		return
	}
	if len(res) == 0 {
		err = fmt.Errorf("получено пустое значение вместо следующего индекса")
		return
	}

	id = res[0]

	return
}

func (t TarantoolDAO) getNextVotingId() (id uint, err error) {
	id, err = t.getNextId(votingNextIndexFuncName)

	return
}

func (t TarantoolDAO) getNextOptionsId() (id uint, err error) {
	id, err = t.getNextId(optionsNextIndexFuncName)

	return
}
