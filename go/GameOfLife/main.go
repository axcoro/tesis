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

	h := flag.Int("h", 10, "Largo del tablero")
	w := flag.Int("w", 10, "Ancho del tablero")
	p := flag.Int("p", 25, "Propabilidad de que una celda este viva al inicio")

	flag.Parse()

	b := board.Board{}
	fmt.Printf(b.Init(*h, *w, *p))

	clear()
	for i := 0; i < 600; i++ {
		fmt.Printf(b.Next())
		wait()
		clear()
	}
}
