package main

import (
	"flag"
	"fmt"
	"github.com/spacemeshos/smrepl/client"
	"github.com/spacemeshos/smrepl/log"
	"github.com/spacemeshos/smrepl/repl"
	"os"
)

func main() {

	var (
		dataDir    string
		walletName string
		be         *client.WalletBackend
	)
	grpcServer := client.DefaultGRPCServer
	secureConnection := client.DefaultSecureConnection

	flag.StringVar(&grpcServer, "server", grpcServer, fmt.Sprintf("The Spacemesh api grpc server host and port. Defaults to %s", client.DefaultGRPCServer))
	flag.BoolVar(&secureConnection, "secure", secureConnection, "Connect securely to the server. Default is false")
	flag.StringVar(&dataDir, "wallet_directory", getwd(), "set default wallet files directory")
	flag.StringVar(&walletName, "wallet", "", "set the name of wallet file to open")

	flag.Parse()

	be, err := client.OpenConnection(grpcServer, secureConnection, dataDir)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}
	if walletName != "" {
		walletPath := dataDir + "/" + walletName
		fmt.Println("opening ", walletPath)
		be, err = client.OpenWalletBackend(walletPath, grpcServer, secureConnection)
		if err != nil {
			fmt.Println("failed to open wallet file : ", err)
			os.Exit(1)
		}
	}

	_, err = be.GetMeshInfo()
	if err != nil {
		log.Error("Failed to connect to mesh service at %v: %v", be.ServerInfo(), err)
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	repl.Start(be)
}

func getwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}
