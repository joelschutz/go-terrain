package util

import "math"

const (
	Sqrt3 = 1.732050807568877293527446341505872366942805253810380628055806
)

type HexPoint [3]float64 // Q R S

func (h HexPoint) Q() float64 {
	return h[0]
}

func (h HexPoint) R() float64 {
	return h[1]
}

func (h HexPoint) S() float64 {
	return h[2]
}

func (h *HexPoint) Add(p HexPoint) HexPoint {
	return HexPoint{
		h[0] + p.Q(),
		h[1] + p.R(),
		h[2] + p.S(),
	}
}

func (h *HexPoint) Sub(p HexPoint) HexPoint {
	return HexPoint{
		h[0] - p.Q(),
		h[1] - p.R(),
		h[2] - p.S(),
	}
}

type Point2D [2]float64 // X Y

func (p Point2D) X() float64 {
	return p[0]
}

func (p Point2D) Y() float64 {
	return p[1]
}

func Pointy_hex_to_pixel(hex HexPoint, size float64) Point2D {
	x := size * (Sqrt3*hex.Q() + Sqrt3/2*hex.R())
	y := size * (3. / 2 * hex.R())
	return Point2D{x, y}
}

func Flat_hex_to_pixel(hex HexPoint, size float64) Point2D {
	x := size * (3. / 2 * hex.Q())
	y := size * (Sqrt3*hex.R() + Sqrt3/2*hex.Q())
	return Point2D{x, y}
}

func Pointy_hex_corner(p Point2D, size float64, i int) Point2D {
	var angle_rad = float64(i)*math.Pi/3 - math.Pi/6.
	return Point2D{p.X() + size*math.Cos(angle_rad), p.Y() + size*math.Sin(angle_rad)}
}

func Flat_hex_corner(p Point2D, size float64, i int) Point2D {
	var angle_rad = float64(i) * math.Pi / 3
	return Point2D{p.X() + size*math.Cos(angle_rad), p.Y() + size*math.Sin(angle_rad)}
}

func Pixel_to_pointy_hex(p Point2D, size float64) HexPoint {
	q := (Sqrt3/3*p.X() - 1./3*p.Y()) / size
	r := (2. / 3 * p.Y()) / size
	s := -q - r
	return HexPoint{q, r, s}
}

func Pixel_to_flat_hex(p Point2D, size float64) HexPoint {
	q := (2. / 3 * p.X()) / size
	r := (Sqrt3/3*p.Y() - 1./3*p.X()) / size
	s := -q - r
	return HexPoint{q, r, s}
}

func Cube_round(hex HexPoint) HexPoint {
	q := math.Round(hex.Q())
	r := math.Round(hex.R())
	s := math.Round(hex.S())

	q_diff := math.Abs(q - hex.Q())
	r_diff := math.Abs(r - hex.R())
	s_diff := math.Abs(s - hex.S())

	if q_diff > r_diff && q_diff > s_diff {
		q = -r - s
	} else if r_diff > s_diff {
		r = -q - s
	} else {
		s = -q - r
	}
	return HexPoint{q, r, s}
}

type HexPlane struct {
	Values []HexPoint
	Origin HexPoint
	Size   int
}

func NewHexBoard(qnt int) *HexPlane {
	b := HexPlane{}
	b.Size = qnt
	b.Origin = HexPoint{-float64(qnt), 0, float64(qnt)}
	for R := -(qnt / 2); R < (qnt / 2); R++ {
		for S := (qnt / 2); S > -(qnt / 2); S-- {
			Q := -R - S
			b.Values = append(b.Values, HexPoint{float64(Q), float64(R), float64(S)})
		}
	}
	return &b
}
