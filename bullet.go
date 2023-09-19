package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Bullet struct {
	GameObject
	image       *ebiten.Image
	speedFactor int
}

func NewBullet(cfg *Config, ship *Ship) *Bullet {
	rect := image.Rect(0, 0, cfg.BulletWidth, cfg.BulletHeight)
	img := ebiten.NewImageWithOptions(rect, nil)
	img.Fill(cfg.BulletColor)
	gameObj := GameObject{
		width:  cfg.BulletWidth,
		height: cfg.BulletHeight,
		x:      ship.x + ((ship.width - cfg.BulletWidth) / 2),
		y:      cfg.ScreenHeight - ship.height - cfg.BulletHeight,
	}
	fmt.Printf("NewBullet -> width: %d, height: %d, x: %d, y: %d \n", gameObj.width, gameObj.height, gameObj.x, gameObj.y)
	return &Bullet{
		GameObject:  gameObj,
		image:       img,
		speedFactor: cfg.BulletSpeedFactor,
	}
}

func (bullet *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(bullet.x), float64(bullet.y))
	screen.DrawImage(bullet.image, op)
}

func (bullet *Bullet) outOfScreen() bool {
	return bullet.y < -bullet.height
}
