package main

import (
	"fmt"
	"hermes/internal/hermes"
)

func main() {
	fmt.Println("Hermes (Ἑρμῆς) - An Key-Val data store")

	h := hermes.New(":3333")
	h.Listen()
}
