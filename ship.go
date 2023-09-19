package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type Ship struct {
	GameObject
	image    *ebiten.Image
	moveStep int
}

func NewShip(cfg *Config) *Ship {
	img, _, err := ebitenutil.NewImageFromFile("images/ship.png")
	if err != nil {
		log.Fatal(err)
	}

	width, height := img.Size()
	gameObj := GameObject{width: width, height: height, x: (cfg.ScreenWidth - width) / 2, y: cfg.ScreenHeight - height}
	fmt.Printf("NewShip: width: %d, height: %d, x: %d, y: %d \n", gameObj.width, gameObj.height, gameObj.x, gameObj.y)

	ship := &Ship{
		GameObject: gameObj,
		image:      img,
		moveStep:   cfg.MoveStep,
	}

	return ship
}
func (ship *Ship) Draw(screen *ebiten.Image, cfg *Config) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(ship.x), float64(ship.y))
	screen.DrawImage(ship.image, op)
}
