package board

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
)

const (
	_defaultWorkerSize = 10

	_alive = true
	_dead  = false

	_actionSearch   = 1
	_actionNoSearch = 0

	_survives   = 1
	_noSurvives = 0
)

type Board struct {
	cells map[[2]int]bool

	aliveCells chan [2]int

	render bool

	h int
	w int

	mutex sync.RWMutex
}

func (b *Board) Init(h, w, prob int, render bool) string {

	b.h = h
	b.w = w
	b.render = render

	b.cells = make(map[[2]int]bool, w*h)
	b.aliveCells = make(chan [2]int, w*h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {

			cell := [2]int{x, y}
			state := (rand.Intn(101) < prob)

			b.cells[cell] = state
			if state == _alive {
				b.aliveCells <- cell
			}
		}
	}

	if b.render {
		return b.String()
	}

	return ""
}

func (b *Board) workerFindNeighbors(cells chan [3]int, neighbors chan [2]int, results chan<- [3]int, wg *sync.WaitGroup) {

	var (
		x          int
		y          int
		aliveCount int

		isAlive bool

		cell     [2]int
		neighbor [2]int

		action int
	)

	for {
		data, ok := <-cells
		if !ok {
			fmt.Printf("workerFindNeighbors - finalizando\n")
			return
		}

		fmt.Printf("workerFindNeighbors - procesando:%v\n", cell)

		x, y, action, aliveCount = data[0], data[1], data[2], 0

		cell = [2]int{x, y}

		for _x := x - 1; _x <= x+1; _x++ { // desplazamiento horizontal de -1,0,1 respecto a x

			if _x < 0 || _x >= b.w {
				continue
			}

			for _y := y - 1; _y <= y+1; _y++ { // desplazamiento horizontal de -1,0,1 respecto a y

				if _y < 0 || _y >= b.h {
					continue
				}

				neighbor = [2]int{_x, _y}

				b.mutex.RLock()
				isAlive = b.cells[neighbor]
				b.mutex.RUnlock()

				if isAlive {
					fmt.Printf("workerFindNeighbors - celda:%v - vecino:%v\n", cell, neighbor)
					aliveCount = aliveCount + 1
				}

				if action == _actionSearch {
					fmt.Printf("workerFindNeighbors - %v notifica a :%v | cells:%v\n", cell, neighbor, len(cells))
					neighbors <- neighbor
				}

			}
		}

		nextStatus := _noSurvives
		if aliveCount == 3 { // B3?
			nextStatus = _survives
		} else {

			b.mutex.RLock()
			currentStatus := b.cells[cell]
			b.mutex.RUnlock()

			if currentStatus && aliveCount == 2 { // S2?
				nextStatus = _survives
			}
		}

		fmt.Printf("workerFindNeighbors - resultado para %v: %v | results:%v\n", cell, nextStatus, len(results))

		results <- [3]int{x, y, nextStatus}

		wg.Done()

	}
}

func (b *Board) Next() string {

	// recorro las vivas
	// cuento cuantos vecinos vivos tiene
	// verifico rule y notifico si queda viva o no
	// por cada vecino
	// cuanto cuentos vecinos vivos tiene
	// verifico rule y notifico si queda viva o no

	// wait

	// marco cada celda viva o muerta como tal
	// mantengo las celdas vivas

	wg := sync.WaitGroup{}
	totalWork := len(b.aliveCells) * 9 // celdas vivas + vecinas

	if totalWork > b.h*b.w { // nunca puede haber mas trabajo que celdas
		totalWork = b.h * b.w
	}

	cells := make(chan [3]int, totalWork)
	neighbors := make(chan [2]int, totalWork)

	results := make(chan [3]int, totalWork)

	aliveCells := make(chan [2]int, totalWork)
	deadCells := make(chan [2]int, totalWork)

	fmt.Printf("Next - totalWork:%v\n", totalWork)

	wg.Add(totalWork)

	// de todas las celdas, las unicas que me interesa validar son:
	// celdas vivas y sus vecinas

	go func() {
		for i := 0; i < _defaultWorkerSize; i++ {
			go b.workerFindNeighbors(cells, neighbors, results, &wg)
			// se crean aca, pero no se usan hasta el final
			go b.workerSetAlive(aliveCells, &wg)
			go b.workerSetDead(deadCells, &wg)
		}
	}()

	processedCells := make(map[[2]int]bool, totalWork)
	m := sync.RWMutex{}

	go func() { // mando cada celda viva a los workers para que busquen los vecinos
		for cell := range b.aliveCells {
			fmt.Printf("Next - mando a buscar vecinos para [%v %v] | b.aliveCells:%v\n", cell[0], cell[1], len(b.aliveCells))
			m.Lock()
			processedCells[cell] = true
			m.Unlock()
			cells <- [3]int{cell[0], cell[1], _actionSearch}
		}
	}()

	go func() { // mando cada celda viva a los workers para que busquen los vecinos
		for neighbor := range neighbors {
			m.RLock()
			processed := processedCells[neighbor]
			m.RUnlock()
			if !processed {
				fmt.Printf("Neighbors - mando a procesar %v\n", neighbor)
				m.Lock()
				processedCells[neighbor] = true
				m.Unlock()
				cells <- [3]int{neighbor[0], neighbor[1], _actionNoSearch}
			} else {
				fmt.Printf("Neighbors - %v ya esta procesado\n", neighbor)
			}
		}
	}()

	wg.Wait()
	close(cells)
	close(results)
	// close(neighbors)

	fmt.Printf("Next - fase 2 | results:%v\n", len(results))

	b.aliveCells = make(chan [2]int, b.w*b.h) // en el unico momento en que deberia estar abierto este canal
	wg.Add(totalWork)

	for result := range results {

		cell := [2]int{result[0], result[1]}
		nextStatus := result[2]

		switch nextStatus {
		case _survives:
			aliveCells <- cell
		default: // no survive
			deadCells <- cell
		}
	}

	wg.Wait()

	close(aliveCells)
	close(deadCells)
	close(b.aliveCells) // siempre deberia estar cerrado

	if b.render {
		return b.String()
	}

	// cuando
	return ""
}

