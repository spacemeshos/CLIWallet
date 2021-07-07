package repl

import (
	"fmt"
	"os"
	"strings"

	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"github.com/spacemeshos/ed25519"
	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/smrepl/common"
	"github.com/spacemeshos/smrepl/log"
	"google.golang.org/genproto/googleapis/rpc/status"

	"github.com/c-bata/go-prompt"
)

const (
	prefix = "$ "
)

const (
	commandStateRoot = iota
	commandStateWallet
	commandStateAccount
	commandStateStatus
	commandStateState
	commandStateMesh
	commandStatePost
	commandStateSmesher
	commandStateDBG
	commandStateLeaf
)

// TestMode variable used for check if unit test is running
var TestMode = false

type command struct {
	parent      int
	text        string
	state       int
	description string
	fn          func()
}

type repl struct {
	commands   []command
	client     Client
	clientOpen bool
	input      string
}

// Client interface to REPL clients.
type Client interface {
	PrintWalletMnemonic()
	WalletInfo()
	IsOpen() bool
	OpenWallet() bool
	NewWallet() bool
	CloseWallet()

	// Local account management methods
	CreateAccount(alias string) (*common.LocalAccount, error)
	CurrentAccount() (*common.LocalAccount, error)
	SetCurrentAccount(accountNumber int) error
	ListAccounts() ([]string, error)
	GetAccount(name string) (*common.LocalAccount, error)
	StoreAccounts() error

	// Local config
	ServerInfo() string

	// Node service
	NodeStatus() (*apitypes.NodeStatus, error)
	NodeInfo() (*common.NodeInfo, error)
	Echo() error

	// Mesh service
	GetMeshTransactions(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Transaction, uint32, error)
	GetMeshActivations(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Activation, uint32, error)
	GetMeshInfo() (*common.NetInfo, error)

	// Transaction service
	Transfer(recipient gosmtypes.Address, nonce, amount, gasPrice, gasLimit uint64, key ed25519.PrivateKey) (*apitypes.TransactionState, error)
	TransactionState(txId []byte, includeTx bool) (*apitypes.TransactionState, *apitypes.Transaction, error)

	// Smesher service
	IsSmeshing() (*apitypes.IsSmeshingResponse, error)
	StartSmeshing(request *apitypes.StartSmeshingRequest) (*status.Status, error)
	StopSmeshing(deleteFiles bool) (*status.Status, error)
	GetPostComputeProviders(benchmark bool) ([]*apitypes.PoSTSetupComputeProvider, error)
	GetSmesherId() ([]byte, error)
	GetRewardsAddress() (*gosmtypes.Address, error)
	SetRewardsAddress(coinbase gosmtypes.Address) (*status.Status, error)
	Config() (*apitypes.PoSTConfigResponse, error)
	PostStatus() (*apitypes.PoSTSetupStatusResponse, error)
	PostDataCreationProgressStream() (apitypes.SmesherService_PoSTSetupStatusStreamClient, error)

	// debug service
	DebugAllAccounts() ([]*apitypes.Account, error)

	// global state service
	AccountState(address gosmtypes.Address) (*apitypes.Account, error)
	AccountRewards(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Reward, uint32, error)
	AccountRewardsStream(address gosmtypes.Address) (apitypes.GlobalStateService_AccountDataStreamClient, error)
	AccountUpdatesStream(address gosmtypes.Address) (apitypes.GlobalStateService_AccountDataStreamClient, error)
	AccountTransactionsReceipts(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.TransactionReceipt, uint32, error)
	GlobalStateHash() (*apitypes.GlobalStateHash, error)
	SmesherRewards(smesherId []byte, offset uint32, maxResults uint32) ([]*apitypes.Reward, uint32, error)
}

