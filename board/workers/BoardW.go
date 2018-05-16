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

func worker(tasksCh <-chan workerParams, wg *sync.WaitGroup) {
	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}
		task.fn(task.b, task.x, task.y)
		wg.Done()
	}
}

type workerParams struct {
	fn   func(b *BoardW, x, y int)
	b    *BoardW
	x, y int
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
	tasksCh := make(chan workerParams)

	for i := 0; i < 5; i++ {
		go worker(tasksCh, &wg)
	}

	for x := 0; x < w; x++ {
		b.cellW[x] = make([]*cellW, h)
		for y := 0; y < h; y++ {
			params := workerParams{
				func(_b *BoardW, _x, _y int) {
					_b.cellW[_x][_y] = &cellW{
						board: _b,
						alive: rand.Intn(101) < prob,
					}
				},
				b, x, y,
			}
			tasksCh <- params
		}
	}
	wg.Wait()
	close(tasksCh)

	b.Render()
}

// Next lleva el tablero al proximo estado
func (b *BoardW) Next() {

	cellW := make(chan *cellW, b.w*b.h)

	wg := sync.WaitGroup{}
	wg.Add(b.w)

	for x := 0; x < b.w; x++ {
		go func(_x int) {
			for y := 0; y < b.h; y++ {
				cell := b.cellW[_x][y]
				cellW <- cell.next(_x, y)
			}
			wg.Done()
		}(x)
	}

	wg.Wait()
	close(cellW)
	for c := range cellW {
		c.apply()
	}

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
