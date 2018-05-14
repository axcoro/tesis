package secuencial

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

// BoardS es una implementacion secuencial del juego de la vida
type BoardS struct {
	cells  [][]cellS
	h, w   int
	render bool
}

var reader = bufio.NewReader(os.Stdin)

// Render dibuja el tablero, solo si b.render == true
func (b BoardS) Render() {
	if b.render {
		// clear
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		// render
		fmt.Print(b.String())

		// wait input
		reader.ReadString('\n')
	}
}

// Init sirve para establecer las condiciones iniciales del tablero
func (b *BoardS) Init(w, h, prob int, render bool) {

	b.w = w
	b.h = h
	b.render = render

	b.cells = make([][]cellS, w)

	for x := 0; x < w; x++ {
		b.cells[x] = make([]cellS, h)
		for y := 0; y < h; y++ {
			b.cells[x][y] = cellS{
				board: b,
				alive: rand.Intn(101) < prob,
			}
		}
	}
}

// Next lleva el tablero al proximo estado
func (b BoardS) Next() {
	cant := b.w * b.h

	cells := make(chan *cellS, cant)

	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {
			cell := b.cells[x][y]
			cells <- cell.next(x, y)
		}
	}

	close(cells)
	for c := range cells {
		c.apply()
	}

	b.Render()
}

func (b BoardS) String() string {
	var buf bytes.Buffer

	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {
			cell := b.cells[x][y]
			if cell.alive {
				buf.WriteByte('*')
			} else {
				buf.WriteByte(' ')
			}
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (b BoardS) getCell(x, y int) *cellS {
	if x < 0 || x > b.w {
		return nil
	}

	if y < 0 || y > b.h {
		return nil
	}

	return &b.cells[x][y]
}