func (r *repl) initializeCommands() {
	firstStageCommands := []command{
		{commandStateRoot, "wallet", commandStateWallet, "Wallet related commands", nil},
		{commandStateRoot, "state", commandStateState, "Global state commands", nil},
		{commandStateRoot, "status", commandStateStatus, "Status commands", nil},
		{commandStateRoot, "mesh", commandStateMesh, "Mesh data", nil},
		{commandStateRoot, "smesher", commandStateSmesher, "Smesher commands", nil},
		{commandStateRoot, "post", commandStatePost, "Proof of spacetime commands", nil},
		{commandStateRoot, "dbg", commandStateDBG, "Debugging commands", nil},
		{commandStateRoot, "quit", commandStateLeaf, "Quit app", r.quit},
	}
	accountCommands := []command{
		// wallets
		{commandStateWallet, "open", commandStateLeaf, "Open a wallet", r.openWallet},
		{commandStateWallet, "create", commandStateLeaf, "Create a wallet", r.createWallet},
	}
	if r.clientOpen {
		firstStageCommands = append(firstStageCommands,
			command{commandStateRoot, "account", commandStateAccount, "Wallet's accounts commands", nil})

		accountCommands = []command{
			// local wallet account commands
			{commandStateWallet, "info", commandStateLeaf, "Display wallet info", r.walletInfo},
			{commandStateWallet, "mnemonic", commandStateLeaf, "Display wallet mnemonic", r.printWalletMnemonic},

			{commandStateWallet, "close", commandStateLeaf, "Close current wallet", r.closeWallet},

			{commandStateAccount, "new", commandStateLeaf, "Create a new account (key pair) and set as current", r.createAccount},
			{commandStateAccount, "set", commandStateLeaf, "Set one of the previously created accounts as current", r.chooseAccount},
			{commandStateAccount, "info", commandStateLeaf, "Display the current account info", r.printAccountInfo},
			{commandStateAccount, "rewards", commandStateLeaf, "Display all rewards awarded to the current account", r.printLocalAccountRewards},
			{commandStateAccount, "sign", commandStateLeaf, "Sign a hex message with the current account private key", r.sign},
			{commandStateAccount, "text-sign", commandStateLeaf, "Sign a text message with the current account private key", r.signText},
			{commandStateAccount, "txs", commandStateLeaf, "Display all outgoing and incoming transactions for the current account that are on the mesh", r.printCurrAccountMeshTransactions},
			{commandStateAccount, "send-coin", commandStateLeaf, "Transfer coins from current account to another account", r.submitCoinTransaction},
		}
	}

	otherCommands := []command{
		// Misc entities status
		{commandStateStatus, "node", commandStateLeaf, "Display node status", r.nodeInfo},
		{commandStateStatus, "net", commandStateLeaf, "Display network information", r.printMeshInfo},
		{commandStateStatus, "tx", commandStateLeaf, "Display a transaction status", r.printTransactionStatus},

		// global state
		{commandStateState, "account", commandStateLeaf, "Display an account balance and nonce", r.printAccountState},
		{commandStateState, "account-txs", commandStateLeaf, "Display account transactions in global state", r.printAccountState},

		{commandStateState, "rewards", commandStateLeaf, "Display an account rewards ", r.printAccountRewards},

		// global state streams
		{commandStateState, "stream-rewards", commandStateLeaf, "Stream new rewards for an account", r.printAccountRewardsStream},
		{commandStateState, "stream-account", commandStateLeaf, "Stream account updates", r.printAccountUpdatesStream},

		{commandStateState, "smesher-rewards", commandStateLeaf, "Display smesher rewards", r.printSmesherRewards},
		{commandStateState, "global", commandStateLeaf, "Display the most recent network global state", r.printGlobalState},

		// mesh
		{commandStateMesh, "transactions", commandStateLeaf, "Display mesh transaction for an account", r.printMeshTransactions},

		// smeshing - smesher ops

		{commandStateSmesher, "id", commandStateLeaf, "Display current smesher id", r.printSmesherId},
		{commandStateSmesher, "rewards-address", commandStateLeaf, "Display current smesher rewards address", r.printRewardsAddress},
		{commandStateSmesher, "set-rewards-address", commandStateLeaf, "Set the smesher rewards address", r.setRewardsAddress},
		{commandStateSmesher, "rewards", commandStateLeaf, "Display current smesher rewards", r.printCurrentSmesherRewards},
		{commandStateSmesher, "status", commandStateLeaf, "Display smesher status", r.printSmeshingStatus},
		{commandStateSmesher, "stop", commandStateLeaf, "Stop smeshing", r.stopSmeshing},

		{commandStatePost, "status", commandStateLeaf, "Display the proof of space status", r.printPostStatus},
		{commandStatePost, "providers", commandStateLeaf, "Display the available proof of space providers", r.printPosProviders},
		{commandStatePost, "setup", commandStateLeaf, "Set up (or change) smesher proof of space data", r.setupPos},

		{commandStatePost, "progress", commandStateLeaf, "Stream proof of space data creation progress", r.printPostDataCreationProgress},

		// debug commands
		{commandStateDBG, "all-accounts", commandStateLeaf, "Display all global state accounts", r.printAllAccounts},
	}
	r.commands = append(firstStageCommands, append(accountCommands, otherCommands...)...)
}

// Start starts the REPL
func Start(c Client) {

	// init logging system
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Aborting. Can't get current dir. Your system is high.", err)
		return
	}

	log.InitSpacemeshLoggingSystem(path, "log.txt")

	log.Info("new session started")

	if !TestMode {
		r := &repl{client: c}
		r.clientOpen = c.IsOpen()
		r.initializeCommands()
		runPrompt(r.executor, r.completer, r.firstTime, uint16(len(r.commands)))
	} else {
		// holds for unit test purposes
		hold := make(chan bool)
		<-hold
	}
}

func (r *repl) executor(text string) {
	// All commands currently follows a format of `FirstStageCommand SecondStageCommand ...`
	textSlice := strings.Split(text, " ")
	parseState := commandStateRoot
	for _, s := range textSlice {
		for _, c := range r.commands {
			if parseState == c.parent && s == c.text {
				if c.state == commandStateLeaf {
					r.input = text
					//log.Debug(userExecutingCommandMsg, c.text)
					c.fn()
					return
				} else {
					parseState = c.state
				}
			}
		}
	}

	fmt.Println("invalid command.")
}

func (r *repl) completer(in prompt.Document) []prompt.Suggest {
	suggests := make([]prompt.Suggest, 0)
	textSliceBeforeCursor := strings.Split(in.TextBeforeCursor(), " ")
	parseState := commandStateRoot
	// get current state by parsing the whole text before the cursor position
	for _, s := range textSliceBeforeCursor {
		for _, c := range r.commands {
			if parseState == c.parent && s == c.text {
				parseState = c.state
			}
		}
	}

	// give suggestions based on current state
	for _, command := range r.commands {
		if command.parent != parseState {
			continue
		}

		s := prompt.Suggest{
			Text:        command.text,
			Description: command.description,
		}

		suggests = append(suggests, s)
	}

	return prompt.FilterHasPrefix(suggests, in.GetWordBeforeCursor(), true)
}

func (r *repl) firstTime() {
	fmt.Print(splash)

	_, err := r.client.GetMeshInfo()
	if err != nil {
		log.Error("Failed to connect to mesh service at %v: %v", r.client.ServerInfo(), err)
		r.quit()
	}

	fmt.Println("Welcome to Spacemesh. Connected to api server at", r.client.ServerInfo())
	r.printMeshInfo()
}

func (r *repl) quit() {
	os.Exit(0)
}
