package redirect

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type (
	Lookup struct {
		GUID string `json:"guid"`
	}
)

//go:embed mintMap.json
var mapRaw []byte
var guidMap = map[string]int{}

func init() {
	err := json.Unmarshal(mapRaw, &guidMap)
	if err != nil {
		panic("unable to build mintMap")
	}
	functions.HTTP("LookupPieceByGuid", LookupPiece)
}
func HandleCors(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func LookupPiece(w http.ResponseWriter, r *http.Request) {
	HandleCors(w, r)
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
	fmt.Fprint(w, val)
}
