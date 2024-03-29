package pof

import "errors"

// consensus errors
var (
	ErrInsufficientCandidateCount    = errors.New("insufficient candidate count")
	ErrExceedCandidateCount          = errors.New("exceed candidate count")
	ErrInvalidMaxBlocksPerFormulator = errors.New("invalid max blocks per formulator")
	ErrInvalidObserverKey            = errors.New("invalid observer key")
	ErrInvalidTopAddress             = errors.New("invalid top address")
	ErrInvalidTopSignature           = errors.New("invalid top signature")
	ErrInvalidSignatureCount         = errors.New("invalid signature count")
	ErrInvalidPhase                  = errors.New("invalid phase")
	ErrExistAddress                  = errors.New("exist address")
	ErrFoundForkedBlockGen           = errors.New("found forked block gen")
	ErrInvalidVote                   = errors.New("invalid vote")
	ErrInvalidRoundState             = errors.New("invalid round state")
	ErrInvalidRequest                = errors.New("invalid request")
	ErrAlreadyVoted                  = errors.New("already voted")
	ErrNotExistObserverPeer          = errors.New("not exist observer peer")
	ErrNotExistFormulatorPeer        = errors.New("not exist formulator peer")
)
