package repl

import (
	"fmt"
	"time"

	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"

	"github.com/spacemeshos/smrepl/log"
)

func (r *repl) printMeshInfo() {
	info, err := r.client.GetMeshInfo()
	if err != nil {
		log.Error("failed to get mesh info: %v", err)
		return
	}

	localGenesisTime := time.Unix(int64(info.GenesisTime), 0)

	fmt.Println("Network info:")
	fmt.Println("Network id:", info.NetId)
	fmt.Println("Max transactions per second:", info.MaxTxsPerSec)
	fmt.Println("Layers per epoch:", info.LayerPerEpoch)
	fmt.Printf("Layer duration: %d seconds\n", info.LayerDuration)
	fmt.Println("Current layer:", info.CurrentLayer)
	fmt.Println("Current epoch:", info.CurrentEpoch)
	fmt.Println("Genesis time:", localGenesisTime.Local().String())
}

// printCurrAccountMeshTransactions displays mesh transactions for the current account
func (r *repl) printCurrAccountMeshTransactions() {
	acc, err := r.getCurrent()
	if err != nil {
		log.Error("failed to get account", err)
		return
	}
	r.printAccountMeshTransactions(acc.Address())
}

// printAccountMeshTransactions displays mesh transactions for an account
func (r *repl) printMeshTransactions() {
	addrStr := inputNotBlank(enterAddressMsg)
	addr := gosmtypes.HexToAddress(addrStr)
	r.printAccountMeshTransactions(addr)
}

// Print transaction for an account from mesh data
func (r *repl) printAccountMeshTransactions(address gosmtypes.Address) {

	// todo: request offset and total from user
	txs, total, err := r.client.GetMeshTransactions(address, 0, 1000)
	if err != nil {
		log.Error("failed to print transactions: %v", err)
		return
	}

	fmt.Printf("Total mesh transactions: %d", total)
	for i, tx := range txs {
		if i != 0 {
			fmt.Println("-----")
		}
		printTransaction(tx)
	}
}
