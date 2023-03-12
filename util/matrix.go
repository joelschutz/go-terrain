package util

import "math/rand"

func MakeMatrixUint8(size int) [][]uint8 {
	m := [][]uint8{}

	for x := 0; x < size; x++ {
		row := []uint8{}
		for y := 0; y < size; y++ {
			row = append(row, 0)
		}
		m = append(m, row)
	}

	return m
}

func MakeMatrixBool(size int) [][]bool {
	m := [][]bool{}

	for x := 0; x < size; x++ {
		row := []bool{}
		for y := 0; y < size; y++ {
			row = append(row, false)
		}
		m = append(m, row)
	}

	return m
}

func MakeRandMatrixBool(size, threshold int) [][]bool {
	m := [][]bool{}

	for x := 0; x < size; x++ {
		row := []bool{}
		for y := 0; y < size; y++ {
			row = append(row, rand.Intn(101) < threshold)
		}
		m = append(m, row)
	}

	return m
}
