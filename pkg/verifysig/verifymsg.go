package verifysig

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"os"

	_ "embed"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	flowGrpc "github.com/onflow/flow-go-sdk/access/grpc"
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

//go:embed verifysignatures.cdc
var verifySignaturesQuery string

func GetAccessNode() string {
	val, has := os.LookupEnv("FLOW_ACCESS_NODE")
	if has {
		return val
	}
	return "access.devnet.nodes.onflow.org:9000"
}

func (s *SignedMessage) VerifySignature(ctx context.Context, w io.Writer) error {
	// bootstrap client
	fc, err := flowGrpc.NewClient(GetAccessNode(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	// fmt.Fprint(w, args[0].String())
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
	addr := flow.HexToAddress(s.Address)
	cAddr := cadence.NewAddress(addr)

	// convert message to hex, then to string
	hx := hex.EncodeToString([]byte(s.Message))
	cMessage, err := cadence.NewString(hx)
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
