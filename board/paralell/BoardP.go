package paralell

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
)

// BoardP es una implementacion paralela del juego de la vida
type BoardP struct {
	cellP [][]*cellP
	h, w  int

	render bool

	elapsed int
	times   int
}

var reader = bufio.NewReader(os.Stdin)

// Render dibuja el tablero, solo si b.render == true
func (b *BoardP) Render() {

	if b.render {
		// clear
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		// render
		fmt.Print(b.String())
		fmt.Printf("%d%% (%d/%d)\n", (((b.elapsed) * 100) / b.times), b.elapsed, b.times)
		// wait input
		reader.ReadString('\n')
	}

	b.elapsed = b.elapsed + 1
}

// Init sirve para establecer las condiciones iniciales del tablero
func (b *BoardP) Init(w, h, prob, times int, render bool) {

	b.w = w
	b.h = h

	b.render = render
	b.elapsed = 0
	b.times = times

	b.cellP = make([][]*cellP, w)

	wg := sync.WaitGroup{}
	wg.Add(w)

	for x := 0; x < w; x++ {
		go func(_x int) {
			b.cellP[_x] = make([]*cellP, h)
			for y := 0; y < h; y++ {
				b.cellP[_x][y] = &cellP{
					board: b,
					alive: rand.Intn(101) < prob,
				}
			}
			wg.Done()
		}(x)
	}

	wg.Wait()

	b.Render()
}

// Next lleva el tablero al proximo estado
func (b *BoardP) Next() {

	cellP := make(chan *cellP, b.w*b.h)

	wg := sync.WaitGroup{}
	wg.Add(b.w)

	for x := 0; x < b.w; x++ {
		go func(_x int) {
			for y := 0; y < b.h; y++ {
				cell := b.cellP[_x][y]
				cellP <- cell.next(_x, y)
			}
			wg.Done()
		}(x)
	}

	wg.Wait()
	close(cellP)
	for c := range cellP {
		c.apply()
	}

	b.Render()
}

func (b BoardP) String() string {
	var buf bytes.Buffer

	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {
			cell := b.cellP[x][y]
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

func (b BoardP) getCell(x, y int) *cellP {
	if x < 0 || x >= b.w {
		return nil
	}

	if y < 0 || y >= b.h {
		return nil
	}

	return b.cellP[x][y]
}
