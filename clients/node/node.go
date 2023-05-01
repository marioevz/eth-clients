package node

import (
	"context"
	"fmt"
	"strings"

	"github.com/marioevz/eth-clients/clients/beacon"
	"github.com/marioevz/eth-clients/clients/execution"
	"github.com/marioevz/eth-clients/clients/utils"
	"github.com/marioevz/eth-clients/clients/validator"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// A node bundles together:
// - Running Execution client
// - Running Beacon client
// - Running Validator client
// Contains a flag that marks a node that can be used to query
// test verification information.
type Node struct {
	// Logging interface for all the events that happen in the node
	Logging utils.Logging
	// Index of the node in the network/testnet
	Index int
	// Clients that comprise the node
	ExecutionClient *execution.ExecutionClient
	BeaconClient    *beacon.BeaconClient
	ValidatorClient *validator.ValidatorClient
	// Whether this node can be used to query test verification information
	Verification bool
}

func (n *Node) Logf(format string, values ...interface{}) {
	if l := n.Logging; l != nil {
		l.Logf(format, values...)
	}
}

// Starts all clients included in the bundle
func (n *Node) Start() error {
	n.Logf("Starting validator client bundle %d", n.Index)
	if n.ExecutionClient != nil {
		if err := n.ExecutionClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No execution client started")
	}
	if n.BeaconClient != nil {
		if err := n.BeaconClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No beacon client started")
	}
	if n.ValidatorClient != nil {
		if err := n.ValidatorClient.Start(); err != nil {
			return err
		}
	} else {
		n.Logf("No validator client started")
	}
	return nil
}

func (n *Node) Shutdown() error {
	if err := n.ExecutionClient.Shutdown(); err != nil {
		return err
	}
	if err := n.BeaconClient.Shutdown(); err != nil {
		return err
	}
	if err := n.ValidatorClient.Shutdown(); err != nil {
		return err
	}
	return nil
}

func (n *Node) ClientNames() string {
	var name string
	if n.ExecutionClient != nil {
		name = n.ExecutionClient.ClientType()
	}
	if n.BeaconClient != nil {
		name = fmt.Sprintf("%s/%s", name, n.BeaconClient.ClientName())
	}
	return name
}

func (n *Node) IsRunning() bool {
	return n.ExecutionClient.IsRunning() && n.BeaconClient.IsRunning()
}

// Validator operations
func (n *Node) SignBLSToExecutionChange(
	ctx context.Context,
	blsToExecutionChangeInfo validator.BLSToExecutionChangeInfo,
) (*common.SignedBLSToExecutionChange, error) {
	vc, bn := n.ValidatorClient, n.BeaconClient
	if !vc.ContainsValidatorIndex(blsToExecutionChangeInfo.ValidatorIndex) {
		return nil, fmt.Errorf(
			"validator does not contain specified validator index %d",
			blsToExecutionChangeInfo.ValidatorIndex,
		)
	}
	if domain, err := bn.ComputeDomain(
		ctx,
		common.DOMAIN_BLS_TO_EXECUTION_CHANGE,
		&bn.Config.Spec.GENESIS_FORK_VERSION,
	); err != nil {
		return nil, err
	} else {
		return vc.SignBLSToExecutionChange(domain, blsToExecutionChangeInfo)
	}
}

func (n *Node) SignSubmitBLSToExecutionChanges(
	ctx context.Context,
	blsToExecutionChangesInfo []validator.BLSToExecutionChangeInfo,
) error {
	l := make(common.SignedBLSToExecutionChanges, 0)
	for _, c := range blsToExecutionChangesInfo {
		blsToExecChange, err := n.SignBLSToExecutionChange(
			ctx,
			c,
		)
		if err != nil {
			return err
		}
		l = append(l, *blsToExecChange)
	}

	return n.BeaconClient.SubmitPoolBLSToExecutionChange(ctx, l)
}

func (n *Node) SignVoluntaryExit(
	ctx context.Context,
	epoch common.Epoch,
	validatorIndex common.ValidatorIndex,
) (*phase0.SignedVoluntaryExit, error) {
	vc, bn := n.ValidatorClient, n.BeaconClient
	if !vc.ContainsValidatorIndex(validatorIndex) {
		return nil, fmt.Errorf(
			"validator does not contain specified validator index %d",
			validatorIndex,
		)
	}
	if domain, err := bn.ComputeDomain(
		ctx,
		common.DOMAIN_VOLUNTARY_EXIT,
		nil,
	); err != nil {
		return nil, err
	} else {
		return vc.SignVoluntaryExit(domain, epoch, validatorIndex)
	}
}

func (n *Node) SignSubmitVoluntaryExit(
	ctx context.Context,
	epoch common.Epoch,
	validatorIndex common.ValidatorIndex,
) error {
	exit, err := n.SignVoluntaryExit(ctx, epoch, validatorIndex)
	if err != nil {
		return err
	}
	return n.BeaconClient.SubmitVoluntaryExit(ctx, exit)
}

// Node cluster operations
type Nodes []*Node

// Return all execution clients, even the ones not currently running
func (all Nodes) ExecutionClients() execution.ExecutionClients {
	en := make(execution.ExecutionClients, 0)
	for _, n := range all {
		if n.ExecutionClient != nil {
			en = append(en, n.ExecutionClient)
		}
	}
	return en
}

// Return all proxy pointers, even the ones not currently running
func (all Nodes) Proxies() execution.Proxies {
	ps := make(execution.Proxies, 0)
	for _, n := range all {
		if n.ExecutionClient != nil {
			ps = append(ps, n.ExecutionClient)
		}
	}
	return ps
}

// Return all beacon clients, even the ones not currently running
func (all Nodes) BeaconClients() beacon.BeaconClients {
	bn := make(beacon.BeaconClients, 0)
	for _, n := range all {
		if n.BeaconClient != nil {
			bn = append(bn, n.BeaconClient)
		}
	}
	return bn
}

// Return all validator clients, even the ones not currently running
func (all Nodes) ValidatorClients() validator.ValidatorClients {
	vc := make(validator.ValidatorClients, 0)
	for _, n := range all {
		if n.ValidatorClient != nil {
			vc = append(vc, n.ValidatorClient)
		}
	}
	return vc
}

// Return subset of nodes which are marked as verification nodes
func (all Nodes) VerificationNodes() Nodes {
	// If none is set as verification, then all are verification nodes
	var any bool
	for _, n := range all {
		if n.Verification {
			any = true
			break
		}
	}
	if !any {
		return all
	}

	res := make(Nodes, 0)
	for _, n := range all {
		if n.Verification {
			res = append(res, n)
		}
	}
	return res
}

// Return subset of nodes that are currently running
func (all Nodes) Running() Nodes {
	res := make(Nodes, 0)
	for _, n := range all {
		if n.IsRunning() {
			res = append(res, n)
		}
	}
	return res
}

func (all Nodes) FilterByCL(filters []string) Nodes {
	ret := make(Nodes, 0)
	for _, n := range all {
		for _, filter := range filters {
			if strings.Contains(n.BeaconClient.ClientName(), filter) {
				ret = append(ret, n)
				break
			}
		}
	}
	return ret
}

func (all Nodes) FilterByEL(filters []string) Nodes {
	ret := make(Nodes, 0)
	for _, n := range all {
		for _, filter := range filters {
			if strings.Contains(n.ExecutionClient.ClientType(), filter) {
				ret = append(ret, n)
				break
			}
		}
	}
	return ret
}

func (all Nodes) RemoveNodeAsVerifier(id int) error {
	if id >= len(all) {
		return fmt.Errorf("node %d does not exist", id)
	}
	var any bool
	for _, n := range all {
		if n.Verification {
			any = true
			break
		}
	}
	if any {
		all[id].Verification = false
	} else {
		// If no node is set as verifier, we will set all other nodes as verifiers then
		for i := range all {
			all[i].Verification = (i != id)
		}
	}
	return nil
}

func (all Nodes) ByValidatorIndex(validatorIndex common.ValidatorIndex) *Node {
	for _, n := range all {
		if n.ValidatorClient.ContainsValidatorIndex(validatorIndex) {
			return n
		}
	}
	return nil
}

func (all Nodes) SignSubmitBLSToExecutionChanges(
	ctx context.Context,
	blsToExecutionChanges []validator.BLSToExecutionChangeInfo,
) error {
	// First gather all signed changes
	l := make(common.SignedBLSToExecutionChanges, 0)
	for _, c := range blsToExecutionChanges {
		n := all.ByValidatorIndex(c.ValidatorIndex)
		if n == nil {
			return fmt.Errorf(
				"validator index %d not found",
				c.ValidatorIndex,
			)
		}
		blsToExecChange, err := n.SignBLSToExecutionChange(
			ctx,
			c,
		)
		if err != nil {
			return err
		}
		l = append(l, *blsToExecChange)
	}
	// Then send the signed changes
	for _, n := range all {
		if err := n.BeaconClient.SubmitPoolBLSToExecutionChange(ctx, l); err != nil {
			return err
		}
	}
	return nil
}
