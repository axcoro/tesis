package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/axcoro/tesis/board"
	"github.com/axcoro/tesis/board/secuencial"
)

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

	b := secuencial.BoardS{}

	b.Init(*h, *w, *p, *t, render)

	for i := 0; i < n; i++ {
		b.Next()
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
