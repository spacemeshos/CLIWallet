package client

import (
	"context"

	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
)

// GlobalStateHash returns the current global state hash
func (c *gRPCClient) GlobalStateHash() (*apitypes.GlobalStateHash, error) {
	gsc := c.getGlobalStateServiceClient()
	if resp, err := gsc.GlobalStateHash(context.Background(), &apitypes.GlobalStateHashRequest{}); err != nil {
		return nil, err
	} else {
		return resp.Response, nil
	}
}

// AccountState returns basic account data such as balance and nonce from the global state
func (c *gRPCClient) AccountState(address gosmtypes.Address) (*apitypes.Account, error) {
	gsc := c.getGlobalStateServiceClient()
	resp, err := gsc.Account(context.Background(), &apitypes.AccountRequest{
		AccountId: &apitypes.AccountId{Address: address.Bytes()}})
	if err != nil {
		return nil, err
	}

	return resp.AccountWrapper, nil
}

// SmesherRewards returns rewards for a smesher identified by a smesher id
func (c *gRPCClient) SmesherRewards(smesherId []byte, offset uint32, maxResults uint32) ([]*apitypes.Reward, uint32, error) {
	gsc := c.getGlobalStateServiceClient()
	resp, err := gsc.SmesherDataQuery(context.Background(), &apitypes.SmesherDataQueryRequest{
		SmesherId:  &apitypes.SmesherId{Id: smesherId},
		MaxResults: maxResults,
		Offset:     offset,
	})

	if err != nil {
		return nil, 0, err
	}

	return resp.Rewards, resp.TotalResults, nil
}

// AccountRewards returns rewards for an account
func (c *gRPCClient) AccountRewards(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Reward, uint32, error) {
	gsc := c.getGlobalStateServiceClient()
	resp, err := gsc.AccountDataQuery(context.Background(), &apitypes.AccountDataQueryRequest{
		Filter: &apitypes.AccountDataFilter{
			AccountId:        &apitypes.AccountId{Address: address.Bytes()},
			AccountDataFlags: uint32(apitypes.AccountDataFlag_ACCOUNT_DATA_FLAG_REWARD),
		},

		MaxResults: maxResults,
		Offset:     offset,
	})

	if err != nil {
		return nil, 0, err
	}

	rewards := make([]*apitypes.Reward, 0)

	for _, data := range resp.AccountItem {
		r := data.GetReward()
		if r != nil {
			rewards = append(rewards, r)
		}
	}

	return rewards, resp.TotalResults, nil
}

// getAccountStream returns an AccountDataStreamClient
func (c *gRPCClient) getAccountStream(address gosmtypes.Address, flag apitypes.AccountDataFlag) (apitypes.GlobalStateService_AccountDataStreamClient, error) {
	gsc := c.getGlobalStateServiceClient()
	return gsc.AccountDataStream(context.Background(), &apitypes.AccountDataStreamRequest{
		Filter: &apitypes.AccountDataFilter{
			AccountId: &apitypes.AccountId{
				Address: address.Bytes()},
			AccountDataFlags: uint32(flag),
		},
	})
}

// AccountRewardsStream returns a stream of account rewards
func (c *gRPCClient) AccountRewardsStream(address gosmtypes.Address) (apitypes.GlobalStateService_AccountDataStreamClient, error) {
	return c.getAccountStream(address, apitypes.AccountDataFlag_ACCOUNT_DATA_FLAG_REWARD)
}

// AccountUpdatesStream returns a stream of account changes
func (c *gRPCClient) AccountUpdatesStream(address gosmtypes.Address) (apitypes.GlobalStateService_AccountDataStreamClient, error) {
	return c.getAccountStream(address, apitypes.AccountDataFlag_ACCOUNT_DATA_FLAG_ACCOUNT)
}

// AccountTransactionsReceipts returns transaction receipts for an account
func (c *gRPCClient) AccountTransactionsReceipts(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.TransactionReceipt, uint32, error) {
	gsc := c.getGlobalStateServiceClient()

	resp, err := gsc.AccountDataQuery(context.Background(), &apitypes.AccountDataQueryRequest{
		Filter: &apitypes.AccountDataFilter{
			AccountId:        &apitypes.AccountId{Address: address.Bytes()},
			AccountDataFlags: uint32(apitypes.AccountDataFlag_ACCOUNT_DATA_FLAG_TRANSACTION_RECEIPT),
		},

		MaxResults: maxResults,
		Offset:     offset,
	})

	if err != nil {
		return nil, 0, err
	}

	receipts := make([]*apitypes.TransactionReceipt, 0)

	for _, data := range resp.AccountItem {
		r := data.GetReceipt()
		if r != nil {
			receipts = append(receipts, r)
		}
	}

	return receipts, resp.TotalResults, nil
}
