package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/axcoro/tesis/go/GameOfLife/board"
)

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

var reader = bufio.NewReader(os.Stdin)

func wait() {
	time.Sleep(200 * time.Millisecond)
	reader.ReadString('\n')
}

func main() {
	defer board.Un(board.Trace("main"))
	h := flag.Int("h", 10, "Largo del tablero")
	w := flag.Int("w", 30, "Ancho del tablero")
	p := flag.Int("p", 25, "Propabilidad de que una celda este viva al inicio")
	r := flag.Bool("r", false, "Mostrar paso a paso la evolucion del tablero")

	flag.Parse()

	render := *r

	b := board.Board{}
	if *r {
		fmt.Printf(b.Init(*h, *w, *p, render))
	} else {
		b.Init(*h, *w, *p, render)
	}

	clear()
	for i := 0; i < 600; i++ {
		fmt.Printf("%d%%\n", (i*100)/600)
		if render {
			fmt.Printf(b.Next(render))
			wait()
		} else {
			b.Next(render)
		}
		clear()
	}
}
