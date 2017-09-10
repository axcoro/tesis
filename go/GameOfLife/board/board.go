package board

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

type Board struct {
	cells map[string]*Cell
	h, w  int
}

func (b *Board) Init(w, h, prob int) string {

	b.w = w
	b.h = h
	b.cells = make(map[string]*Cell)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			alive := rand.Intn(101) < prob
			cell := Cell{state: alive, x: x, y: y}
			cell.Board(b)
		}
	}

	return b.String()
}

func (b Board) Next() string {

	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {
			cell := b.cell(x, y)
			cell.Next()
			defer cell.Apply()
		}
	}

	return b.String()
}

func (b Board) cell(x, y int) *Cell {
	point := strconv.Itoa(x) + strconv.Itoa(y)
	return b.cells[point]
}

func (b Board) String() string {
	var buf bytes.Buffer

	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {
			point := strconv.Itoa(x) + strconv.Itoa(y)
			cell := b.cells[point]
			if cell != nil {
				b := cell.Row()
				buf.WriteByte(b)
			}
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
