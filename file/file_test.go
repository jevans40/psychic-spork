package file

import (
	"image"
	"testing"
)

func TestLoadImageFromFile(t *testing.T) {
	img1 := LoadImageFromFile("./../res/image/God2.png")
	img2 := LoadImageFromFile("./../res/image/debt.PNG")
	img1 = DrawToImage(img1, img2, image.Point{0, 0})
	SaveImageToFile("./../res/image/CaptureGod.png", img1)
}
