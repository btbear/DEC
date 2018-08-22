// Copyright 2016 The go-DEWH Authors
// This file is part of the go-DEWH library.
//
// The go-DEWH library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-DEWH library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-DEWH library. If not, see <http://www.gnu.org/licenses/>.

// Contains all the wrappers from the node package to support client side node
// management on mobile platforms.

package geth

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/DEWH/go-DEWH/core"
	"github.com/DEWH/go-DEWH/eth"
	"github.com/DEWH/go-DEWH/eth/downloader"
	"github.com/DEWH/go-DEWH/ethclient"
	"github.com/DEWH/go-DEWH/ethstats"
	"github.com/DEWH/go-DEWH/internal/debug"
	"github.com/DEWH/go-DEWH/les"
	"github.com/DEWH/go-DEWH/node"
	"github.com/DEWH/go-DEWH/p2p"
	"github.com/DEWH/go-DEWH/p2p/nat"
	"github.com/DEWH/go-DEWH/params"
	whisper "github.com/DEWH/go-DEWH/whisper/whisperv6"
)

// NoDEWHonfig represents the collection of configuration values to fine tune the Geth
// node embedded into a mobile process. The available values are a subset of the
// entire API provided by go-DEWH to reduce the maintenance surface and dev
// complexity.
type NoDEWHonfig struct {
	// Bootstrap nodes used to establish connectivity with the rest of the network.
	BootstrapNodes *Enodes

	// MaxPeers is the maximum number of peers that can be connected. If this is
	// set to zero, then only the configured static and trusted peers can connect.
	MaxPeers int

	// DEWHEnabled specifies whether the node should run the DEWH protocol.
	DEWHEnabled bool

	// DEWHNetworkID is the network identifier used by the DEWH protocol to
	// DEWHide if remote peers should be accepted or not.
	DEWHNetworkID int64 // uint64 in truth, but Java can't handle that...

	// DEWHGenesis is the genesis JSON to use to seed the blockchain with. An
	// empty genesis state is equivalent to using the mainnet's state.
	DEWHGenesis string

	// DEWHDatabaseCache is the system memory in MB to allocate for database caching.
	// A minimum of 16MB is always reserved.
	DEWHDatabaseCache int

	// DEWHNetStats is a netstats connection string to use to report various
	// chain, transaction and node stats to a monitoring server.
	//
	// It has the form "nodename:secret@host:port"
	DEWHNetStats string

	// WhisperEnabled specifies whether the node should run the Whisper protocol.
	WhisperEnabled bool

	// Listening address of pprof server.
	PprofAddress string
}

// defaultNoDEWHonfig contains the default node configuration values to use if all
// or some fields are missing from the user's specified list.
var defaultNoDEWHonfig = &NoDEWHonfig{
	BootstrapNodes:        FoundationBootnodes(),
	MaxPeers:              25,
	DEWHEnabled:       true,
	DEWHNetworkID:     1,
	DEWHDatabaseCache: 16,
}

// NewNoDEWHonfig creates a new node option set, initialized to the default values.
func NewNoDEWHonfig() *NoDEWHonfig {
	config := *defaultNoDEWHonfig
	return &config
}

// Node represents a Geth DEWH node instance.
type Node struct {
	node *node.Node
}

// NewNode creates and configures a new Geth node.
func NewNode(datadir string, config *NoDEWHonfig) (stack *Node, _ error) {
	// If no or partial configurations were specified, use defaults
	if config == nil {
		config = NewNoDEWHonfig()
	}
	if config.MaxPeers == 0 {
		config.MaxPeers = defaultNoDEWHonfig.MaxPeers
	}
	if config.BootstrapNodes == nil || config.BootstrapNodes.Size() == 0 {
		config.BootstrapNodes = defaultNoDEWHonfig.BootstrapNodes
	}

	if config.PprofAddress != "" {
		debug.StartPProf(config.PprofAddress)
	}

	// Create the empty networking stack
	noDEWHonf := &node.Config{
		Name:        clientIdentifier,
		Version:     params.VersionWithMeta,
		DataDir:     datadir,
		KeyStoreDir: filepath.Join(datadir, "keystore"), // Mobile should never use internal keystores!
		P2P: p2p.Config{
			NoDiscovery:      true,
			DiscoveryV5:      true,
			BootstrapNodesV5: config.BootstrapNodes.nodes,
			ListenAddr:       ":0",
			NAT:              nat.Any(),
			MaxPeers:         config.MaxPeers,
		},
	}
	rawStack, err := node.New(noDEWHonf)
	if err != nil {
		return nil, err
	}

	debug.Memsize.Add("node", rawStack)

	var genesis *core.Genesis
	if config.DEWHGenesis != "" {
		// Parse the user supplied genesis spec if not mainnet
		genesis = new(core.Genesis)
		if err := json.Unmarshal([]byte(config.DEWHGenesis), genesis); err != nil {
			return nil, fmt.Errorf("invalid genesis spec: %v", err)
		}
		// If we have the testnet, hard code the chain configs too
		if config.DEWHGenesis == TestnetGenesis() {
			genesis.Config = params.TestnetChainConfig
			if config.DEWHNetworkID == 1 {
				config.DEWHNetworkID = 3
			}
		}
	}
	// Register the DEWH protocol if requested
	if config.DEWHEnabled {
		ethConf := eth.DefaultConfig
		ethConf.Genesis = genesis
		ethConf.SyncMode = downloader.LightSync
		ethConf.NetworkId = uint64(config.DEWHNetworkID)
		ethConf.DatabaseCache = config.DEWHDatabaseCache
		if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			return les.New(ctx, &ethConf)
		}); err != nil {
			return nil, fmt.Errorf("DEWH init: %v", err)
		}
		// If netstats reporting is requested, do it
		if config.DEWHNetStats != "" {
			if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
				var lesServ *les.LightDEWH
				ctx.Service(&lesServ)

				return ethstats.New(config.DEWHNetStats, nil, lesServ)
			}); err != nil {
				return nil, fmt.Errorf("netstats init: %v", err)
			}
		}
	}
	// Register the Whisper protocol if requested
	if config.WhisperEnabled {
		if err := rawStack.Register(func(*node.ServiceContext) (node.Service, error) {
			return whisper.New(&whisper.DefaultConfig), nil
		}); err != nil {
			return nil, fmt.Errorf("whisper init: %v", err)
		}
	}
	return &Node{rawStack}, nil
}

// Start creates a live P2P node and starts running it.
func (n *Node) Start() error {
	return n.node.Start()
}

// Stop terminates a running node along with all it's services. If the node was
// not started, an error is returned.
func (n *Node) Stop() error {
	return n.node.Stop()
}

// GetDEWHClient retrieves a client to access the DEWH subsystem.
func (n *Node) GetDEWHClient() (client *DEWHClient, _ error) {
	rpc, err := n.node.Attach()
	if err != nil {
		return nil, err
	}
	return &DEWHClient{ethclient.NewClient(rpc)}, nil
}

// GetNodeInfo gathers and returns a collection of metadata known about the host.
func (n *Node) GetNodeInfo() *NodeInfo {
	return &NodeInfo{n.node.Server().NodeInfo()}
}

// GetPeersInfo returns an array of metadata objects describing connected peers.
func (n *Node) GetPeersInfo() *PeerInfos {
	return &PeerInfos{n.node.Server().PeersInfo()}
}
