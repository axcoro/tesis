package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/axcoro/tesis/go/GameOfLife/board"
)

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

var reader = bufio.NewReader(os.Stdin)

func printAndWait(board string) {
	fmt.Print(board)
	// time.Sleep(200 * time.Millisecond)
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

	str := b.Init(*h, *w, *p, render)
	if render {
		clear()
		fmt.Println("0%")
		printAndWait(str)
	}

	clear()
	for i := 0; i < 600; i++ {
		fmt.Printf("%d%%\n", (i*100)/600)
		str := b.Next(render)
		if render {
			printAndWait(str)
		}
		clear()
	}
}
