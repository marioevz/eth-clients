package beacon

import (
	"fmt"

	api "github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
	"github.com/protolambda/zrnt/eth2/beacon/capella"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/ztyp/tree"
)

type VersionedSignedBeaconBlock struct {
	*eth2api.VersionedSignedBeaconBlock
	spec *common.Spec
}

func (versionedBlock *VersionedSignedBeaconBlock) ContainsExecutionPayload() bool {
	return versionedBlock.Version == "bellatrix" ||
		versionedBlock.Version == "capella"
}

func (versionedBlock *VersionedSignedBeaconBlock) ExecutionPayload() (api.ExecutableData, error) {
	result := api.ExecutableData{}
	switch v := versionedBlock.Data.(type) {
	case *bellatrix.SignedBeaconBlock:
		execPayload := v.Message.Body.ExecutionPayload
		copy(result.ParentHash[:], execPayload.ParentHash[:])
		copy(result.FeeRecipient[:], execPayload.FeeRecipient[:])
		copy(result.StateRoot[:], execPayload.StateRoot[:])
		copy(result.ReceiptsRoot[:], execPayload.ReceiptsRoot[:])
		copy(result.LogsBloom[:], execPayload.LogsBloom[:])
		copy(result.Random[:], execPayload.PrevRandao[:])
		result.Number = uint64(execPayload.BlockNumber)
		result.GasLimit = uint64(execPayload.GasLimit)
		result.GasUsed = uint64(execPayload.GasUsed)
		result.Timestamp = uint64(execPayload.Timestamp)
		copy(result.ExtraData[:], execPayload.ExtraData[:])
		result.BaseFeePerGas = (*uint256.Int)(&execPayload.BaseFeePerGas).ToBig()
		copy(result.BlockHash[:], execPayload.BlockHash[:])
		result.Transactions = make([][]byte, 0)
		for _, t := range execPayload.Transactions {
			result.Transactions = append(result.Transactions, t)
		}
	case *capella.SignedBeaconBlock:
		execPayload := v.Message.Body.ExecutionPayload
		copy(result.ParentHash[:], execPayload.ParentHash[:])
		copy(result.FeeRecipient[:], execPayload.FeeRecipient[:])
		copy(result.StateRoot[:], execPayload.StateRoot[:])
		copy(result.ReceiptsRoot[:], execPayload.ReceiptsRoot[:])
		copy(result.LogsBloom[:], execPayload.LogsBloom[:])
		copy(result.Random[:], execPayload.PrevRandao[:])
		result.Number = uint64(execPayload.BlockNumber)
		result.GasLimit = uint64(execPayload.GasLimit)
		result.GasUsed = uint64(execPayload.GasUsed)
		result.Timestamp = uint64(execPayload.Timestamp)
		copy(result.ExtraData[:], execPayload.ExtraData[:])
		result.BaseFeePerGas = (*uint256.Int)(&execPayload.BaseFeePerGas).ToBig()
		copy(result.BlockHash[:], execPayload.BlockHash[:])
		result.Transactions = make([][]byte, 0)
		for _, t := range execPayload.Transactions {
			result.Transactions = append(result.Transactions, t)
		}
		result.Withdrawals = make([]*types.Withdrawal, 0)
		for _, w := range execPayload.Withdrawals {
			withdrawal := new(types.Withdrawal)
			withdrawal.Index = uint64(w.Index)
			withdrawal.Validator = uint64(w.ValidatorIndex)
			copy(withdrawal.Address[:], w.Address[:])
			withdrawal.Amount = uint64(w.Amount)
			result.Withdrawals = append(result.Withdrawals, withdrawal)
		}
	default:
		return result, fmt.Errorf(
			"beacon block version can't contain execution payload",
		)
	}
	return result, nil
}

func (versionedBlock *VersionedSignedBeaconBlock) Withdrawals() (common.Withdrawals, error) {
	switch v := versionedBlock.Data.(type) {
	case *capella.SignedBeaconBlock:
		return v.Message.Body.ExecutionPayload.Withdrawals, nil
	}
	return nil, nil
}

func (b *VersionedSignedBeaconBlock) Root() tree.Root {
	switch v := b.Data.(type) {
	case *phase0.SignedBeaconBlock:
		return v.Message.HashTreeRoot(b.spec, tree.GetHashFn())
	case *altair.SignedBeaconBlock:
		return v.Message.HashTreeRoot(b.spec, tree.GetHashFn())
	case *bellatrix.SignedBeaconBlock:
		return v.Message.HashTreeRoot(b.spec, tree.GetHashFn())
	case *capella.SignedBeaconBlock:
		return v.Message.HashTreeRoot(b.spec, tree.GetHashFn())
	}
	panic("badly formatted beacon block")
}

func (b *VersionedSignedBeaconBlock) StateRoot() tree.Root {
	switch v := b.Data.(type) {
	case *phase0.SignedBeaconBlock:
		return v.Message.StateRoot
	case *altair.SignedBeaconBlock:
		return v.Message.StateRoot
	case *bellatrix.SignedBeaconBlock:
		return v.Message.StateRoot
	case *capella.SignedBeaconBlock:
		return v.Message.StateRoot
	}
	panic("badly formatted beacon block")
}

func (b *VersionedSignedBeaconBlock) ParentRoot() tree.Root {
	switch v := b.Data.(type) {
	case *phase0.SignedBeaconBlock:
		return v.Message.ParentRoot
	case *altair.SignedBeaconBlock:
		return v.Message.ParentRoot
	case *bellatrix.SignedBeaconBlock:
		return v.Message.ParentRoot
	case *capella.SignedBeaconBlock:
		return v.Message.ParentRoot
	}
	panic("badly formatted beacon block")
}

func (b *VersionedSignedBeaconBlock) Slot() common.Slot {
	switch v := b.Data.(type) {
	case *phase0.SignedBeaconBlock:
		return v.Message.Slot
	case *altair.SignedBeaconBlock:
		return v.Message.Slot
	case *bellatrix.SignedBeaconBlock:
		return v.Message.Slot
	case *capella.SignedBeaconBlock:
		return v.Message.Slot
	}
	panic("badly formatted beacon block")
}

func (b *VersionedSignedBeaconBlock) ProposerIndex() common.ValidatorIndex {
	switch v := b.Data.(type) {
	case *phase0.SignedBeaconBlock:
		return v.Message.ProposerIndex
	case *altair.SignedBeaconBlock:
		return v.Message.ProposerIndex
	case *bellatrix.SignedBeaconBlock:
		return v.Message.ProposerIndex
	case *capella.SignedBeaconBlock:
		return v.Message.ProposerIndex
	}
	panic("badly formatted beacon block")
}
