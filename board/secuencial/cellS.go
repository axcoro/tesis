package secuencial

type cellS struct {
	alive     bool
	nextState bool
	board     *BoardS
}

func (c *cellS) next(x, y int) *cellS {

	b := c.board

	neighbors := [8]*cellS{
		b.getCell(x-1, y-1), b.getCell(x, y-1), b.getCell(x+1, y-1),
		b.getCell(x-1, y) /*b.getCell(x,y]*/, b.getCell(x+1, y),
		b.getCell(x-1, y+1), b.getCell(x, y+1), b.getCell(x+1, y+1),
	}

	alive := 0
	for _, neighbor := range neighbors {
		if neighbor != nil && neighbor.alive {
			alive = alive + 1
		}
	}

	// regla: S2B3
	c.nextState = (alive == 3) || (alive == 2 && c.alive)

	return c
}

func (c *cellS) apply() {
	c.alive = c.nextState
	c.nextState = false
}
