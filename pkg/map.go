package pkg

import (
	"image"
	"image/color"
	"math"

	perlin "github.com/aquilax/go-perlin"
	"github.com/go-gl/mathgl/mgl64"
)

type HeightMap [][]float64

func (hm HeightMap) Get(x, y int) float64 {
	return hm[x][y]
}

func (hm *HeightMap) Set(x, y int, val float64) {
	c := (*hm)
	c[x][y] = val
	hm = &c
}

// ColorModel returns the Image's color model.
func (hm HeightMap) ColorModel() color.Model {
	return color.GrayModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (hm HeightMap) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(hm), len(hm))
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (hm HeightMap) At(x int, y int) color.Color {
	return color.Gray{uint8((hm[x][y] + 1) * 128)}
}

func (hm HeightMap) Truncate() []float64 {
	l := len(hm)
	n := make([]float64, 0)

	for i := 0; i < l; i++ {
		n = append(n, hm[i]...)
	}
	return n
}

func NewPerlinMap(size, scale int) HeightMap {
	if scale == 0 {
		panic("cant be zero")
	}
	m := make(HeightMap, size)

	noise := perlin.NewPerlin(2, 2, 10, 420)
	for x := 0; x < size; x++ {
		m[x] = make([]float64, size)
		for y := 0; y < size; y++ {
			m[x][y] = noise.Noise2D(float64(x)/float64(scale), float64(y)/float64(scale))
		}
	}
	return m
}

func NewHeightMap(size int) HeightMap {
	m := make(HeightMap, size)

	for i := 0; i < size; i++ {
		m[i] = make([]float64, size)
	}
	return m
}

type NormalMap [][]mgl64.Vec3

// ColorModel returns the Image's color model.
func (nm NormalMap) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (nm NormalMap) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(nm), len(nm))
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (nm NormalMap) At(x int, y int) color.Color {
	return color.RGBA{uint8((nm[x][y].X() + 1) * 128), uint8((nm[x][y].Y() + 1) * 128), uint8(math.Abs(nm[x][y].Z()) * 255), 255}
}

func (nm NormalMap) Truncate() []mgl64.Vec3 {
	l := len(nm)
	n := make([]mgl64.Vec3, 0)

	for i := 0; i < l; i++ {
		n = append(n, nm[i]...)
	}
	return n
}

func NewNormalMap(hm HeightMap, length float64) NormalMap {
	m := make(NormalMap, len(hm))
	m[0] = make([]mgl64.Vec3, len(hm[0]))
	m[len(m)-1] = make([]mgl64.Vec3, len(hm[0]))

	for x := 1; x < len(m)-1; x++ {
		m[x] = make([]mgl64.Vec3, len(hm[0]))
		for y := 1; y < len(m[x])-1; y++ {
			m[x][y] = CalculateNormal(x, y, hm, length)
		}
	}
	return m
}

func CalculateNormal(x, y int, hm HeightMap, length float64) mgl64.Vec3 {
	n1 := mgl64.Vec3{float64(x), float64(y - 1), (hm[x][y-1] + 1) * length}
	n2 := mgl64.Vec3{float64(x - 1), float64(y), (hm[x-1][y] + 1) * length}
	n3 := mgl64.Vec3{float64(x + 1), float64(y), (hm[x+1][y] + 1) * length}
	n4 := mgl64.Vec3{float64(x), float64(y + 1), (hm[x][y+1] + 1) * length}
	return n2.Sub(n3).Cross(n1.Sub(n4)).Normalize()
}
