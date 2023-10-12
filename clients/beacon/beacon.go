package beacon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/marioevz/eth-clients/clients"
	"github.com/marioevz/eth-clients/clients/utils"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/client/beaconapi"
	"github.com/protolambda/eth2api/client/builderapi"
	"github.com/protolambda/eth2api/client/debugapi"
	"github.com/protolambda/eth2api/client/nodeapi"
	"github.com/protolambda/eth2api/client/validatorapi"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/deneb"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/configs"
	"github.com/protolambda/ztyp/tree"
)

const (
	PortBeaconTCP    = 9000
	PortBeaconUDP    = 9000
	PortBeaconAPI    = 4000
	PortBeaconGRPC   = 4001
	PortMetrics      = 8080
	PortValidatorAPI = 5000
)

var EMPTY_TREE_ROOT = tree.Root{}

type BeaconClientConfig struct {
	ClientIndex             int
	TerminalTotalDifficulty int64
	BeaconAPIPort           int
	Spec                    *common.Spec
	GenesisValidatorsRoot   *tree.Root
	GenesisTime             *common.Timestamp
	Subnet                  string
}

type BeaconClient struct {
	clients.Client
	Logger  utils.Logging
	Config  BeaconClientConfig
	Builder interface{}

	api *eth2api.Eth2HttpClient
}

func (bn *BeaconClient) Logf(format string, values ...interface{}) {
	if l := bn.Logger; l != nil {
		l.Logf(format, values...)
	}
}

func (bn *BeaconClient) Start() error {
	if !bn.IsRunning() {
		if managedClient, ok := bn.Client.(clients.ManagedClient); !ok {
			return fmt.Errorf("attempted to start an unmanaged client")
		} else {
			if err := managedClient.Start(); err != nil {
				return err
			}
		}

	}

	return bn.Init(context.Background())
}

func (bn *BeaconClient) Init(ctx context.Context) error {
	if bn.api == nil {
		port := bn.Config.BeaconAPIPort
		if port == 0 {
			port = PortBeaconAPI
		}
		bn.api = &eth2api.Eth2HttpClient{
			Addr: fmt.Sprintf(
				"http://%s:%d",
				bn.GetIP(),
				port,
			),
			Cli:   &http.Client{},
			Codec: eth2api.JSONCodec{},
		}
	}

	var wg sync.WaitGroup
	var errs = make(chan error, 2)
	if bn.Config.Spec == nil {
		// Try to fetch config directly from the client
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if cfg, err := bn.BeaconConfig(ctx); err == nil && cfg != nil {
					if spec, err := SpecFromConfig(cfg); err != nil {
						errs <- err
						return
					} else {
						bn.Config.Spec = spec
						return
					}
				}
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case <-time.After(time.Second):
				}
			}
		}()
	}

	if bn.Config.GenesisTime == nil || bn.Config.GenesisValidatorsRoot == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if gen, err := bn.GenesisConfig(ctx); err == nil &&
					gen != nil {
					bn.Config.GenesisTime = &gen.GenesisTime
					bn.Config.GenesisValidatorsRoot = &gen.GenesisValidatorsRoot
					return
				}
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case <-time.After(time.Second):
				}
			}
		}()
	}
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}

func (bn *BeaconClient) Shutdown() error {
	if managedClient, ok := bn.Client.(clients.ManagedClient); !ok {
		return fmt.Errorf("attempted to shutdown an unmanaged client")
	} else {
		return managedClient.Shutdown()
	}
}

func (bn *BeaconClient) ENR(parentCtx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	var out eth2api.NetworkIdentity
	if err := nodeapi.Identity(ctx, bn.api, &out); err != nil {
		return "", err
	}
	return out.ENR, nil
}

func (bn *BeaconClient) P2PAddr(parentCtx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	var out eth2api.NetworkIdentity
	if err := nodeapi.Identity(ctx, bn.api, &out); err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"/ip4/%s/tcp/%d/p2p/%s",
		bn.GetIP().String(),
		PortBeaconTCP,
		out.PeerID,
	), nil
}

