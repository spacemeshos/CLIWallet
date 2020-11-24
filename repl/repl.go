package repl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spacemeshos/CLIWallet/common"
	"github.com/spacemeshos/CLIWallet/log"
	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"github.com/spacemeshos/ed25519"
	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
	"google.golang.org/genproto/googleapis/rpc/status"

	"github.com/c-bata/go-prompt"
)

const (
	prefix           = "$ "
	printPrefix      = ">"
	minVerifiedLayer = 575
)

// TestMode variable used for check if unit test is running
var TestMode = false

type command struct {
	text        string
	description string
	fn          func()
}

type repl struct {
	commands []command
	client   Client
	input    string
}

// Client interface to REPL clients.
type Client interface {

	// Local account management methods
	CreateAccount(alias string) *common.LocalAccount
	CurrentAccount() *common.LocalAccount
	SetCurrentAccount(a *common.LocalAccount)
	ListAccounts() []string
	GetAccount(name string) (*common.LocalAccount, error)
	StoreAccounts() error

	// Local config
	ServerUrl() string

	// Node service
	NodeStatus() (*apitypes.NodeStatus, error)
	NodeInfo() (*common.NodeInfo, error)
	Sanity() error

	// Mesh service
	GetMeshTransactions(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Transaction, uint32, error)
	GetMeshActivations(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Activation, uint32, error)

	// Transaction service
	Transfer(recipient gosmtypes.Address, nonce, amount, gasPrice, gasLimit uint64, key ed25519.PrivateKey) (*apitypes.TransactionState, error)
	TransactionState(txId []byte, includeTx bool) (*apitypes.TransactionState, *apitypes.Transaction, error)

	// Smesher service
	GetSmesherId() ([]byte, error)
	IsSmeshing() (bool, error)
	StartSmeshing(address gosmtypes.Address, dataDir string, dataSizeBytes uint64) (*status.Status, error)
	StopSmeshing(deleteFiles bool) (*status.Status, error)
	GetCoinbase() (*gosmtypes.Address, error)
	SetCoinbase(coinbase gosmtypes.Address) (*status.Status, error)

	// debug service
	DebugAllAccounts() ([]*apitypes.Account, error)

	// global state service
	AccountState(address gosmtypes.Address) (*common.AccountState, error)
	AccountRewards(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Reward, uint32, error)
	AccountTransactionsReceipts(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.TransactionReceipt, uint32, error)
	GlobalStateHash() (*apitypes.GlobalStateHash, error)
	SmesherRewards(smesherId []byte, offset uint32, maxResults uint32) ([]*apitypes.Reward, uint32, error)

	//Unlock(passphrase string) error
	//IsAccountUnLock(id string) bool
	//Lock(passphrase string) error
	//SetVariables(params, flags []string) error
	//GetVariable(key string) string
	//Restart(params, flags []string) error
	//NeedRestartNode(params, flags []string) bool
	//Setup(allocation string) error
}

func (r *repl) initializeCommands() {
	r.commands = []command{
		// account commands
		{"new", "Create a new account (key pair) and set as current", r.createAccount},
		{"set", "Set one of the previously created accounts as current", r.chooseAccount},
		{"info", "Display the current account info", r.accountInfo},

		// transactions
		{"send-coin", "Transfer coins from current account to another account", r.submitCoinTransaction},
		{"tx-status", "Display a transaction status", r.printTransactionStatus},
		{"txs", "Print all outgoing and incoming transactions for the current account that are on the mesh", r.printAccountTransactions},
		{"rewards", "Print all rewards awarded to the current account", r.printAccountRewards},

		{"sign", "Sign a hex message with the current account private key", r.sign},
		{"text-sign", "Sign a text message with the current account private key", r.textsign},

		// printing status and state of things
		{"accounts", "Display all mesh accounts (debug)", r.printAllAccounts},
		{"node", "Display node status", r.nodeInfo},
		{"state", "Display the current global state", r.printGlobalState},

		// smeshing operations
		{"get-rewards-account", "Get current account as the node smesher's rewards account", r.printCoinbase},
		{"set-rewards-account", "Set current account as the node smesher's rewards account", r.setCoinbase},
		{"get-smesher-id", "Get the node smesher's current rewards account", r.printSmesherId},
		{"set-rewards-account", "Set current account as the node smesher's rewards account", r.setCoinbase},
		{"smeshing-status", "Set current account as the node smesher's rewards account", r.printSmeshingStatus},
		{"start-smeshing", "Start smeshing using the current account as the rewards account", r.startSmeshing},
		{"stop-smeshing", "Stop smeshing", r.stopSmeshing},
		{"smesher-rewards", "Print rewards for a smesher", r.printSmesherRewards},

		//{"smesh", "Start smeshing", r.smesh},

		{"quit", "Quit the CLI", r.quit},

		//{"unlock accountInfo", "Unlock accountInfo.", r.unlockAccount},
		//{"lock accountInfo", "Lock LocalAccount.", r.lockAccount},
		//{"setup", "Setup POST.", r.setup},
		//{"restart node", "Restart node.", r.restartNode},
		//{"set", "change CLI flag or param. E.g. set param a=5 flag c=5 or E.g. set param a=5", r.setCLIFlagOrParam},
		//{"echo", "Echo runtime variable.", r.echoVariable},
	}
}

// Start starts REPL.
func Start(c Client) {
	if !TestMode {
		r := &repl{client: c}
		r.initializeCommands()

		runPrompt(r.executor, r.completer, r.firstTime, uint16(len(r.commands)))
	} else {
		// holds for unit test purposes
		hold := make(chan bool)
		<-hold
	}
}

func (r *repl) executor(text string) {
	for _, c := range r.commands {
		if len(text) >= len(c.text) && text[:len(c.text)] == c.text {
			r.input = text
			//log.Debug(userExecutingCommandMsg, c.text)
			c.fn()
			return
		}
	}

	fmt.Println(printPrefix, "invalid command.")
}

func (r *repl) completer(in prompt.Document) []prompt.Suggest {
	suggets := make([]prompt.Suggest, 0)
	for _, command := range r.commands {
		s := prompt.Suggest{
			Text:        command.text,
			Description: command.description,
		}

		suggets = append(suggets, s)
	}

	return prompt.FilterHasPrefix(suggets, in.GetWordBeforeCursor(), true)
}

func (r *repl) commandLineParams(idx int, input string) string {
	c := r.commands[idx]
	params := strings.Replace(input, c.text, "", -1)

	return strings.TrimSpace(params)
}

func (r *repl) firstTime() {
	fmt.Print(printPrefix, splash)

	if err := r.client.Sanity(); err != nil {
		log.Error("Failed to connect to node at %v: %v", r.client.ServerUrl(), err)
		r.quit()
	}

	fmt.Println("Welcome to Spacemesh. Connected to node at", r.client.ServerUrl())
}

func (r *repl) quit() {
	os.Exit(0)
}
