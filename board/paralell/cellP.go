package paralell

type cellP struct {
	alive     bool
	nextState bool
	board     *BoardP
}

func (c *cellP) next(x, y int) *cellP {

	b := c.board

	neighbors := [8]*cellP{
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

func (c *cellP) apply() {
	c.alive = c.nextState
	c.nextState = false
}