func (bn *BeaconClient) BeaconAPIURL() (string, error) {
	if bn.api == nil {
		return "", fmt.Errorf("api not initialized")
	}
	return bn.api.Addr, nil
}

func (bn *BeaconClient) EnodeURL() (string, error) {
	return "", errors.New(
		"beacon node does not have an discv4 Enode URL, use ENR or multi-address instead",
	)
}

func (bn *BeaconClient) ClientName() string {
	name := bn.ClientType()
	if len(name) > 3 && name[len(name)-3:] == "-bn" {
		name = name[:len(name)-3]
	}
	return name
}

func (bn *BeaconClient) API() *eth2api.Eth2HttpClient {
	return bn.api
}

func SpecFromConfig(cfg *common.Config) (*common.Spec, error) {
	if cfg == nil {
		return nil, fmt.Errorf("empty cfg")
	}
	var spec *common.Spec
	if cfg.PRESET_BASE == "mainnet" {
		specCpy := *configs.Mainnet
		spec = &specCpy
	} else if cfg.PRESET_BASE == "minimal" {
		specCpy := *configs.Minimal
		spec = &specCpy
	} else {
		return nil, fmt.Errorf("invalid preset base: %s", cfg.PRESET_BASE)
	}
	spec.Config = *cfg
	return spec, nil
}

