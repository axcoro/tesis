package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
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
	// time.Sleep(200 * time.Millisecond)
	reader.ReadString('\n')
}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	h          = flag.Int("h", 10, "Largo del tablero")
	w          = flag.Int("w", 30, "Ancho del tablero")
	p          = flag.Int("p", 25, "Propabilidad de que una celda este viva al inicio")
	t          = flag.Int("t", 600, "Cantidad de ciclos")
	s          = flag.Int64("s", 0, "Semilla para rnd")
	r          = flag.Bool("r", false, "Mostrar paso a paso la evolucion del tablero")
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	render := *r
	n := *t

	if !render {
		defer board.Un(board.Trace("main"))
	}

	seed := *s
	if seed == 0 { // sino se define una semilla definir algo pseudo rnd
		seed = time.Now().UTC().UnixNano()
	}
	rand.Seed(seed)

	b := board.Board{}
	str := b.Init(*h, *w, *p, render)

	if render {
		clear()
		fmt.Printf("0%% (0/%d)\n", n)
		printAndWait(str)
	}

	clear()
	for i := 0; i < n; i++ {
		fmt.Printf("%d%% (%d/%d)\n", (((i + 1) * 100) / n), i+1, n)
		str := b.Next(render)
		if render {
			printAndWait(str)
		}
		clear()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
