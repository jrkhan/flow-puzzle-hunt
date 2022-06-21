package funcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	_ "embed"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	SignedMessage struct {
		Address     string   `json:"address"`
		Message     string   `json:"message"`
		KeyIndicies []int    `json:"keyIndices"`
		Signatures  []string `json:"signatures"`
	}
)

func init() {
	functions.HTTP("MintFuzzle", HandleMintRequest)
}

func HandleMintRequest(w http.ResponseWriter, r *http.Request) {
	var envelope struct {
		SignedMessage SignedMessage `json:"signedMessage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		fmt.Fprint(w, "error decoding message")
		return
	}
	ctx := context.Background()
	if err := envelope.SignedMessage.VerifySignature(ctx); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
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

//go:embed verifysignatures.cdc
var verifySignaturesQuery string

func (s *SignedMessage) VerifySignature(ctx context.Context) error {
	// bootstrap client
	fc, err := client.New(GetAccessNode(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	// expand script using env variables
	q := os.ExpandEnv(verifySignaturesQuery)

	// convert args to cadence values
	args, err := s.ToCadenceValues()
	if err != nil {
		return err
	}

	// execute script
	val, err := fc.ExecuteScriptAtLatestBlock(ctx, []byte(q), args)
	if err != nil {
		return err
	}

	if !val.ToGoValue().(bool) {
		return errors.New("signed message was invalid")
	}
	return nil
}

func (s *SignedMessage) ToCadenceValues() ([]cadence.Value, error) {
	// convert address
	cAddr, err := cadence.NewString(s.Address)
	if err != nil {
		return nil, err
	}
	// convert message
	cMessage, err := cadence.NewString(s.Message)
	if err != nil {
		return nil, err
	}
	// convert signatures
	cSigs := make([]cadence.Value, len(s.Signatures))
	for i, sig := range s.Signatures {
		cSigs[i], err = cadence.NewString(sig)
		if err != nil {
			return nil, err
		}
	}
	caSigs := cadence.NewArray(cSigs)

	// convert key indicies
	ckeyIndices := make([]cadence.Value, len(s.KeyIndicies))
	for i, ki := range s.KeyIndicies {
		ckeyIndices[i] = cadence.NewInt(ki)
	}
	caKeyIndices := cadence.NewArray(ckeyIndices)

	return []cadence.Value{
		cAddr, cMessage, caKeyIndices, caSigs,
	}, nil
}
