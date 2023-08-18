package beacon

import (
	"fmt"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
	"github.com/protolambda/zrnt/eth2/beacon/capella"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/deneb"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/ztyp/tree"
)

type VersionedBeaconStateResponse struct {
	*eth2api.VersionedBeaconState
	spec *common.Spec
}

func (vbs *VersionedBeaconStateResponse) Root() tree.Root {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.HashTreeRoot(vbs.spec, tree.GetHashFn())
	case *altair.BeaconState:
		return state.HashTreeRoot(vbs.spec, tree.GetHashFn())
	case *bellatrix.BeaconState:
		return state.HashTreeRoot(vbs.spec, tree.GetHashFn())
	case *capella.BeaconState:
		return state.HashTreeRoot(vbs.spec, tree.GetHashFn())
	case *deneb.BeaconState:
		return state.HashTreeRoot(vbs.spec, tree.GetHashFn())
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) GenesisTime() common.Timestamp {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.GenesisTime
	case *altair.BeaconState:
		return state.GenesisTime
	case *bellatrix.BeaconState:
		return state.GenesisTime
	case *capella.BeaconState:
		return state.GenesisTime
	case *deneb.BeaconState:
		return state.GenesisTime
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) GenesisValidatorsRoot() common.Root {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.GenesisValidatorsRoot
	case *altair.BeaconState:
		return state.GenesisValidatorsRoot
	case *bellatrix.BeaconState:
		return state.GenesisValidatorsRoot
	case *capella.BeaconState:
		return state.GenesisValidatorsRoot
	case *deneb.BeaconState:
		return state.GenesisValidatorsRoot
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) Fork() common.Fork {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Fork
	case *altair.BeaconState:
		return state.Fork
	case *bellatrix.BeaconState:
		return state.Fork
	case *capella.BeaconState:
		return state.Fork
	case *deneb.BeaconState:
		return state.Fork
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) CurrentVersion() common.Version {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Fork.CurrentVersion
	case *altair.BeaconState:
		return state.Fork.CurrentVersion
	case *bellatrix.BeaconState:
		return state.Fork.CurrentVersion
	case *capella.BeaconState:
		return state.Fork.CurrentVersion
	case *deneb.BeaconState:
		return state.Fork.CurrentVersion
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) PreviousVersion() common.Version {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Fork.PreviousVersion
	case *altair.BeaconState:
		return state.Fork.PreviousVersion
	case *bellatrix.BeaconState:
		return state.Fork.PreviousVersion
	case *capella.BeaconState:
		return state.Fork.PreviousVersion
	case *deneb.BeaconState:
		return state.Fork.PreviousVersion
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) LatestBlockHeader() common.BeaconBlockHeader {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.LatestBlockHeader
	case *altair.BeaconState:
		return state.LatestBlockHeader
	case *bellatrix.BeaconState:
		return state.LatestBlockHeader
	case *capella.BeaconState:
		return state.LatestBlockHeader
	case *deneb.BeaconState:
		return state.LatestBlockHeader
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) BlockRoots() phase0.HistoricalBatchRoots {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.BlockRoots
	case *altair.BeaconState:
		return state.BlockRoots
	case *bellatrix.BeaconState:
		return state.BlockRoots
	case *capella.BeaconState:
		return state.BlockRoots
	case *deneb.BeaconState:
		return state.BlockRoots
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) StateRoots() phase0.HistoricalBatchRoots {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.StateRoots
	case *altair.BeaconState:
		return state.StateRoots
	case *bellatrix.BeaconState:
		return state.StateRoots
	case *capella.BeaconState:
		return state.StateRoots
	case *deneb.BeaconState:
		return state.StateRoots
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) HistoricalRoots() phase0.HistoricalRoots {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.HistoricalRoots
	case *altair.BeaconState:
		return state.HistoricalRoots
	case *bellatrix.BeaconState:
		return state.HistoricalRoots
	case *capella.BeaconState:
		return state.HistoricalRoots
	case *deneb.BeaconState:
		return state.HistoricalRoots
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) Eth1Data() common.Eth1Data {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Eth1Data
	case *altair.BeaconState:
		return state.Eth1Data
	case *bellatrix.BeaconState:
		return state.Eth1Data
	case *capella.BeaconState:
		return state.Eth1Data
	case *deneb.BeaconState:
		return state.Eth1Data
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) Eth1DataVotes() phase0.Eth1DataVotes {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Eth1DataVotes
	case *altair.BeaconState:
		return state.Eth1DataVotes
	case *bellatrix.BeaconState:
		return state.Eth1DataVotes
	case *capella.BeaconState:
		return state.Eth1DataVotes
	case *deneb.BeaconState:
		return state.Eth1DataVotes
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) Eth1DepositIndex() common.DepositIndex {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Eth1DepositIndex
	case *altair.BeaconState:
		return state.Eth1DepositIndex
	case *bellatrix.BeaconState:
		return state.Eth1DepositIndex
	case *capella.BeaconState:
		return state.Eth1DepositIndex
	case *deneb.BeaconState:
		return state.Eth1DepositIndex
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) Balances() phase0.Balances {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Balances
	case *altair.BeaconState:
		return state.Balances
	case *bellatrix.BeaconState:
		return state.Balances
	case *capella.BeaconState:
		return state.Balances
	case *deneb.BeaconState:
		return state.Balances
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) Balance(
	id common.ValidatorIndex,
) common.Gwei {
	balances := vbs.Balances()
	if int(id) >= len(balances) {
		panic("invalid validator requested")
	}
	return balances[id]
}

