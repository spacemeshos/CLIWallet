package repl

const (
	initialTransferMsg    = "Transfer coins from local account to another account."
	destAddressMsg        = "Enter destination address: "
	enterAddressMsg       = "Enter an address: "
	enterRewardsAddress   = "Enter an address for smesher rewards: "
	txIdMsg               = "Enter transaction id: "
	smesherIdMsg          = "Enter Smesher id: "
	amountToTransferMsg   = "Enter amount to transfer in Smidge: "
	confirmTransactionMsg = "Confirm transaction (y/n): "
	confirmDeleteDataMsg  = "Delete the proof of space data file(s)? (y/n)"
	createAccountMsg      = "Account alias (name): "
	useDefaultGasMsg      = "Use default transaction fee of 1 Smidge? (y/n) "
	enterGasPrice         = "Enter transaction fee (Smidge):"

	posDataDirMsg  = "Enter proof of space data directory (relative to node or absolute): "
	posSizeMsg     = "Enter number of units. (%f GiB per unit. Min units: %d, Max units: %d): "
	posProviderMsg = "Enter proof of space compute provider id number: "

	msgSignMsg     = "Enter message to sign (in hex): "
	msgTextSignMsg = "Enter text message to sign: "
	coinUnitName   = "Smidge"
)

const splash = `

                                    .++++++++++++++++++++++++++.
                                    %@@@@@@@@@@@@@@@@@@@@@@@@@@%
                                   -@@@@@@@##############@@@@@@@-
                                     +@@@@@*.          .*@@@@@+
                                      .+@@@@@*.      .*@@@@@+.
                                        .*@@@@@+.  .+@@@@@*.
                                          .*@@@@@++@@@@@*.
                                            .*@@@@@@@@*.
                                              *@@@@@@*
                                            =@@@@@@@@@@=
                                          =@@@@@#::#@@@@@=
                                        =%@@@@%:    :#@@@@%=
                                      -%@@@@%-        -%@@@@%-
                                    -%@@@@%-            -%@@@@%-
                                   *@@@@%-                -%@@@@*
                                   *@@@@#:                :#@@@@*
                                    =@@@@@#:            :#@@@@@=
                                      =@@@@@#:        :#@@@@@=
                                        =@@@@@#:    :#@@@@@=
                                          +@@@@@*..*@@@@@+
                                            +@@@@@@@@@@+
                                             .*@@@@@@*
                                            .*@@@@@@@@*.
                                          .+@@@@@**@@@@@+.
                                         +@@@@@*.  .*@@@@@+
                                       +@@@@@*.      .*@@@@@+
                                     +@@@@@*.          .*@@@@@+
                                   -@@@@@@@##############@@@@@@@-
                                    %@@@@@@@@@@@@@@@@@@@@@@@@@@%
                                    .++++++++++++++++++++++++++.

`