// Beacon API wrappers
func (bn *BeaconClient) BeaconConfig(
	parentCtx context.Context,
) (*common.Config, error) {
	var (
		cfg    = new(common.Config)
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = eth2api.SimpleRequest(
		ctx,
		bn.api,
		eth2api.FmtGET("/eth/v1/config/spec"),
		eth2api.Wrap(cfg),
	)

	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return cfg, err
}

func (bn *BeaconClient) GenesisConfig(
	parentCtx context.Context,
) (*eth2api.GenesisResponse, error) {
	var (
		dest   = new(eth2api.GenesisResponse)
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()

	exists, err = beaconapi.Genesis(ctx, bn.api, dest)

	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return dest, err
}

func (bn *BeaconClient) BlockV2Root(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (tree.Root, error) {
	var (
		root   tree.Root
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	root, exists, err = beaconapi.BlockRoot(ctx, bn.api, blockId)
	if !exists {
		return root, fmt.Errorf(
			"endpoint not found on beacon client",
		)
	}
	return root, err
}

func (bn *BeaconClient) BlockV2(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (*VersionedSignedBeaconBlock, error) {
	var (
		versionedBlock = new(eth2api.VersionedSignedBeaconBlock)
		exists         bool
		err            error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.BlockV2(ctx, bn.api, blockId, versionedBlock)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return &VersionedSignedBeaconBlock{
		VersionedSignedBeaconBlock: versionedBlock,
		spec:                       bn.Config.Spec,
	}, err
}

type BlockV2OptimisticResponse struct {
	Version             string `json:"version"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
}

func (bn *BeaconClient) BlockIsOptimistic(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (bool, error) {
	var (
		blockOptResp = new(BlockV2OptimisticResponse)
		exists       bool
		err          error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = eth2api.SimpleRequest(
		ctx,
		bn.api,
		eth2api.FmtGET("/eth/v2/beacon/blocks/%s", blockId.BlockId()),
		blockOptResp,
	)
	if !exists {
		return false, fmt.Errorf("endpoint not found on beacon client")
	}
	return blockOptResp.ExecutionOptimistic, err
}

func (bn *BeaconClient) BlockHeader(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (*eth2api.BeaconBlockHeaderAndInfo, error) {
	var (
		headInfo = new(eth2api.BeaconBlockHeaderAndInfo)
		exists   bool
		err      error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.BlockHeader(ctx, bn.api, blockId, headInfo)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return headInfo, err
}

func (bn *BeaconClient) BlobSidecars(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) ([]deneb.BlobSidecar, error) {
	var (
		blobSidecars = new([]deneb.BlobSidecar)
		exists       bool
		err          error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.BlobSidecars(ctx, bn.api, blockId, blobSidecars)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return *blobSidecars, err
}

func (bn *BeaconClient) StateValidator(
	parentCtx context.Context,
	stateId eth2api.StateId,
	validatorId eth2api.ValidatorId,
) (*eth2api.ValidatorResponse, error) {
	var (
		stateValidatorResponse = new(eth2api.ValidatorResponse)
		exists                 bool
		err                    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.StateValidator(
		ctx,
		bn.api,
		stateId,
		validatorId,
		stateValidatorResponse,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return stateValidatorResponse, err
}

func (bn *BeaconClient) ProposerIndex(
	parentCtx context.Context,
	slot common.Slot,
) (common.ValidatorIndex, error) {
	var (
		proposerDutyResponse = new(eth2api.DependentProposerDuty)
		epoch                = bn.Config.Spec.SlotToEpoch(slot)
		syncing              bool
		err                  error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	syncing, err = validatorapi.ProposerDuties(
		ctx,
		bn.api,
		epoch,
		proposerDutyResponse,
	)
	if err != nil {
		return 0, err
	}
	if syncing {
		return 0, fmt.Errorf("beacon client is syncing")
	}
	if proposerDutyResponse.Data == nil {
		return 0, fmt.Errorf("no proposer duty data")
	}
	for _, duty := range proposerDutyResponse.Data {
		if duty.Slot == slot {
			return duty.ValidatorIndex, nil
		}
	}
	return 0, fmt.Errorf("no proposer duty found for slot %d", slot)
}

func (bn *BeaconClient) StateFinalityCheckpoints(
	parentCtx context.Context,
	stateId eth2api.StateId,
) (*eth2api.FinalityCheckpoints, error) {
	var (
		finalityCheckpointsResponse = new(eth2api.FinalityCheckpoints)
		exists                      bool
		err                         error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.FinalityCheckpoints(
		ctx,
		bn.api,
		stateId,
		finalityCheckpointsResponse,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return finalityCheckpointsResponse, err
}

func (bn *BeaconClient) StateFork(
	parentCtx context.Context,
	stateId eth2api.StateId,
) (*common.Fork, error) {
	var (
		fork   = new(common.Fork)
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.Fork(
		ctx,
		bn.api,
		stateId,
		fork,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return fork, err
}

func (bn *BeaconClient) StateRandaoMix(
	parentCtx context.Context,
	stateId eth2api.StateId,
) (*tree.Root, error) {
	var (
		resp   = new(eth2api.RandaoMixResponse)
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.RandaoMix(
		ctx,
		bn.api,
		stateId,
		resp,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	if err != nil {
		return nil, err
	}
	return &resp.RandaoMix, err
}

func (bn *BeaconClient) ExpectedWithdrawals(
	parentCtx context.Context,
	stateId eth2api.StateId,
) (common.Withdrawals, error) {
	var (
		resp   = new(common.Withdrawals)
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = builderapi.ExpectedWithdrawals(
		ctx,
		bn.api,
		stateId,
		resp,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("nil expected withdrawals")
	}
	return *resp, err
}

func (bn *BeaconClient) BlockFinalityCheckpoints(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (*eth2api.FinalityCheckpoints, error) {
	var (
		headInfo                    *eth2api.BeaconBlockHeaderAndInfo
		finalityCheckpointsResponse *eth2api.FinalityCheckpoints
		err                         error
	)
	headInfo, err = bn.BlockHeader(parentCtx, blockId)
	if err != nil {
		return nil, err
	}
	finalityCheckpointsResponse, err = bn.StateFinalityCheckpoints(
		parentCtx,
		eth2api.StateIdRoot(headInfo.Header.Message.StateRoot),
	)
	if err != nil {
		// Try again using slot number
		return bn.StateFinalityCheckpoints(
			parentCtx,
			eth2api.StateIdSlot(headInfo.Header.Message.Slot),
		)
	}
	return finalityCheckpointsResponse, err
}

func Eth1WithdrawalCredential(validator *phase0.Validator) common.Eth1Address {
	var address common.Eth1Address
	copy(address[:], validator.WithdrawalCredentials[12:])
	return address
}

func IsFullyWithdrawableValidator(
	validator *phase0.Validator,
	balance common.Gwei,
	epoch common.Epoch,
) bool {
	return HasEth1WithdrawalCredential(validator) &&
		validator.WithdrawableEpoch <= epoch &&
		balance > 0
}

func IsPartiallyWithdrawableValidator(
	spec *common.Spec,
	validator *phase0.Validator,
	balance common.Gwei,
	epoch common.Epoch,
) bool {
	effectiveBalance := validator.EffectiveBalance
	hasMaxEffectiveBalance := effectiveBalance == spec.MAX_EFFECTIVE_BALANCE
	hasExcessBalance := balance > spec.MAX_EFFECTIVE_BALANCE
	return HasEth1WithdrawalCredential(validator) && hasMaxEffectiveBalance &&
		hasExcessBalance
}

func HasEth1WithdrawalCredential(validator *phase0.Validator) bool {
	return bytes.Equal(
		validator.WithdrawalCredentials[:1],
		[]byte{common.ETH1_ADDRESS_WITHDRAWAL_PREFIX},
	)
}

func (bn *BeaconClient) BeaconStateV2(
	parentCtx context.Context,
	stateId eth2api.StateId,
) (*VersionedBeaconStateResponse, error) {
	var (
		versionedBeaconStateResponse = new(eth2api.VersionedBeaconState)
		exists                       bool
		err                          error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = debugapi.BeaconStateV2(
		ctx,
		bn.api,
		stateId,
		versionedBeaconStateResponse,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return &VersionedBeaconStateResponse{
		VersionedBeaconState: versionedBeaconStateResponse,
		spec:                 bn.Config.Spec,
	}, err
}

func (bn *BeaconClient) BeaconStateV2ByBlock(
	parentCtx context.Context,
	blockId eth2api.BlockId,
) (*VersionedBeaconStateResponse, error) {
	var (
		headInfo *eth2api.BeaconBlockHeaderAndInfo
		err      error
	)
	headInfo, err = bn.BlockHeader(parentCtx, blockId)
	if err != nil {
		return nil, err
	}
	return bn.BeaconStateV2(
		parentCtx,
		eth2api.StateIdRoot(headInfo.Header.Message.StateRoot),
	)
}

func (bn *BeaconClient) StateValidators(
	parentCtx context.Context,
	stateId eth2api.StateId,
	validatorIds []eth2api.ValidatorId,
	statusFilter []eth2api.ValidatorStatus,
) ([]eth2api.ValidatorResponse, error) {
	var (
		stateValidatorResponse = make(
			[]eth2api.ValidatorResponse,
			0,
		)
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.StateValidators(
		ctx,
		bn.api,
		stateId,
		validatorIds,
		statusFilter,
		&stateValidatorResponse,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return stateValidatorResponse, err
}

func (bn *BeaconClient) StateValidatorBalances(
	parentCtx context.Context,
	stateId eth2api.StateId,
	validatorIds []eth2api.ValidatorId,
) ([]eth2api.ValidatorBalanceResponse, error) {
	var (
		stateValidatorBalanceResponse = make(
			[]eth2api.ValidatorBalanceResponse,
			0,
		)
		exists bool
		err    error
	)
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	exists, err = beaconapi.StateValidatorBalances(
		ctx,
		bn.api,
		stateId,
		validatorIds,
		&stateValidatorBalanceResponse,
	)
	if !exists {
		return nil, fmt.Errorf("endpoint not found on beacon client")
	}
	return stateValidatorBalanceResponse, err
}

func (bn *BeaconClient) ComputeDomain(
	ctx context.Context,
	typ common.BLSDomainType,
	version *common.Version,
) (common.BLSDomain, error) {
	if bn.Config.GenesisTime == nil {
		panic(fmt.Errorf("init not called yet"))
	}
	if version != nil {
		return common.ComputeDomain(
			typ,
			*version,
			*bn.Config.GenesisValidatorsRoot,
		), nil
	}
	// We need to request for head state to know current active version
	state, err := bn.BeaconStateV2ByBlock(ctx, eth2api.BlockHead)
	if err != nil {
		return common.BLSDomain{}, err
	}
	return common.ComputeDomain(
		typ,
		state.CurrentVersion(),
		*bn.Config.GenesisValidatorsRoot,
	), nil
}

func (bn *BeaconClient) SubmitPoolBLSToExecutionChange(
	parentCtx context.Context,
	l common.SignedBLSToExecutionChanges,
) error {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return beaconapi.SubmitBLSToExecutionChanges(ctx, bn.api, l)
}

func (bn *BeaconClient) SubmitVoluntaryExit(
	parentCtx context.Context,
	exit *phase0.SignedVoluntaryExit,
) error {
	ctx, cancel := utils.ContextTimeoutRPC(parentCtx)
	defer cancel()
	return beaconapi.SubmitVoluntaryExit(ctx, bn.api, exit)
}

func (b *BeaconClient) WaitForExecutionPayload(
	ctx context.Context,
) (ethcommon.Hash, error) {
	if b.Config.GenesisTime == nil {
		panic(fmt.Errorf("init not called yet"))
	}
	b.Logf(
		"Waiting for execution payload on beacon %d (%s)\n",
		b.Config.ClientIndex,
		b.ClientName(),
	)
	slotDuration := time.Duration(b.Config.Spec.SECONDS_PER_SLOT) * time.Second
	timer := time.NewTicker(slotDuration)

	for {
		select {
		case <-ctx.Done():
			return ethcommon.Hash{}, ctx.Err()
		case <-timer.C:
			realTimeSlot := b.Config.Spec.TimeToSlot(
				common.Timestamp(time.Now().Unix()),
				*b.Config.GenesisTime,
			)
			var (
				headInfo  *eth2api.BeaconBlockHeaderAndInfo
				err       error
				execution ethcommon.Hash
			)
			if headInfo, err = b.BlockHeader(ctx, eth2api.BlockHead); err != nil {
				return ethcommon.Hash{}, err
			}
			if !headInfo.Canonical {
				continue
			}

			if versionedBlock, err := b.BlockV2(ctx, eth2api.BlockIdRoot(headInfo.Root)); err != nil {
				continue
			} else if executionPayload, _, _, err := versionedBlock.ExecutionPayload(); err == nil {
				copy(
					execution[:],
					executionPayload.BlockHash[:],
				)
			}
			zero := ethcommon.Hash{}
			b.Logf(
				"WaitForExecutionPayload: beacon %d (%s): slot=%d, realTimeSlot=%d, head=%s, exec=%s\n",
				b.Config.ClientIndex,
				b.ClientName(),
				headInfo.Header.Message.Slot,
				realTimeSlot,
				utils.Shorten(headInfo.Root.String()),
				utils.Shorten(execution.Hex()),
			)
			if !bytes.Equal(execution[:], zero[:]) {
				return execution, nil
			}
		}
	}
}

func (b *BeaconClient) WaitForOptimisticState(
	ctx context.Context,
	blockID eth2api.BlockId,
	optimistic bool,
) (*eth2api.BeaconBlockHeaderAndInfo, error) {
	b.Logf("Waiting for optimistic sync on beacon %d (%s)\n",
		b.Config.ClientIndex,
		b.ClientName(),
	)

	slotDuration := time.Duration(b.Config.Spec.SECONDS_PER_SLOT) * time.Second
	timer := time.NewTicker(slotDuration)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
			var headOptStatus BlockV2OptimisticResponse
			if exists, err := eth2api.SimpleRequest(ctx, b.api, eth2api.FmtGET("/eth/v2/beacon/blocks/%s", blockID.BlockId()), &headOptStatus); err != nil {
				// Block still not synced
				continue
			} else if !exists {
				// Block still not synced
				continue
			}
			if headOptStatus.ExecutionOptimistic != optimistic {
				continue
			}
			// Return the block
			var blockInfo eth2api.BeaconBlockHeaderAndInfo
			if exists, err := beaconapi.BlockHeader(ctx, b.api, blockID, &blockInfo); err != nil {
				return nil, fmt.Errorf(
					"WaitForExecutionPayload: failed to poll block: %v",
					err,
				)
			} else if !exists {
				return nil, fmt.Errorf("WaitForExecutionPayload: failed to poll block: !exists")
			}
			return &blockInfo, nil
		}
	}
}

func (bn *BeaconClient) GetLatestExecutionBeaconBlock(
	parentCtx context.Context,
) (*VersionedSignedBeaconBlock, error) {
	headInfo, err := bn.BlockHeader(parentCtx, eth2api.BlockHead)
	if err != nil {
		return nil, fmt.Errorf("failed to poll head: %v", err)
	}
	for slot := headInfo.Header.Message.Slot; slot > 0; slot-- {
		versionedBlock, err := bn.BlockV2(parentCtx, eth2api.BlockIdSlot(slot))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve block: %v", err)
		}
		if executionPayload, _, _, err := versionedBlock.ExecutionPayload(); err == nil {
			if !bytes.Equal(
				executionPayload.BlockHash[:],
				EMPTY_TREE_ROOT[:],
			) {
				return versionedBlock, nil
			}
		}
	}
	return nil, nil
}

func (bn *BeaconClient) GetFirstExecutionBeaconBlock(
	parentCtx context.Context,
) (*VersionedSignedBeaconBlock, error) {
	if bn.Config.GenesisTime == nil {
		panic(fmt.Errorf("init not called yet"))
	}
	lastSlot := bn.Config.Spec.TimeToSlot(
		common.Timestamp(time.Now().Unix()),
		*bn.Config.GenesisTime,
	)
	for slot := common.Slot(0); slot <= lastSlot; slot++ {
		versionedBlock, err := bn.BlockV2(parentCtx, eth2api.BlockIdSlot(slot))
		if err != nil {
			continue
		}
		if executionPayload, _, _, err := versionedBlock.ExecutionPayload(); err == nil {
			if !bytes.Equal(
				executionPayload.BlockHash[:],
				EMPTY_TREE_ROOT[:],
			) {
				return versionedBlock, nil
			}
		}
	}
	return nil, nil
}

func (bn *BeaconClient) GetBeaconBlockByExecutionHash(
	parentCtx context.Context,
	hash ethcommon.Hash,
) (*VersionedSignedBeaconBlock, error) {
	headInfo, err := bn.BlockHeader(parentCtx, eth2api.BlockHead)
	if err != nil {
		return nil, fmt.Errorf("failed to poll head: %v", err)
	}
	for slot := int(headInfo.Header.Message.Slot); slot > 0; slot -= 1 {
		versionedBlock, err := bn.BlockV2(parentCtx, eth2api.BlockIdSlot(slot))
		if err != nil {
			continue
		}
		if executionPayload, _, _, err := versionedBlock.ExecutionPayload(); err == nil {
			if !bytes.Equal(executionPayload.BlockHash[:], hash[:]) {
				return versionedBlock, nil
			}
		}
	}
	return nil, nil
}

func (bn *BeaconClient) GetFilledSlotsCountPerEpoch(
	parentCtx context.Context,
) (map[common.Epoch]uint64, error) {
	headInfo, err := bn.BlockHeader(parentCtx, eth2api.BlockHead)
	epochMap := make(map[common.Epoch]uint64)
	for {
		if err != nil {
			return nil, fmt.Errorf("failed to poll head: %v", err)
		}
		epoch := common.Epoch(
			headInfo.Header.Message.Slot / bn.Config.Spec.SLOTS_PER_EPOCH,
		)
		if prev, ok := epochMap[epoch]; ok {
			epochMap[epoch] = prev + 1
		} else {
			epochMap[epoch] = 1
		}
		if bytes.Equal(
			headInfo.Header.Message.ParentRoot[:],
			EMPTY_TREE_ROOT[:],
		) {
			break
		}
		headInfo, err = bn.BlockHeader(
			parentCtx,
			eth2api.BlockIdRoot(headInfo.Header.Message.ParentRoot),
		)
	}

	return epochMap, nil
}

type BeaconClients []*BeaconClient

// Return subset of clients that are currently running
func (all BeaconClients) Running() BeaconClients {
	res := make(BeaconClients, 0)
	for _, bc := range all {
		if bc.IsRunning() {
			res = append(res, bc)
		}
	}
	return res
}

// Return subset of clients that are part of an specific subnet
func (all BeaconClients) Subnet(subnet string) BeaconClients {
	if subnet == "" {
		return all
	}
	res := make(BeaconClients, 0)
	for _, bn := range all {
		if bn.Config.Subnet == subnet {
			res = append(res, bn)
		}
	}
	return res
}

// Returns comma-separated ENRs of all running beacon nodes
func (beacons BeaconClients) ENRs(parentCtx context.Context) (string, error) {
	if len(beacons) == 0 {
		return "", nil
	}
	enrs := make([]string, 0)
	for _, bn := range beacons {
		if bn.IsRunning() {
			enr, err := bn.ENR(parentCtx)
			if err != nil {
				return "", err
			}
			enrs = append(enrs, enr)
		}
	}
	return strings.Join(enrs, ","), nil
}

// Returns comma-separated P2PAddr of all running beacon nodes
func (beacons BeaconClients) P2PAddrs(
	parentCtx context.Context,
) (string, error) {
	if len(beacons) == 0 {
		return "", nil
	}
	staticPeers := make([]string, 0)
	for _, bn := range beacons {
		if bn.IsRunning() {
			p2p, err := bn.P2PAddr(parentCtx)
			if err != nil {
				return "", err
			}
			staticPeers = append(staticPeers, p2p)
		}
	}
	return strings.Join(staticPeers, ","), nil
}

func (b BeaconClients) GetBeaconBlockByExecutionHash(
	parentCtx context.Context,
	hash ethcommon.Hash,
) (*VersionedSignedBeaconBlock, error) {
	for _, bn := range b {
		block, err := bn.GetBeaconBlockByExecutionHash(parentCtx, hash)
		if err != nil || block != nil {
			return block, err
		}
	}
	return nil, nil
}

func (runningBeacons BeaconClients) PrintStatus(
	ctx context.Context,
) {
	for i, b := range runningBeacons {
		var (
			slot      common.Slot
			version   string
			head      string
			justified string
			finalized string
			execution = "0x0000..0000"
		)

		if headInfo, err := b.BlockHeader(ctx, eth2api.BlockHead); err == nil {
			slot = headInfo.Header.Message.Slot
			head = utils.Shorten(headInfo.Root.String())
		}
		checkpoints, err := b.BlockFinalityCheckpoints(
			ctx,
			eth2api.BlockHead,
		)
		if err == nil {
			justified = utils.Shorten(
				checkpoints.CurrentJustified.String(),
			)
			finalized = utils.Shorten(checkpoints.Finalized.String())
		}
		if versionedBlock, err := b.BlockV2(ctx, eth2api.BlockHead); err == nil {
			version = versionedBlock.Version
			if executionPayload, _, _, err := versionedBlock.ExecutionPayload(); err == nil {
				execution = utils.Shorten(
					executionPayload.BlockHash.String(),
				)
			}
		}

		b.Logf(
			"beacon %d (%s): fork=%s, slot=%d, head=%s, exec_payload=%s, justified=%s, finalized=%s",
			i,
			b.ClientName(),
			version,
			slot,
			head,
			execution,
			justified,
			finalized,
		)
	}
}

func (all BeaconClients) SubmitPoolBLSToExecutionChange(
	parentCtx context.Context,
	l common.SignedBLSToExecutionChanges,
) error {
	for _, b := range all {
		if err := b.SubmitPoolBLSToExecutionChange(parentCtx, l); err != nil {
			return err
		}
	}
	return nil
}
