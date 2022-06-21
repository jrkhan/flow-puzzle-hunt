package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed flow.json
var config []byte

func checkErr(msg string, err error) {
	if err != nil {
		panic(msg + err.Error())
	}
}
func main() {
	ctx := context.Background()
	cfg := &FlowConfig{}
	err := json.Unmarshal(config, cfg)
	checkErr("could not unmarshal config", err)

	fc, err := client.New(cfg.Networks.Testnet, grpc.WithTransportCredentials(insecure.NewCredentials()))
	checkErr("could not connect to client", err)

	latestBlock, err := fc.GetLatestBlock(ctx, true)
	printBlock(latestBlock, err)
}

type (
	FlowConfig struct {
		Networks Networks `json:"networks"`
	}
	Networks struct {
		Mainnet string `json:"mainnet"`
		Testnet string `json:"testnet"`
	}
)

func printBlock(block *flow.Block, err error) {
	checkErr("error getting latest blocK", err)
	fmt.Printf("\nID: %s\n", block.ID)
	fmt.Printf("height: %d\n", block.Height)
	fmt.Printf("timestamp: %s\n\n", block.Timestamp)
}
