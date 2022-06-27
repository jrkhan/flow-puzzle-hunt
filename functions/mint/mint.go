package mint

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/cors"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/mint"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/verifysig"
)

//go:embed mintMap.json
var mintMap []byte

type (
	Envelope struct {
		SignedMessage verifysig.SignedMessage `json:"signedMessage"`
	}
	MinterHandler struct {
		mint.PieceMinter
	}
)

func init() {
	minter := &MinterHandler{mint.NewMinter(mintMap)}
	functions.HTTP("MintFuzzle", minter.HandleMintRequest)
}

func (m *MinterHandler) HandleMintRequest(w http.ResponseWriter, r *http.Request) {
	cors.HandleCors(w, r)
	var envelope = &Envelope{}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(envelope); err != nil {
		fmt.Fprint(w, "error decoding message")
		return
	}
	ctx := context.Background()
	signedMsg := envelope.SignedMessage
	if err := signedMsg.VerifySignature(ctx, w); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	// message has been verified ^
	id, err := m.MintPiece(signedMsg.Address, signedMsg.Message)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, id.Hex())
}
