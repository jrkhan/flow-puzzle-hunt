package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/png"
	"net/http"
	"os"

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
		e := Encoder{
			AlphaThreshold: 2000,
			GreyThreshold:  50,
			QRLevel:        qr.Highest, // recommended, as logo steals some redundant space
		}

		rsz := resize.Resize(480, 0, img, resize.Lanczos3)
		qr, err := e.Encode(v.MintURL, rsz, 1280)
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

// Encoder defines settings for QR/Overlay encoder.
type Encoder struct {
	AlphaThreshold int
	GreyThreshold  int
	QRLevel        qr.RecoveryLevel
}

// DefaultEncoder is the encoder with default settings.
var DefaultEncoder = Encoder{
	AlphaThreshold: 2000,       // FIXME: don't remember where this came from
	GreyThreshold:  30,         // in percent
	QRLevel:        qr.Highest, // recommended, as logo steals some redundant space
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	return DefaultEncoder.Encode(str, logo, size)
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func (e Encoder) Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	code, err := qr.New(str, e.QRLevel)
	if err != nil {
		return nil, err
	}

	img := code.Image(size)

	gs := image.NewPaletted(img.Bounds(), palette.Plan9)
	draw.Draw(gs, img.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Over)
	e.overlayLogo(gs, logo)

	err = png.Encode(&buf, gs)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// overlayLogo blends logo to the center of the QR code,
// changing all colors to black.
func (e Encoder) overlayLogo(dst, src image.Image) {
	//grey := uint32(^uint16(0)) * uint32(e.GreyThreshold) / 100
	alphaOffset := uint32(e.AlphaThreshold)
	offset := dst.Bounds().Max.X/2 - src.Bounds().Max.X/2
	for x := 0; x < src.Bounds().Max.X; x++ {
		for y := 0; y < src.Bounds().Max.Y; y++ {
			if _, _, _, alpha := src.At(x, y).RGBA(); alpha > alphaOffset {

				dst.(*image.Paletted).Set(x+offset, y+offset, src.At(x, y))
			}
		}
	}
}
