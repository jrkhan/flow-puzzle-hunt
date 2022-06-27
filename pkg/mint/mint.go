package mint

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	flowGrpc "github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var contractMap = map[string]string{
	`"../../contracts/NonFungibleToken.cdc"`: `${NON_FUNGIBLE_TOKEN_ADDRESS}`,
	`"../../contracts/FuzzlePieceV2.cdc"`:    `${FUZZLE_PIECE_V2_ADDRESS}`,
}

//go:embed mint.cdc
var mintTx string

type (
	PieceMinter struct {
		pieceMap map[string]Piece
	}
	Piece struct {
		IPFS        string `json:"ipfs"`
		PieceID     int    `json:"pieceID"`
		PuzzleID    int    `json:"puzzleID"`
		DisplayName string `json:"displayName"`
		MintURL     string `json:"mintURL"`
	}

	MintRequest struct {
		Address string
		Piece   Piece
	}
)

func txScript() []byte {
	formatted := mintTx
	for k, v := range contractMap {
		formatted = strings.Replace(mintTx, k, v, -1)
	}
	return []byte(os.ExpandEnv(formatted))
}

func NewMinter(source []byte) PieceMinter {
	mp := &map[string]Piece{}
	json.Unmarshal(source, mp)
	return PieceMinter{*mp}
}

func (m *PieceMinter) MintPiece(addr string, key string) (*flow.Identifier, error) {
	piece := m.pieceMap[key]
	mr := MintRequest{Address: addr, Piece: piece}
	cv, err := mr.ToCadenceValues()
	if err != nil {
		return nil, err
	}

	// bootstrap client
	fc, err := flowGrpc.NewClient(GetAccessNode(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	latestBlock, err := fc.GetLatestBlockHeader(context.Background(), true)
	proposerAddress := flow.HexToAddress(GetMinterAddress())
	proposerAccount, err := fc.GetAccountAtLatestBlock(context.Background(), proposerAddress)
	proposerKeyIndex := 0
	// what guarentees do we have about this sequence number?
	// do we need a distributed lock here?
	sequenceNumber := proposerAccount.Keys[proposerKeyIndex].SequenceNumber

	tx := flow.NewTransaction().
		SetScript(txScript()).
		SetGasLimit(100).
		SetReferenceBlockID(latestBlock.ID).
		SetProposalKey(proposerAddress, proposerKeyIndex, sequenceNumber).
		SetPayer(proposerAddress).
		AddAuthorizer(proposerAddress)
	for _, arg := range cv {
		tx.AddArgument(arg)
	}
	sk := proposerAccount.Keys[0]
	// construct a signer from your private key and configured hash algorithm
	pk, err := crypto.DecodePrivateKey(sk.SigAlgo, GetMinterKey())
	signer, err := crypto.NewInMemorySigner(pk, proposerAccount.Keys[0].HashAlgo)
	if err != nil {
		panic("failed to create a signer")
	}
	tx.SignEnvelope(proposerAddress, 0, signer)

	err = fc.SendTransaction(context.Background(), *tx)
	id := tx.ID()
	return &id, nil
}

func GetAccessNode() string {
	val, has := os.LookupEnv("FLOW_ACCESS_NODE")
	if has {
		return val
	}
	return "access.devnet.nodes.onflow.org:9000"
}

func GetMinterAddress() string {
	val, has := os.LookupEnv("FUZZLE_MINTER_ADDRESS")
	if !has {
		return ""
	}
	return val
}
func GetMinterKey() []byte {
	val, has := os.LookupEnv("FUZZLE_MINTER_KEY")
	if !has {
		panic("FUZZLE_MINTER_KEY is required")
	}
	return []byte(val)
}

func (m *MintRequest) ToCadenceValues() ([]cadence.Value, error) {
	// convert address
	addr := flow.HexToAddress(m.Address)
	cAddr := cadence.NewAddress(addr)

	// display name
	dn, err := cadence.NewString(m.Piece.DisplayName)
	if err != nil {
		return nil, err
	}

	puzId := cadence.NewUInt64(uint64(m.Piece.PuzzleID))
	pieceId := cadence.NewUInt64(uint64(m.Piece.PieceID))

	return []cadence.Value{
		cAddr, dn, puzId, pieceId,
	}, nil
}
