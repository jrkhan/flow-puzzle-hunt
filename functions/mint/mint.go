package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "embed"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/jrkhan/flow-puzzle-hunt/cors"
	"github.com/jrkhan/flow-puzzle-hunt/verifysig"
)

//go:embed mint.cdc
var mintTx string

type (
	Envelope struct {
		SignedMessage verifysig.SignedMessage `json:"signedMessage"`
	}
)

func init() {
	functions.HTTP("MintFuzzle", HandleMintRequest)
}

func tx() string {
	formatted := strings.Replace(mintTx, `"../../contracts/NonFungibleToken.cdc"`, `${NON_FUNGIBLE_TOKEN_ADDRESS}`, -1)
	formatted = strings.Replace(formatted, `"../../contracts/FuzzlePieceV2.cdc"`, `${FUZZLE_PIECE_V2_ADDRESS}`, -1)
	return os.ExpandEnv(formatted)
}

func HandleMintRequest(w http.ResponseWriter, r *http.Request) {
	cors.HandleCors(w, r)
	var envelope = &Envelope{}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(envelope); err != nil {
		fmt.Fprint(w, "error decoding message")
		return
	}
	ctx := context.Background()
	if err := envelope.SignedMessage.VerifySignature(ctx, w); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	// message has been verified

	fmt.Fprint(w, "message verified!")
}

func GetAccessNode() string {
	val, has := os.LookupEnv("FLOW_ACCESS_NODE")
	if has {
		return val
	}
	return "access.devnet.nodes.onflow.org:9000"
}

func GetMinter() string {
	val, has := os.LookupEnv("FUZZLE_MINTER_KEY")
	if !has {
		return ""
	}
	return val
}
