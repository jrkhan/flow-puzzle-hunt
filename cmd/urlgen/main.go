package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

//go:embed input.json
var input []byte

type (
	Piece struct {
		IPFS        string `json:"ipfs"`
		PieceID     int    `json:"pieceID"`
		PuzzleID    int    `json:"puzzleID"`
		DisplayName string `json:"displayName"`
		MintURL     string `json:"mintURL"`
	}
	Puzzle struct {
		CAR      string `json:"car"`      // the content-addressed archive to use see -> https://car.ipfs.io/
		Gateway  string `json:"gateway"`  // will be included in the ipfs url
		Count    int    `json:"count"`    // piece count
		PuzzleID int    `json:"puzzleId"` // id of the puzzle - will be eventually tied to a separate mintable nft
		MintPath string `json:"mintpath"` // the base domain/path of the url for QR codes
		Prefix   string `json:"prefix"`   // puzzle-1-PREFIX-1 will be the display name and gateway/CAR-CID/PREFIX-1 the url
		Ext      string `json:"ext"`      // file extension for each piece
	}
	Input struct {
		Puzzle Puzzle `json:"puzzle"`
	}
)

func ParseInput(inp []byte) Puzzle {
	input := &Input{}
	err := json.Unmarshal(inp, input)
	if err != nil {
		panic(err)
	}
	return input.Puzzle
}

func (p Puzzle) IPFS(pieceID int) string {
	return p.Gateway + p.CAR + "/" + p.Prefix + strconv.Itoa(pieceID) + "." + p.Ext
}

func (p Puzzle) PieceDisplayName(pieceID int) string {
	return fmt.Sprintf("puzzle-%v-%v%v", p.PuzzleID, p.Prefix, pieceID)
}

func (p Puzzle) MintURL(guid string) string {
	return fmt.Sprintf("%v%v", p.MintPath, guid)
}

func main() {
	p := ParseInput(input)
	// one of our outputs will be a list of pieces
	mp := map[string]Piece{}
	for i := 1; i <= p.Count; i++ {
		guid := uuid.New().String()
		mp[guid] = Piece{
			IPFS:        p.IPFS(i),
			PieceID:     i,
			PuzzleID:    p.PuzzleID,
			DisplayName: p.PieceDisplayName(i),
			MintURL:     p.MintURL(guid),
		}
	}
	rs, err := json.MarshalIndent(mp, "", "\t")
	if err != nil {
		panic("unable to build result")
	}
	fmt.Println(string(rs))
}
