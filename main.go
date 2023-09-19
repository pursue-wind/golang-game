package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
