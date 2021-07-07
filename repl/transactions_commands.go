package repl

import (
	"encoding/hex"
	"fmt"
	"strconv"

	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"github.com/spacemeshos/go-spacemesh/common/util"

	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/smrepl/log"
)

var transactionStateDisStringsMap = map[int32]string{
	0: "Unspecified state",
	1: "Rejected",
	2: "Insufficient funds",
	3: "Conflicting",
	4: "Submitted to the network",
	5: "On the mesh but not yet processed",
	6: "Processed",
}

// Print a transaction status
func (r *repl) printTransactionStatus() {
	txIdStr := inputNotBlank(txIdMsg)
	txId := util.FromHex(txIdStr)
	txState, tx, err := r.client.TransactionState(txId, true)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if txState != nil {
		txStateDispString := transactionStateDisStringsMap[int32(txState.State.Number())]
		fmt.Println("State:", txStateDispString)
	} else {
		fmt.Println("Unknown transaction state")
	}

	if tx != nil {
		printTransaction(tx)
	} else {
		fmt.Println("Unknown transaction")
	}
}

// canSubmitTransactions returns true if the node is accepting transactions.
// todo: this should move to a method in the transactions service.
func (r *repl) canSubmitTransactions() bool {

	status, err := r.client.NodeStatus()
	if err != nil {
		log.Error("failed to get node status: %v", err)
		return false
	}

	// for now, we allow to submit txs if the node is synced
	return status.IsSynced //&& status.TopLayer.Number > minVerifiedLayer
}

func (r *repl) submitCoinTransaction() {

	if !r.canSubmitTransactions() {
		fmt.Println("Can't submit a new transaction. Please try again later")
		return
	}
	fmt.Println(initialTransferMsg)
	acc, err := r.getCurrent()
	if err != nil {
		log.Error("failed to get account", err)
		return
	}

	srcAddress := gosmtypes.BytesToAddress(acc.PubKey)
	acctState, err := r.client.AccountState(srcAddress)
	if err != nil {
		log.Error("failed to get account info: %v", err)
		return
	}

	destAddressStr := inputNotBlank(destAddressMsg)
	destAddress := gosmtypes.HexToAddress(destAddressStr)

	amountStr := inputNotBlank(amountToTransferMsg)

	gas := uint64(1)
	if yesOrNoQuestion(useDefaultGasMsg) == "n" {
		gasStr := inputNotBlank(enterGasPrice)
		gas, err = strconv.ParseUint(gasStr, 10, 64)
		if err != nil {
			log.Error("invalid transaction fee", err)
			return
		}
	}

	fmt.Println("New transaction summary:")
	fmt.Println("From:  ", srcAddress.String())
	fmt.Println("To:    ", destAddress.String())
	fmt.Println("Amount:", amountStr, coinUnitName)
	fmt.Println("Fee:   ", gas, coinUnitName)
	fmt.Println("Nonce: ", acctState.StateProjected.Counter)

	amount, _ := strconv.ParseUint(amountStr, 10, 64)
	// todo: handle error here!

	if yesOrNoQuestion(confirmTransactionMsg) == "y" {
		txState, err := r.client.Transfer(destAddress, acctState.StateProjected.Counter, amount, gas, 100, acc.PrivKey)
		if err != nil {
			log.Error(err.Error())
			return
		}

		txStateDispString := transactionStateDisStringsMap[int32(txState.State.Number())]

		fmt.Println("Transaction submitted.")
		fmt.Printf("Transaction id: 0x%v\n", hex.EncodeToString(txState.Id.Id))
		fmt.Println("Transaction state:", txStateDispString)
	}
}

// helper method - prints tx info
func printTransaction(t *apitypes.Transaction) {

	txIdStr := "0x" + util.Bytes2Hex(t.Id.Id)
	fmt.Println("Transaction summary:")
	fmt.Printf("Id: %v\n", txIdStr)
	fmt.Println("From:", gosmtypes.BytesToAddress(t.Sender.Address).String())

	ct := t.GetCoinTransfer()
	if ct != nil {
		fmt.Println("To (coin account):", gosmtypes.BytesToAddress(ct.Receiver.Address).String())
		fmt.Println("Nonce:", t.Counter)
		fmt.Println("Amount:", t.Amount.Value, coinUnitName)
		fmt.Println("Fee:", t.GasOffered.GasProvided, coinUnitName)
		return
	}

	sct := t.GetSmartContract()
	if sct == nil {
		log.Error("expected a smart contract transaction type")
		return
	}

	// todo: printout smart contract transaction data here
}
