package main

import (
	_ "embed"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/mint"
)

//go:embed mintMap.json
var mintMap []byte
var minter mint.PieceMinter

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	minter = mint.NewMinter(mintMap)

	addr, err := minter.MintPiece(`0xf8d6e0586b0a20c7`, `49b89eea-4c1e-44f1-ab94-b9fdd0f72457`)
	if err != nil {
		panic(err)
	}
	fmt.Print(addr)
}
