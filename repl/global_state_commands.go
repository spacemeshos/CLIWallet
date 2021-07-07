package repl

import (
	"encoding/hex"
	"fmt"
	"io"

	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/common/util"

	"github.com/spacemeshos/smrepl/log"
)

// printRewards prints all rewards awarded to an account
func (r *repl) printRewards(address gosmtypes.Address) {
	// todo: request offset and total from user
	rewards, total, err := r.client.AccountRewards(address, 0, 0)
	if err != nil {
		log.Error("failed to get rewards: %v", err)
		return
	}

	fmt.Printf("Total rewards: %d\n", total)
	for i, r := range rewards {
		if i != 0 {
			fmt.Println("-----")
		}
		printReward(r)
	}
}

// printAccountRewards prints all rewards awarded to an account
func (r *repl) printAccountRewards() {
	addrStr := inputNotBlank(enterAddressMsg)
	addr := gosmtypes.HexToAddress(addrStr)
	r.printRewards(addr)
}

// printAccountRewardsStream prints new rewards awarded to an account
func (r *repl) printAccountRewardsStream() {
	addrStr := inputNotBlank(enterAddressMsg)
	addr := gosmtypes.HexToAddress(addrStr)
	streamClient, err := r.client.AccountRewardsStream(addr)
	if err != nil {
		log.Error("failed to get rewards stream for account: %v", err)
		return
	}

	fmt.Println("Listening to new rewards for address: ", addr.String())

	done := make(chan bool)
	go func() {
		for {
			resp, err := streamClient.Recv()
			if err == io.EOF {
				// server closed the stream
				log.Info("api server closed the server-side stream")
				done <- true
			} else if err != nil {
				log.Error("error reading from rewards stream: %v", err)
				done <- true
			}

			reward := resp.GetDatum().GetReward()
			printReward(reward)
		}
	}()
}

// printAccountRewardsStream prints account state updates
func (r *repl) printAccountUpdatesStream() {
	addrStr := inputNotBlank(enterAddressMsg)
	address := gosmtypes.HexToAddress(addrStr)
	streamClient, err := r.client.AccountRewardsStream(address)
	if err != nil {
		log.Error("failed to get updates stream for account: %v", err)
		return
	}

	fmt.Println("Listening for new updates for address: ", address.String())

	done := make(chan bool)
	go func() {
		for {
			resp, err := streamClient.Recv()
			if err == io.EOF {
				// server closed the stream
				log.Info("api server closed the server-side stream")
				done <- true
			} else if err != nil {
				log.Error("error reading from stream: %v", err)
				done <- true
			}

			account := resp.GetDatum().GetAccountWrapper()
			printAccount(account, address)
		}
	}()
}

// printGlobalState prints the current global state
func (r *repl) printGlobalState() {
	resp, err := r.client.GlobalStateHash()
	if err != nil {
		log.Error("failed to get global state: %v", err)
		return
	}

	fmt.Println("Hash:", "0x"+hex.EncodeToString(resp.RootHash))
	fmt.Println("Layer:", resp.Layer.Number)
}

// printAccountState prints an account's global state
func (r *repl) printAccountState() {
	addressStr := inputNotBlank(enterAddressMsg)
	address := gosmtypes.BytesToAddress(util.FromHex(addressStr))
	account, err := r.client.AccountState(address)
	if err != nil {
		log.Error("failed to get account info: %v", err)
		return
	}

	printAccount(account, address)
}
