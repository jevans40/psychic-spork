package file

import (
	"fmt"
	"image"
	"testing"
	"unicode"

	"github.com/jevans40/psychic-spork/linmath"
	"github.com/jevans40/psychic-spork/logging"
)

func TestLoadImageFromFile(t *testing.T) {
	img1 := LoadImageFromFile("./../res/image/God2.png")
	img2 := LoadImageFromFile("./../res/image/debt.PNG")
	img1 = DrawToImage(img1, img2, image.Point{0, 0})
	SaveImageToFile("./../res/image/CaptureGod.png", img1)
}

func TestLoadingFonts(t *testing.T) {
	font := LoadFont("./../testres/font/OpenSans-SemiboldItalic.ttf", 1000)
	font.SetPos(linmath.NewPSPoint(10, 10))
	testText := "The quick brown fox jumped over the lazy dog"
	lastnum := rune(-1)
	for _, c := range testText {
		fontCoords := font.DrawChar(c, lastnum)
		lastnum = c
		logging.Log.Noticef("Draw At %s, %v,%v\n", string(c), fontCoords.X(), fontCoords.Y())
	}
	font.SetPos(linmath.NewPSPoint(10, 10))
	points := font.DrawString(testText)

	for i, c := range testText {

		logging.Log.Noticef("Draw At %s, %v,%v\n", string(c), points[i].X(), points[i].Y())
	}

	logging.Log.Notice("Done")

	for _, c := range testText {
		if !unicode.IsSpace(c) {
			img := font.GetGlyphImage(c)
			SaveImageToFile(fmt.Sprintf("./../testres/alphabet/%s.png", string(c)), img)
		}
	}
}
