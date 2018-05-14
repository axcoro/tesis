package board

// Board reprensenta un tablero del juego de la vida
type Board interface {
	Init(width int, height int, prop int, times int, render bool) string
	Next(render bool) string
	String() string
	Render() string
}