func (b Board) String() string {
	var buf bytes.Buffer

	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {

			if b.cells[[2]int{x, y}] { // alive?
				buf.WriteByte('*')
			} else {
				buf.WriteByte(' ')
			}
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

// workers

func (b *Board) workerSearchTargetCells(cells <-chan [2]int, neighbors chan<- [2]int, wg *sync.WaitGroup) {

	notified := make(map[[2]int]bool, b.w*b.h)
	mutex := sync.Mutex{}
	for {
		cell, ok := <-cells
		if !ok {
			return
		}

		x, y := cell[0], cell[1]

		for _x := x - 1; _x <= x+1; _x++ { // desplazamiento horizontal de -1,0,1 respecto a x

			if _x < 0 || _x > b.w {
				continue
			}

			for _y := y - 1; _y <= y+1; _y++ { // desplazamiento horizontal de -1,0,1 respecto a y

				if _y < 0 || _y > b.h {
					continue
				}

				c := [2]int{_x, _y}

				mutex.Lock()
				_, processed := notified[cell]
				if processed {
					mutex.Unlock()
					continue
				}

				notified[cell] = true
				mutex.Unlock()

				// incluyo la celda (x,y) adrede
				neighbors <- c
			}
		}

		wg.Done()

	}

}

func (b *Board) workerCheckCells(cells <-chan [2]int, aliveCells chan<- [2]int, deadCells chan<- [2]int, wg *sync.WaitGroup) {

	// validar si la celda esta viva y si cumple con las reglas

	// marcar la celda en b.cells como viva / muerta

	// si esta viva, mandarla a aliveCells, sino, a deadCells

	for {
		cell, ok := <-cells
		if !ok {
			return
		}

		x, y, alive := cell[0], cell[1], 0

		for _x := x - 1; _x <= x+1; _x++ { // desplazamiento horizontal de -1,0,1 respecto a x

			if _x < 0 || _x > b.w { // _x esta fuera de los limites
				continue
			}

			for _y := y - 1; _y <= y+1; _y++ { // desplazamiento horizontal de -1,0,1 respecto a y

				if _y < 0 || _y > b.h || (_x == x && _y == y) {
					// _x esta fuera de los limites
					// (_x,_y) == (x,y) es la celda a evaluar, por lo que la ignoro
					continue
				}

				neighbor := [2]int{_x, _y}
				b.mutex.RLock()
				neighborAlive := b.cells[neighbor]
				b.mutex.RUnlock()

				if neighborAlive { // la celda vecina esta viva?
					alive = alive + 1
				}

			}
		}

		// validar regla (B3/S2)

		if alive == 3 { // B3?
			aliveCells <- cell
		} else {

			b.mutex.RLock()
			c := b.cells[cell]
			b.mutex.RUnlock()

			if c && alive == 2 { // S2?
				aliveCells <- cell
			} else {
				deadCells <- cell
			}
		}

		wg.Done()

	}
}

func (b *Board) workerSetAlive(cells <-chan [2]int, wg *sync.WaitGroup) {
	for { // marco como viva + agrego al listado de celdas vivas para la prox ronda
		cell, ok := <-cells
		if !ok {
			return
		}
		b.mutex.Lock()
		b.cells[cell] = true
		b.mutex.Unlock()
		b.aliveCells <- cell
		wg.Done()
	}
}

func (b *Board) workerSetDead(cells <-chan [2]int, wg *sync.WaitGroup) {
	for { // marco como viva + agrego al listado de celdas vivas para la prox ronda
		cell, ok := <-cells
		if !ok {
			return
		}
		b.mutex.Lock()
		b.cells[cell] = false
		b.mutex.Unlock()
		wg.Done()
	}
}