func (vbs *VersionedBeaconStateResponse) Validators() phase0.ValidatorRegistry {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Validators
	case *altair.BeaconState:
		return state.Validators
	case *bellatrix.BeaconState:
		return state.Validators
	case *capella.BeaconState:
		return state.Validators
	case *deneb.BeaconState:
		return state.Validators
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) RandaoMixes() phase0.RandaoMixes {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.RandaoMixes
	case *altair.BeaconState:
		return state.RandaoMixes
	case *bellatrix.BeaconState:
		return state.RandaoMixes
	case *capella.BeaconState:
		return state.RandaoMixes
	case *deneb.BeaconState:
		return state.RandaoMixes
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) Slashings() phase0.SlashingsHistory {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Slashings
	case *altair.BeaconState:
		return state.Slashings
	case *bellatrix.BeaconState:
		return state.Slashings
	case *capella.BeaconState:
		return state.Slashings
	case *deneb.BeaconState:
		return state.Slashings
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) PreviousEpochAttestations() phase0.PendingAttestations {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.PreviousEpochAttestations
	case *altair.BeaconState:
		return nil
	case *bellatrix.BeaconState:
		return nil
	case *capella.BeaconState:
		return nil
	case *deneb.BeaconState:
		return nil
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) CurrentEpochAttestations() phase0.PendingAttestations {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.CurrentEpochAttestations
	case *altair.BeaconState:
		return nil
	case *bellatrix.BeaconState:
		return nil
	case *capella.BeaconState:
		return nil
	case *deneb.BeaconState:
		return nil
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) PreviousEpochParticipation() altair.ParticipationRegistry {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return nil
	case *altair.BeaconState:
		return state.PreviousEpochParticipation
	case *bellatrix.BeaconState:
		return state.PreviousEpochParticipation
	case *capella.BeaconState:
		return state.PreviousEpochParticipation
	case *deneb.BeaconState:
		return state.PreviousEpochParticipation
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) CurrentEpochParticipation() altair.ParticipationRegistry {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return nil
	case *altair.BeaconState:
		return state.CurrentEpochParticipation
	case *bellatrix.BeaconState:
		return state.CurrentEpochParticipation
	case *capella.BeaconState:
		return state.CurrentEpochParticipation
	case *deneb.BeaconState:
		return state.CurrentEpochParticipation
	}
	panic("badly formatted beacon state")
}

