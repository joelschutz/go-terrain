package util

import "image"

func AsRGBA(img image.Image) *image.RGBA {
	rect := img.Bounds()
	n := image.NewRGBA(rect)
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			n.Set(x, y, img.At(x, y))
		}
	}
	return n
}
