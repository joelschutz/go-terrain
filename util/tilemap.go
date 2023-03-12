package util

import (
	"image"
	"image/draw"
)

type TileMap struct {
	src              image.Image
	tXCount, tYCount uint
}

func (tm *TileMap) TXSize() uint {
	return uint(tm.src.Bounds().Dx() / int(tm.tXCount))
}

func (tm *TileMap) TYSize() uint {
	return uint(tm.src.Bounds().Dy() / int(tm.tYCount))
}

func (tm *TileMap) GetTile(index uint) image.Image {
	sx := (index % tm.tXCount) * tm.TXSize()
	sy := (index / tm.tXCount) * tm.TYSize()

	r := image.Rect(int(sx), int(sy), int(sx+tm.TXSize()), int(sy+tm.TYSize()))
	t := image.NewRGBA(r)
	draw.Draw(t, r, tm.src, image.Point{}, draw.Src)
	return t
}

func NewTileMap(src image.Image, tXCount, tYCount uint) *TileMap {
	return &TileMap{src, tXCount, tYCount}
}