// Finality
func (vbs *VersionedBeaconStateResponse) JustificationBits() common.JustificationBits {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.JustificationBits
	case *altair.BeaconState:
		return state.JustificationBits
	case *bellatrix.BeaconState:
		return state.JustificationBits
	case *capella.BeaconState:
		return state.JustificationBits
	case *deneb.BeaconState:
		return state.JustificationBits
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) PreviousJustifiedCheckpoint() common.Checkpoint {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.PreviousJustifiedCheckpoint
	case *altair.BeaconState:
		return state.PreviousJustifiedCheckpoint
	case *bellatrix.BeaconState:
		return state.PreviousJustifiedCheckpoint
	case *capella.BeaconState:
		return state.PreviousJustifiedCheckpoint
	case *deneb.BeaconState:
		return state.PreviousJustifiedCheckpoint
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) CurrentJustifiedCheckpoint() common.Checkpoint {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.CurrentJustifiedCheckpoint
	case *altair.BeaconState:
		return state.CurrentJustifiedCheckpoint
	case *bellatrix.BeaconState:
		return state.CurrentJustifiedCheckpoint
	case *capella.BeaconState:
		return state.CurrentJustifiedCheckpoint
	case *deneb.BeaconState:
		return state.CurrentJustifiedCheckpoint
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) FinalizedCheckpoint() common.Checkpoint {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.FinalizedCheckpoint
	case *altair.BeaconState:
		return state.FinalizedCheckpoint
	case *bellatrix.BeaconState:
		return state.FinalizedCheckpoint
	case *capella.BeaconState:
		return state.FinalizedCheckpoint
	case *deneb.BeaconState:
		return state.FinalizedCheckpoint
	}
	panic("badly formatted beacon state")
}

// Altair
// Inactivity
func (vbs *VersionedBeaconStateResponse) InactivityScores() altair.InactivityScores {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return nil
	case *altair.BeaconState:
		return state.InactivityScores
	case *bellatrix.BeaconState:
		return state.InactivityScores
	case *capella.BeaconState:
		return state.InactivityScores
	case *deneb.BeaconState:
		return state.InactivityScores
	}
	panic("badly formatted beacon state")
}

// Sync
func (vbs *VersionedBeaconStateResponse) CurrentSyncCommittee() *common.SyncCommittee {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return nil
	case *altair.BeaconState:
		return &state.CurrentSyncCommittee
	case *bellatrix.BeaconState:
		return &state.CurrentSyncCommittee
	case *capella.BeaconState:
		return &state.CurrentSyncCommittee
	case *deneb.BeaconState:
		return &state.CurrentSyncCommittee
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) NextSyncCommittee() *common.SyncCommittee {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return nil
	case *altair.BeaconState:
		return &state.NextSyncCommittee
	case *bellatrix.BeaconState:
		return &state.NextSyncCommittee
	case *capella.BeaconState:
		return &state.NextSyncCommittee
	case *deneb.BeaconState:
		return &state.NextSyncCommittee
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) StateSlot() common.Slot {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return state.Slot
	case *altair.BeaconState:
		return state.Slot
	case *bellatrix.BeaconState:
		return state.Slot
	case *capella.BeaconState:
		return state.Slot
	case *deneb.BeaconState:
		return state.Slot
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) LatestExecutionPayloadHeaderHash() tree.Root {
	switch state := vbs.Data.(type) {
	case *phase0.BeaconState:
		return tree.Root{}
	case *altair.BeaconState:
		return tree.Root{}
	case *bellatrix.BeaconState:
		return state.LatestExecutionPayloadHeader.BlockHash
	case *capella.BeaconState:
		return state.LatestExecutionPayloadHeader.BlockHash
	case *deneb.BeaconState:
		return state.LatestExecutionPayloadHeader.BlockHash
	}
	panic("badly formatted beacon state")
}

