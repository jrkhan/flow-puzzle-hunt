package mint_test

import (
	"embed"
	"testing"

	"github.com/jrkhan/flow-puzzle-hunt/pkg/mint"
)

//go:embed testdata/puzzles/*.json
var files embed.FS

//go:embed testdata/mintMap.json
var mintMap []byte

func TestNewMinter(t *testing.T) {
	mm := mint.NewMinter(mintMap, files)
	p := mm.PieceMap("0d679fe8-f231-40df-8366-e56d32260206")
	if p.PieceID != 1 {
		t.Logf("Expected %v to be %v", p.PieceID, 1)
		t.Fail()
	}
	p = mm.PieceMap("99f1ed71-4be4-44f4-8138-6465f12aaeed")
	if p.PieceID != 5 {
		t.Logf("Expected %v to be %v", p.PieceID, 5)
		t.Fail()
	}
}
