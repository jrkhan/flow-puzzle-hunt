package redirect

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/jrkhan/flow-puzzle-hunt/pkg/cors"
)

type (
	Lookup struct {
		GUID string `json:"guid"`
	}
	MV struct {
		PieceID int `json:"pieceId"`
	}
)

//go:embed mintMap.json
var mapRaw []byte
var guidMap = map[string]MV{}

func init() {
	err := json.Unmarshal(mapRaw, &guidMap)
	if err != nil {
		panic("unable to build mintMap")
	}
	functions.HTTP("LookupPieceByGuid", LookupPiece)
}

func LookupPiece(w http.ResponseWriter, r *http.Request) {
	cors.HandleCors(w, r)
	var lookup = &Lookup{}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(lookup); err != nil {
		fmt.Fprint(w, "error decoding message")
		return
	}

	val, has := guidMap[lookup.GUID]
	if !has {
		fmt.Fprintf(w, "0")
		return
	}
	fmt.Fprint(w, val.PieceID)
}
