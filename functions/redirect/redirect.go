package redirect

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/cors"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/mint"
)

type (
	Lookup struct {
		GUID string `json:"guid"`
	}
	MV struct {
		PieceID  int `json:"pieceID"`
		PuzzleID int `json:"puzzleID"`
	}
	Locator struct {
		mint.PieceMinter
	}
)

//go:embed mintMap.json
var mintMap []byte

//go:embed puzzles/*.json
var files embed.FS

func init() {
	minter := mint.NewMinter(mintMap, files)
	loc := Locator{PieceMinter: minter}
	functions.HTTP("LookupPieceByGuid", loc.LookupPiece)
}

func (l *Locator) LookupPiece(w http.ResponseWriter, r *http.Request) {
	cors.HandleCors(w, r)
	var lookup = &Lookup{}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(lookup); err != nil {
		fmt.Fprint(w, "error decoding message")
		return
	}

	val := l.PieceMap(lookup.GUID)

	res, err := json.Marshal(val)
	if err != nil {
		fmt.Fprintf(w, "unable to marshal result")
	}
	fmt.Fprint(w, string(res))
}
