package client

import (
	"context"

	"github.com/spacemeshos/smrepl/common"

	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
)

// GetMeshTransactions returns the transactions on the mesh to or from an address.
func (c *gRPCClient) GetMeshTransactions(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Transaction, uint32, error) {
	ms := c.getMeshServiceClient()
	resp, err := ms.AccountMeshDataQuery(context.Background(), &apitypes.AccountMeshDataQueryRequest{
		Filter: &apitypes.AccountMeshDataFilter{
			AccountId:            &apitypes.AccountId{Address: address.Bytes()},
			AccountMeshDataFlags: uint32(apitypes.AccountMeshDataFlag_ACCOUNT_MESH_DATA_FLAG_TRANSACTIONS),
		},
		MinLayer:   &apitypes.LayerNumber{Number: 0},
		MaxResults: maxResults,
		Offset:     offset,
	})

	if err != nil {
		return nil, 0, err
	}

	txsMap := make(map[string]bool)
	txs := make([]*apitypes.Transaction, 0)

	for _, data := range resp.Data {
		tx := data.GetTransaction()
		if tx != nil {
			if !txsMap[string(tx.Id.Id)] {
				txsMap[string(tx.Id.Id)] = true
				txs = append(txs, tx)
			}
		}
	}
	// hack alert: for now, we return the number of filtered results and not the results returned from the api
	// because they include duplicated transactions in case where a transaction is on more than one mesh block
	// todo: think about default sorting

	return txs, uint32(len(txs)), nil
}

// GetMeshActivations returns activations where the address is the coinbase
func (c *gRPCClient) GetMeshActivations(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Activation, uint32, error) {
	ms := c.getMeshServiceClient()

	resp, err := ms.AccountMeshDataQuery(context.Background(), &apitypes.AccountMeshDataQueryRequest{
		Filter: &apitypes.AccountMeshDataFilter{
			AccountId:            &apitypes.AccountId{Address: address.Bytes()},
			AccountMeshDataFlags: uint32(apitypes.AccountMeshDataFlag_ACCOUNT_MESH_DATA_FLAG_ACTIVATIONS),
		},
		MinLayer:   &apitypes.LayerNumber{Number: 0},
		MaxResults: maxResults,
		Offset:     offset,
	})

	if err != nil {
		return nil, 0, err
	}

	activations := make([]*apitypes.Activation, 0)

	for _, data := range resp.Data {
		a := data.GetActivation()
		if a != nil {
			activations = append(activations, a)
		}
	}

	return activations, resp.TotalResults, nil
}

// GetMeshInfo returns basic mesh data
func (c *gRPCClient) GetMeshInfo() (*common.NetInfo, error) {
	netInfo := &common.NetInfo{}
	ms := c.getMeshServiceClient()

	res, err := ms.GenesisTime(context.Background(), &apitypes.GenesisTimeRequest{})
	if err != nil {
		return nil, err
	}
	netInfo.GenesisTime = res.Unixtime.Value

	currLayer, err := ms.CurrentLayer(context.Background(), &apitypes.CurrentLayerRequest{})
	if err != nil {
		return nil, err
	}
	netInfo.CurrentLayer = currLayer.Layernum.Number

	epochNum, err := ms.CurrentEpoch(context.Background(), &apitypes.CurrentEpochRequest{})
	if err != nil {
		return nil, err
	}
	netInfo.CurrentEpoch = epochNum.Epochnum.Value

	netId, err := ms.NetID(context.Background(), &apitypes.NetIDRequest{})
	if err != nil {
		return nil, err
	}
	netInfo.NetId = netId.Netid.Value

	layersPerEpoch, err := ms.EpochNumLayers(context.Background(), &apitypes.EpochNumLayersRequest{})
	if err != nil {
		return nil, err
	}
	netInfo.LayerPerEpoch = layersPerEpoch.Numlayers.Value

	layerDuration, err := ms.LayerDuration(context.Background(), &apitypes.LayerDurationRequest{})
	if err != nil {
		return nil, err
	}
	netInfo.LayerDuration = layerDuration.Duration.Value

	maxTxsPerSec, err := ms.MaxTransactionsPerSecond(context.Background(), &apitypes.MaxTransactionsPerSecondRequest{})
	if err != nil {
		return nil, err
	}
	netInfo.MaxTxsPerSec = maxTxsPerSec.MaxTxsPerSecond.Value

	return netInfo, nil
}
