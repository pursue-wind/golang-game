package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type Alien struct {
	GameObject
	image       *ebiten.Image
	speedFactor int
}

func NewAlien(cfg *Config) *Alien {
	img, _, err := ebitenutil.NewImageFromFile("images/alien.png")
	if err != nil {
		log.Fatal(err)
	}
	width, height := img.Size()

	gameObj := GameObject{width: width, height: height, x: 0, y: 0}
	//fmt.Printf("NewAlien: width: %d, height: %d, x: %d, y: %d \n", gameObj.width, gameObj.height, gameObj.x, gameObj.y)

	return &Alien{
		GameObject:  gameObj,
		image:       img,
		speedFactor: cfg.AlienSpeedFactor,
	}
}

func (alien *Alien) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(alien.x), float64(alien.y))
	screen.DrawImage(alien.image, op)
}

func (alien *Alien) outOfScreen(cfg *Config) bool {
	return cfg.ScreenHeight < alien.y
}
