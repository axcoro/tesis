package board

import (
	"bytes"
	"math/rand"
	"strconv"
	"sync"
)

type Board struct {
	cells map[string]*Cell
	h, w  int
}

func (b *Board) Init(w, h, prob int) string {

	b.w = w
	b.h = h
	b.cells = make(map[string]*Cell)

	wg := sync.WaitGroup{}
	wg.Add(w * h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			go func(_x, _y int) {
				alive := rand.Intn(101) < prob
				cell := Cell{state: alive, x: _x, y: _y}
				cell.Board(b)
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()
	return b.String()
}

func (b Board) Next() string {

	// cant := b.w * b.h

	// wg := sync.WaitGroup{}
	// wg.Add(cant)
	// cells := make(chan *Cell, cant)

	for x := 0; x < b.w; x++ {
		// go func(_x int) {
		for y := 0; y < b.h; y++ {
			cell := b.cell(x, y)
			cell.Next()
			defer cell.Apply()
			// wg.Done()
		}
		// }(x)
	}

	// wg.Wait()
	// close(cells)

	// for c := range cells {
	// 	c.Apply()
	// }

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
