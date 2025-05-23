package DAO

import "voting_bot/Models"

type DAO interface {
	CreateVoting(creatorID string, question string, options []string) (votingID uint, resOptions []Models.Option, err error)
	ReadResults(votingID uint) (question string, resOptions []Models.Option, err error)
	Vote(votingID, optionID uint, votedId string) (err error)
	DeleteVoting(votingID uint, initiatorId string) (err error)
	StopVoting(votingID uint, initiatorId string) (err error)
	Close()
}
