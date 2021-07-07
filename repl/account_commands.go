package repl

import (
	"encoding/hex"
	"fmt"

	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"github.com/spacemeshos/ed25519"
	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/smrepl/common"
	"github.com/spacemeshos/smrepl/log"
)

func (r *repl) printWalletMnemonic() {
	r.client.PrintWalletMnemonic()
}

func (r *repl) walletInfo() {
	r.client.WalletInfo()
}

// openWallet opens a wallet from locally stored wallet data file
func (r *repl) openWallet() {
	r.clientOpen = r.client.OpenWallet()
	if !r.clientOpen {
		fmt.Println("Wallet NOT opened")
		return
	}
	r.client.WalletInfo()
	r.initializeCommands()
}

// createWallet creates a new wallet
func (r *repl) createWallet() {
	r.clientOpen = r.client.NewWallet()
	if !r.clientOpen {
		fmt.Println("Wallet NOT created")
		return
	}
	r.client.WalletInfo()
	r.initializeCommands()
}

// closeWallet closes an open wallet
func (r *repl) closeWallet() {
	r.client.CloseWallet()
	r.clientOpen = false
	r.initializeCommands()
}

// chooseAccount sets the current account to one of the open wallet's accounts
func (r *repl) chooseAccount() {
	accs, err := r.client.ListAccounts()
	if err != nil {
		log.Error("failure to choose account", err)
		return
	}
	if len(accs) == 0 {
		r.createAccount()
		return
	}

	fmt.Println("Choose an account to load:")
	accNumber := multipleChoice(accs)
	if accNumber == 0 {
		fmt.Println("none selected")
		return
	}
	accNumber = accNumber - 1
	err = r.client.SetCurrentAccount(accNumber)
	if err != nil {
		log.Error("failure to set current account", err)
		return
	}

	account, err := r.client.CurrentAccount()
	if err != nil {
		log.Error("error getting current account", err)
		return
	}

	fmt.Printf("Loaded account alias: `%s`, address: %s\n", account.Name, account.Address().String())

}

// createAccount creates a new account in the currently open wallet
func (r *repl) createAccount() {
	fmt.Println("Create a new account")
	alias := inputNotBlank(createAccountMsg)

	ac, err := r.client.CreateAccount(alias)
	if err != nil {
		log.Error("Failed to create a new account: %v", err)
		return
	}
	err = r.client.StoreAccounts()
	if err != nil {
		log.Error("Failed to save the new account: %v", err)
		return
	}

	fmt.Printf("Created account: %s, address: %s\n", ac.Name, ac.Address().String())
}

// One smesh in base coin units
const onesmh = 1000000000000

// coinAmount formats an amount in base coin units to a display string
func coinAmount(val uint64) string {
	if val >= 1000000000000 {
		return fmt.Sprintf("%d.%012d SMH", val/onesmh, val%onesmh)
	} else if val >= 10000000000 {
		return fmt.Sprintf("0.%012d SMH", val%onesmh)
	} else {
		return fmt.Sprint(val, " Smidge")
	}
}

// printAccountInfo prints current wallet's account info from global state
func (r *repl) printAccountInfo() {
	acc, err := r.getCurrent()
	if err != nil {
		log.Error("failed to get account", err)
		return
	}

	address := gosmtypes.BytesToAddress(acc.PubKey)
	account, err := r.client.AccountState(address)
	if err != nil {
		log.Error("failed to get account info: %v", err)
		return
	}

	fmt.Println("Local alias:", acc.Name)
	printAccount(account, address)
	fmt.Printf("Public key: 0x%s\n", hex.EncodeToString(acc.PubKey))
	fmt.Printf("Private key: 0x%s\n", hex.EncodeToString(acc.PrivKey))
}

// printAccountRewards prints all rewards awarded to the current account
func (r *repl) printLocalAccountRewards() {
	acc, err := r.getCurrent()
	if err != nil {
		log.Error("failed to get account", err)
		return
	}
	r.printRewards(acc.Address())
}

// printAccountState prints the account data member
func printAccount(account *apitypes.Account, address gosmtypes.Address) {
	currBalance := uint64(0)
	if account.StateCurrent.Balance != nil {
		currBalance = account.StateCurrent.Balance.Value
	}

	projectedBalance := uint64(0)
	if account.StateProjected.Balance != nil {
		projectedBalance = account.StateProjected.Balance.Value
	}

	fmt.Println("Address:", address.String())
	fmt.Println("Balance:", coinAmount(currBalance)) // currBalance, coinUnitName)
	fmt.Println("Nonce:", account.StateCurrent.Counter)
	fmt.Println("Projected balance:", coinAmount(projectedBalance)) // projectedBalance, coinUnitName)
	fmt.Println("Projected nonce:", account.StateProjected.Counter)
	fmt.Println("Projected state includes all pending transactions that haven't been added to the mesh yet.")
}

// printReward prints a Reward
func printReward(r *apitypes.Reward) {
	fmt.Println("Rewarded on layer:", r.Layer.Number)
	//fmt.Println("Rewarded for layer:", r.LayerComputed.Number)
	fmt.Println("Layer reward", r.LayerReward.Value, coinUnitName)
	fmt.Println("Transaction fees", r.Total.Value-r.LayerReward.Value, coinUnitName)
	fmt.Println("Total reward", r.Total.Value, coinUnitName)
	//fmt.Println("Smesher id", "0x"+hex.EncodeToString(r.Smesher.Id))
	fmt.Println("Rewards account:", gosmtypes.BytesToAddress(r.Coinbase.Address).String())
}

// getCurrent returns the current open wallet's account. If there is no current account
// then it prompts the user to choose one of the wallet's accounts.
func (r *repl) getCurrent() (acc *common.LocalAccount, err error) {
	acc, err = r.client.CurrentAccount()
	if err != nil {
		r.chooseAccount()
		acc, err = r.client.CurrentAccount()
	}
	return
}

// sign signs a hex string with the current account
func (r *repl) sign() {
	acc, err := r.getCurrent()
	if err != nil {
		log.Error("failed to get account", err)
		return
	}

	msgStr := inputNotBlank(msgSignMsg)
	msg, err := hex.DecodeString(msgStr)
	if err != nil {
		log.Error("failed to decode msg hex string: %v", err)
		return
	}
	signature := ed25519.Sign2(acc.PrivKey, msg)
	fmt.Printf("signature (in hex): %x\n", signature)
}

// signText signs a string with the current account
func (r *repl) signText() {
	acc, err := r.getCurrent()
	if err != nil {
		log.Error("failed to get account", err)
		return
	}
	msg := inputNotBlank(msgTextSignMsg)
	signature := ed25519.Sign2(acc.PrivKey, []byte(msg))
	fmt.Printf("signature (in hex): %x\n", signature)
}
