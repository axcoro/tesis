package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
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

func printAndWait(board string) {
	fmt.Print(board)
	time.Sleep(200 * time.Millisecond)
	// reader.ReadString('\n')
}

var (
	h = flag.Int("h", 10, "Largo del tablero")
	w = flag.Int("w", 30, "Ancho del tablero")
	p = flag.Int("p", 25, "Propabilidad de que una celda este viva al inicio")
	t = flag.Int("t", 600, "Cantidad de ciclos")
	s = flag.Int64("s", 0, "Semilla para rnd")
	r = flag.Bool("r", false, "Mostrar paso a paso la evolucion del tablero")
)

func main() {

	flag.Parse()

	n := *t
	seed := *s

	if seed == 0 { // sino se define una semilla definir algo pseudo rnd
		seed = time.Now().UTC().UnixNano()
	}
	rand.Seed(seed)

	b := board.Board{}
	str := b.Init(*h, *w, *p, *r)
	if str != "" {
		clear()
		fmt.Printf("0%% (0/%d)\n", n)
		printAndWait(str)
	}

	clear()
	for i := 0; i < n; i++ {
		fmt.Printf("%d%% (%d/%d)\n", (((i + 1) * 100) / n), i+1, n)
		str := b.Next()
		if str != "" {
			printAndWait(str)
		}
		clear()
	}
}
