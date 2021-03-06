package monitor

import (
	"context"

	"github.com/celo-org/kliento/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

type electionProcessor struct {
	ctx             context.Context
	logger          log.Logger
	electionAddress common.Address
	election        *contracts.Election
}

func NewElectionProcessor(ctx context.Context, logger log.Logger, electionAddress common.Address, election *contracts.Election) *electionProcessor {
	return &electionProcessor{
		ctx:             ctx,
		logger:          logger,
		electionAddress: electionAddress,
		election:        election,
	}
}

func (p electionProcessor) ObserveState(opts *bind.CallOpts) error {
	logger := p.logger.New("contract", "Election")

	// Election.getActiveVotes
	activeVotes, err := p.election.GetActiveVotes(opts)
	if err != nil {
		return err
	}

	logStateViewCall(logger, "method", "getActiveVotes", "activeVotes", activeVotes)

	// Election.getTotalVotes
	totalVotes, err := p.election.GetTotalVotes(opts)
	if err != nil {
		return err
	}

	logStateViewCall(logger, "method", "getTotalVotes", "totalVotes", totalVotes)

	// Election.getElectableValidators
	electableValidatorsMin, electableValidatorsMax, err := p.election.GetElectableValidators(opts)
	if err != nil {
		return err
	}

	logStateViewCall(logger, "method", "getElectableValidators", "electableValidatorsMin", electableValidatorsMin.Uint64())
	logStateViewCall(logger, "method", "getElectableValidators", "electableValidatorsMax", electableValidatorsMax.Uint64())

	return nil
}

func (p electionProcessor) HandleLog(eventLog *types.Log) {
	logger := p.logger.New("contract", "Election")
	if eventLog.Address == p.electionAddress {
		eventName, eventRaw, ok, err := p.election.TryParseLog(*eventLog)
		if err != nil {
			logger.Warn("Ignoring event: Error parsing election event", "err", err, "eventId", eventLog.Topics[0].Hex())
			return
		}
		if !ok {
			return
		}

		switch eventName {
		case "ValidatorGroupVoteCast":
			event := eventRaw.(*contracts.ElectionValidatorGroupVoteCast)
			logEventLog(logger, "eventName", eventName, "account", event.Account, "group", event.Group, "value", event.Value)
		case "ValidatorGroupVoteActivated":
			event := eventRaw.(*contracts.ElectionValidatorGroupVoteActivated)
			logEventLog(logger, "eventName", eventName, "account", event.Account, "group", event.Group, "value", event.Value, "voteUnits", event.Units)
		case "ValidatorGroupPendingVoteRevoked":
			event := eventRaw.(*contracts.ElectionValidatorGroupPendingVoteRevoked)
			logEventLog(logger, "eventName", eventName, "account", event.Account, "group", event.Group, "value", event.Value)
		case "ValidatorGroupActiveVoteRevoked":
			event := eventRaw.(*contracts.ElectionValidatorGroupActiveVoteRevoked)
			logEventLog(logger, "eventName", eventName, "account", event.Account, "group", event.Group, "value", event.Value, "voteUnits", event.Units)
		case "EpochRewardsDistributedToVoters":
			event := eventRaw.(*contracts.ElectionEpochRewardsDistributedToVoters)
			logEventLog(logger, "eventName", eventName, "group", event.Group, "value", event.Value)
		}
	}
}
