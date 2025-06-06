package tarantool

import (
	"github.com/tarantool/go-tarantool/v2"
)

const (
	votingsSpace = "votings"
	optionsSpace = "options"
	votedSpace   = "voted"
	optionsIndex = "options_id_idx"
	votingsIndex = "primary"
	votedIndex   = "voted_id_idx"

	votingNextIndexFuncName  = "get_next_voting_id"
	optionsNextIndexFuncName = "get_next_options_id"
)

type TarantoolDAO struct {
	conn *tarantool.Connection
}

func CreateTarantoolDAO(address, user, password string) (dao TarantoolDAO, err error) {
	dao = TarantoolDAO{}
	err = dao.init(address, user, password)

	return
}
