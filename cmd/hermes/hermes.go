package main

import (
	"flag"
	"fmt"
	"hermes/internal/hermes"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	fmt.Println("Hermes (Ἑρμῆς) - An Key-Val data store")

	if *cpuprofile != "" {
		fmt.Println("Profiling enabled")
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		h := hermes.New(":3333")
		go h.Listen()
		var input string
		fmt.Scanln(&input)
		fmt.Println("Shutting down...")
		return
	}

	h := hermes.New(":3333")
	h.Listen()

}
