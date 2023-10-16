package beacon

import (
	"crypto/sha256"
	"fmt"

	api "github.com/ethereum/go-ethereum/beacon/engine"
	el_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
	"github.com/protolambda/zrnt/eth2/beacon/capella"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/deneb"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/ztyp/tree"
)

var (
	BLOB_COMMITMENT_VERSION_KZG = byte(0x01)
)

type VersionedSignedBeaconBlock struct {
	*eth2api.VersionedSignedBeaconBlock
	spec *common.Spec
}

func (versionedBlock *VersionedSignedBeaconBlock) ContainsExecutionPayload() bool {
	return versionedBlock.Version == "bellatrix" ||
		versionedBlock.Version == "capella" ||
		versionedBlock.Version == "deneb"
}

func (versionedBlock *VersionedSignedBeaconBlock) ContainsKZGCommitments() bool {
	return versionedBlock.Version == "deneb"
}

func KZGCommitmentsToVersionedHashes(kzgCommitments []common.KZGCommitment) []el_common.Hash {
	versionedHashes := make([]el_common.Hash, len(kzgCommitments))
	for i, kzgCommitment := range kzgCommitments {
		sha256Hash := sha256.Sum256(kzgCommitment[:])
		versionedHashes[i] = el_common.BytesToHash(append([]byte{BLOB_COMMITMENT_VERSION_KZG}, sha256Hash[1:]...))
	}
	return versionedHashes
}

func (versionedBlock *VersionedSignedBeaconBlock) ExecutionPayload() (api.ExecutableData, []el_common.Hash, *el_common.Hash, error) {
	var (
		result          api.ExecutableData
		versionedHashes []el_common.Hash
		beaconRoot      *el_common.Hash
		err             error
	)

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
	case *deneb.SignedBeaconBlock:
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
		versionedHashes = KZGCommitmentsToVersionedHashes(v.Message.Body.BlobKZGCommitments)
		beaconRoot = &el_common.Hash{}
		copy(beaconRoot[:], v.Message.ParentRoot[:])
	default:
		err = fmt.Errorf(
			"beacon block version can't contain execution payload",
		)
	}
	return result, versionedHashes, beaconRoot, err
}

func (versionedBlock *VersionedSignedBeaconBlock) Withdrawals() (common.Withdrawals, error) {
	switch v := versionedBlock.Data.(type) {
	case *capella.SignedBeaconBlock:
		return v.Message.Body.ExecutionPayload.Withdrawals, nil
	case *deneb.SignedBeaconBlock:
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
	case *deneb.SignedBeaconBlock:
		return v.Message.HashTreeRoot(b.spec, tree.GetHashFn())
	}
	panic(fmt.Errorf("badly formatted beacon block, type=%T", b.Data))
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
	case *deneb.SignedBeaconBlock:
		return v.Message.StateRoot
	}
	panic(fmt.Errorf("badly formatted beacon block, type=%T", b.Data))
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
	case *deneb.SignedBeaconBlock:
		return v.Message.ParentRoot
	}
	panic(fmt.Errorf("badly formatted beacon block, type=%T", b.Data))
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
	case *deneb.SignedBeaconBlock:
		return v.Message.Slot
	}
	panic(fmt.Errorf("badly formatted beacon block, type=%T", b.Data))
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
	case *deneb.SignedBeaconBlock:
		return v.Message.ProposerIndex
	}
	panic(fmt.Errorf("badly formatted beacon block, type=%T", b.Data))
}

func (b *VersionedSignedBeaconBlock) ExecutionPayloadBlockHash() *tree.Root {
	switch v := b.Data.(type) {
	case *phase0.SignedBeaconBlock:
		return nil
	case *altair.SignedBeaconBlock:
		return nil
	case *bellatrix.SignedBeaconBlock:
		return &v.Message.Body.ExecutionPayload.BlockHash
	case *capella.SignedBeaconBlock:
		return &v.Message.Body.ExecutionPayload.BlockHash
	case *deneb.SignedBeaconBlock:
		return &v.Message.Body.ExecutionPayload.BlockHash
	}
	panic(fmt.Errorf("badly formatted beacon block, type=%T", b.Data))
}

func (b *VersionedSignedBeaconBlock) KZGCommitments() common.KZGCommitments {
	switch v := b.Data.(type) {
	case *phase0.SignedBeaconBlock:
		return nil
	case *altair.SignedBeaconBlock:
		return nil
	case *bellatrix.SignedBeaconBlock:
		return nil
	case *capella.SignedBeaconBlock:
		return nil
	case *deneb.SignedBeaconBlock:
		return v.Message.Body.BlobKZGCommitments
	}
	panic(fmt.Errorf("badly formatted beacon block, type=%T", b.Data))
}
