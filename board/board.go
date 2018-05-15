package board

// Board reprensenta un tablero del juego de la vida
type Board interface {
	Init(width, height, prop, times int, render bool)
	Next()
	String() string
	Render()
}
