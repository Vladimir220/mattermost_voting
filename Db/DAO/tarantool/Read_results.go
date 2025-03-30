package tarantool

import (
	"voting_bot/Models"
)

func (t TarantoolDAO) ReadResults(votingID uint) (question string, resOptions []Models.Option, err error) {
	// Читаем данные по указанному голосованию
	voting, err := t.readVoting(votingID)
	if err != nil {
		return
	}

	question = voting.Question

	// Читаем результаты голосования
	resOptions, err = t.readOptions(voting.OptionsId)
	if err != nil {
		return
	}

	return
}
