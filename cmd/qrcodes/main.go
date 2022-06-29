package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os"

	"github.com/divan/qrlogo"
	"github.com/nfnt/resize"
	qr "github.com/skip2/go-qrcode"
)

//go:embed mintMap.json
var mintMap []byte

type Piece struct {
	IPFS        string `json:"ipfs"`
	PieceID     int    `json:"pieceID"`
	PuzzleID    int    `json:"puzzleID"`
	DisplayName string `json:"displayName"`
	MintURL     string `json:"mintURL"`
}

func main() {
	o := &map[string]Piece{}
	err := json.Unmarshal(mintMap, o)
	if err != nil {
		panic(err)
	}
	for k, v := range *o {
		img := imageFromURL(v.IPFS)
		e := qrlogo.Encoder{
			AlphaThreshold: 200,        // FIXME: don't remember where this came from
			GreyThreshold:  80,         // in percent
			QRLevel:        qr.Highest, // recommended, as logo steals some redundant space
		}

		rsz := resize.Resize(580, 0, img, resize.Lanczos3)
		qr, err := e.Encode(v.MintURL, rsz, 2048)
		if err != nil {
			panic(err)
		}
		out, err := os.Create(fmt.Sprintf("./images/qr-%v.png", k))
		if err != nil {
			panic(err)
		}
		out.Write(qr.Bytes())
	}
}

func imageFromURL(url string) image.Image {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		panic(err)
	}
	return m
}
