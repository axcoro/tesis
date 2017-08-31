package board

import "strconv"

type Cell struct {
	state bool // true: alive, false: dead
	next  bool // true: alive, false: dead
	x     int
	y     int
	board *Board
}

func (c Cell) Alive() bool {
	return c.state == true
}

func (c Cell) Board(b *Board) {
	c.board = b
	b.cells[c.point()] = &c
}

func (c Cell) point() string {
	return strconv.Itoa(c.x) + strconv.Itoa(c.y)
}

func (c Cell) neighbor(x, y int) *Cell {

	_x := c.x + x
	_y := c.y + y

	point := strconv.Itoa(_x) + strconv.Itoa(_y)

	return c.board.cells[point]
}

func (c Cell) neighbors() []*Cell {

	return []*Cell{
		c.neighbor(-1, -1), c.neighbor(0, -1), c.neighbor(1, -1),
		c.neighbor(-1, 0), c.neighbor(1, 0),
		c.neighbor(-1, 1), c.neighbor(0, 1), c.neighbor(1, 1),
	}

}

func (c Cell) neighborsAlive() int {

	neighbors := c.neighbors()

	alive := 0
	for _, neighbor := range neighbors {
		if neighbor != nil && neighbor.Alive() {
			alive = alive + 1
		}
	}

	return alive
}

func (c *Cell) Next() {

	alive := c.neighborsAlive()

	// Rules:
	//   3 neighbors: alive/on,
	//   2 neighbors: maintain,
	//   otherwise: dead/off.

	// calcular luego de que se leyeron todos los valores
	c.next = (alive == 3) || (alive == 2 && c.state)
}

func (c *Cell) Apply() {
	c.state = c.next
}

func (c *Cell) Row() byte {
	if c.state {
		return '*'
	}

	return 'X'
}