func (vbs *VersionedBeaconStateResponse) NextWithdrawalIndex() (common.WithdrawalIndex, error) {
	var wIndex common.WithdrawalIndex
	switch state := vbs.Data.(type) {
	case *capella.BeaconState:
		wIndex = state.NextWithdrawalIndex
	case *deneb.BeaconState:
		wIndex = state.NextWithdrawalIndex
	}
	return wIndex, nil
}

func (vbs *VersionedBeaconStateResponse) NextWithdrawalValidatorIndex() (common.ValidatorIndex, error) {
	var wIndex common.ValidatorIndex
	switch state := vbs.Data.(type) {
	case *capella.BeaconState:
		wIndex = state.NextWithdrawalValidatorIndex
	case *deneb.BeaconState:
		wIndex = state.NextWithdrawalValidatorIndex
	}
	return wIndex, nil
}

func (vbs *VersionedBeaconStateResponse) NextWithdrawals(
	slot common.Slot,
) (common.Withdrawals, error) {
	var (
		withdrawalIndex common.WithdrawalIndex
		validatorIndex  common.ValidatorIndex
		validators      phase0.ValidatorRegistry
		balances        phase0.Balances
		epoch           = vbs.spec.SlotToEpoch(slot)
	)
	switch state := vbs.Data.(type) {
	case *bellatrix.BeaconState:
		// withdrawalIndex and validatorIndex start at zero
		validators = state.Validators
		balances = state.Balances
	case *capella.BeaconState:
		withdrawalIndex = state.NextWithdrawalIndex
		validatorIndex = state.NextWithdrawalValidatorIndex
		validators = state.Validators
		balances = state.Balances
	case *deneb.BeaconState:
		withdrawalIndex = state.NextWithdrawalIndex
		validatorIndex = state.NextWithdrawalValidatorIndex
		validators = state.Validators
		balances = state.Balances
	default:
		return nil, fmt.Errorf("badly formatted beacon state")
	}
	validatorCount := uint64(len(validators))
	withdrawals := make(common.Withdrawals, 0)

	i := uint64(0)
	for {
		if validatorIndex >= common.ValidatorIndex(validatorCount) ||
			validatorIndex >= common.ValidatorIndex(len(balances)) {
			return nil, fmt.Errorf("invalid validator index")
		}
		validator := validators[validatorIndex]
		if validator == nil {
			return nil, fmt.Errorf("invalid validator")
		}
		balance := balances[validatorIndex]
		if i >= validatorCount ||
			i >= uint64(vbs.spec.MAX_VALIDATORS_PER_WITHDRAWALS_SWEEP) {
			break
		}
		if IsFullyWithdrawableValidator(validator, balance, epoch) {
			withdrawals = append(withdrawals, common.Withdrawal{
				Index:          withdrawalIndex,
				ValidatorIndex: validatorIndex,
				Address:        Eth1WithdrawalCredential(validator),
				Amount:         balance,
			})
			withdrawalIndex += 1
		} else if IsPartiallyWithdrawableValidator(vbs.spec, validator, balance, epoch) {
			withdrawals = append(withdrawals, common.Withdrawal{
				Index:          withdrawalIndex,
				ValidatorIndex: validatorIndex,
				Address:        Eth1WithdrawalCredential(validator),
				Amount:         balance - vbs.spec.MAX_EFFECTIVE_BALANCE,
			})
			withdrawalIndex += 1
		}
		if len(withdrawals) == int(vbs.spec.MAX_WITHDRAWALS_PER_PAYLOAD) {
			break
		}
		validatorIndex = common.ValidatorIndex(
			uint64(validatorIndex+1) % validatorCount,
		)
		i += 1
	}
	return withdrawals, nil
}
