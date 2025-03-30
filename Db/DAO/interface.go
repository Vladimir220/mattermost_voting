package DAO

import "voting_bot/Models"

type DAO interface {
	CreateVoting(creatorID uint, question string, options []string) (votingID uint, resOptions []Models.Option, err error)
	ReadResults(votingID uint) (question string, resOptions []Models.Option, err error)
	Vote(votingID, optionID uint) (err error)
	DeleteVoting(votingID uint) (err error)
	StopVoting(votingID uint) (err error)
	Close()
}
