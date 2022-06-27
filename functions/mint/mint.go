package mint

import (
	_ "embed"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/mint"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/verifysig"
)

//go:embed mintMap.json
var mintMap []byte

type (
	Envelope struct {
		SignedMessage verifysig.SignedMessage `json:"signedMessage"`
	}
)

func init() {
	minter := mint.NewMinter(mintMap)
	functions.HTTP("MintFuzzle", minter.HandleMintRequest)
}
