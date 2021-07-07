package repl

import (
	"fmt"

	"github.com/spacemeshos/smrepl/log"
)

/// nodeInfo displays basic node info
func (r *repl) nodeInfo() {
	info, err := r.client.NodeInfo()
	if err != nil {
		log.Error("failed to get node info: %v", err)
		return
	}
	fmt.Println("Node info:")
	fmt.Println("Version:", info.Version)
	fmt.Println("Build:", info.Build)
	fmt.Println("API server:", r.client.ServerInfo())

	status, err := r.client.NodeStatus()
	if err != nil {
		log.Error("failed to get node status: %v", err)
		return
	}

	fmt.Println("Synced:", status.IsSynced)
	fmt.Println("Synced layer:", status.SyncedLayer.Number)
	fmt.Println("Current layer:", status.TopLayer.Number)
	fmt.Println("Verified layer:", status.VerifiedLayer.Number)
	fmt.Println("Peers:", status.ConnectedPeers)
}
