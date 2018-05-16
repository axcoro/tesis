package workers

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
)

// BoardW es una implementacion paralela del juego de la vida
type BoardW struct {
	cellW [][]*cellW
	h, w  int

	render bool

	elapsed int
	times   int
}

var reader = bufio.NewReader(os.Stdin)

const minWorkers = 1

// Render dibuja el tablero, solo si b.render == true
func (b *BoardW) Render() {

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

func workerInit(tasksInit <-chan [2]int, board *BoardW, prob int, wg *sync.WaitGroup) {
	for {
		coord, ok := <-tasksInit
		if !ok {
			return
		}

		x, y := coord[0], coord[1]

		board.cellW[x][y] = &cellW{
			board: board,
			alive: rand.Intn(101) < prob,
		}

		wg.Done()
	}
}

func workerNext(tasksNext <-chan [2]int, cells chan<- *cellW, board *BoardW, wg *sync.WaitGroup) {
	for {
		coord, ok := <-tasksNext
		if !ok {
			return
		}

		x, y := coord[0], coord[1]

		cell := board.cellW[x][y]
		cells <- cell.next(x, y)

		wg.Done()
	}
}

func workerApply(cells <-chan *cellW, wg *sync.WaitGroup) {
	for {
		cell, ok := <-cells
		if !ok {
			return
		}

		cell.apply()

		wg.Done()
	}
}

// Init sirve para establecer las condiciones iniciales del tablero
func (b *BoardW) Init(w, h, prob, times int, render bool) {

	b.w = w
	b.h = h

	b.render = render
	b.elapsed = 0
	b.times = times

	b.cellW = make([][]*cellW, w)

	wg := sync.WaitGroup{}
	wg.Add(w * h)
	tasksInit := make(chan [2]int)

	for i := 0; i < minWorkers; i++ {
		go workerInit(tasksInit, b, prob, &wg)
	}

	for x := 0; x < w; x++ {
		b.cellW[x] = make([]*cellW, h)
		for y := 0; y < h; y++ {
			tasksInit <- [2]int{x, y}
		}
	}

	wg.Wait()
	close(tasksInit)

	b.Render()
}

// Next lleva el tablero al proximo estado
func (b *BoardW) Next() {

	tasksNext := make(chan [2]int)
	cells := make(chan *cellW, b.w*b.h)

	wg := sync.WaitGroup{}
	wg.Add(b.w * b.h)

	for i := 0; i < minWorkers; i++ {
		go workerNext(tasksNext, cells, b, &wg)
	}

	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {
			tasksNext <- [2]int{x, y}
		}
	}

	wg.Wait()
	close(cells)

	wg.Add(b.w * b.h)
	for i := 0; i < minWorkers; i++ {
		go workerApply(cells, &wg)
	}

	wg.Wait()
	b.Render()
}

func (b BoardW) String() string {
	var buf bytes.Buffer

	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {
			cell := b.cellW[x][y]
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

func (b BoardW) getCell(x, y int) *cellW {
	if x < 0 || x >= b.w {
		return nil
	}

	if y < 0 || y >= b.h {
		return nil
	}

	return b.cellW[x][y]
}
