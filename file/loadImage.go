package file

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/jevans40/psychic-spork/logging"
)

type notOpaqueRGBA struct {
	*image.RGBA
}

func (i *notOpaqueRGBA) Opaque() bool {
	return false
}

func LoadImageFromFile(filename string) *image.RGBA {
	f, err := os.Open(filename)
	if err != nil {
		logging.Log.Panic(err)
	}
	defer f.Close()

	img, fmtName, err := image.Decode(f)
	logging.Log.Info(fmtName)
	if err != nil {
		logging.Log.Panic(err)
	}
	b := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	return m
}

/*
func SaveImageToFile(filename string, img image.RGBA) {
	f, err := os.Create(filename)
	if err != nil {
		logging.Log.Panic(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		logging.Log.Panic(err)
	}
}
*/
func DrawToImage(src, dst image.Image, sp image.Point) *image.RGBA {
	b := dst.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), dst, b.Min, draw.Src)
	r := image.Rectangle{sp, sp.Add(src.Bounds().Size())}
	draw.Draw(m, r, src, src.Bounds().Min, draw.Src)
	return m
}
