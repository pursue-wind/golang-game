package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeOver
)

var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
)

type GameAttr interface {
	X() int
	Y() int
	Width() int
	Height() int
}

type GameObject struct {
	width  int
	height int
	x      int
	y      int
}

func (gameObj *GameObject) Width() int {
	return gameObj.width
}

func (gameObj *GameObject) Height() int {
	return gameObj.height
}

func (gameObj *GameObject) X() int {
	return gameObj.x
}

func (gameObj *GameObject) Y() int {
	return gameObj.y
}

type Input struct {
	msg            string
	lastBulletTime time.Time
}

func (i *Input) Update(g *Game) {
	ship := g.ship
	cfg := g.cfg
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		//fmt.Println("←←←←←←←←←←←←←←←←←←←←←←←")
		i.msg = "left pressed"

		if (0 - g.ship.width/2) < ship.x {
			ship.x -= ship.moveStep
		}

	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		//fmt.Println("→→→→→→→→→→→→→→→→→→→→→→→")
		i.msg = "right pressed"
		if ship.x < (cfg.ScreenWidth - g.ship.width/2) {
			ship.x += ship.moveStep
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		//fmt.Println("-----------------------")
		i.msg = "space pressed"
		if time.Now().Sub(i.lastBulletTime).Milliseconds() > g.cfg.BulletInterval {
			bullet := NewBullet(g.cfg, g.ship)
			g.addBullet(bullet)
			i.lastBulletTime = time.Now()
		}
	}
}

type Game struct {
	mode       Mode
	score      int64
	input      *Input
	cfg        *Config
	ship       *Ship
	bullets    map[*Bullet]struct{}
	aliens     map[*Alien]bool
	aliensLock sync.Mutex
}

func NewGame() *Game {
	cfg := loadConfig()
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle(cfg.Title)
	g := &Game{
		score:   0,
		input:   &Input{},
		cfg:     cfg,
		ship:    NewShip(cfg),
		bullets: make(map[*Bullet]struct{}),
		aliens:  make(map[*Alien]bool),
	}
	g.CreateFonts()
	go g.createAliens()
	return g
}

func (g *Game) ResetGame() {
	g.mode = ModeTitle
	g.score = 0
	g.input = &Input{}
	g.bullets = make(map[*Bullet]struct{})
	g.aliens = make(map[*Alien]bool)
}

func (g *Game) Update() error {
	for bullet := range g.bullets {
		if bullet.outOfScreen() {
			fmt.Println("delete bullet")
			delete(g.bullets, bullet)
		}
		bullet.y -= g.cfg.BulletSpeedFactor
	}
	g.aliensLock.Lock()
	for alien, f := range g.aliens {
		if alien.outOfScreen(g.cfg) {
			fmt.Println("delete alien")
			delete(g.aliens, alien)
		}
		alien.y += g.cfg.AlienSpeedFactor + int(g.score/5)
		if f {
			alien.x = alien.x - 1
		} else {
			alien.x = alien.x + 1
		}
	}
	g.CheckCollision()
	g.CheckCollision2()
	g.aliensLock.Unlock()
	g.input.Update(g)
	return nil
}

func (g *Game) CreateFonts() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(32),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(32),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.cfg.BgColor)
	var texts []string
	switch g.mode {
	case ModeTitle:
		texts = []string{"PRESS SPACE KEY", "START GAME"}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.mode = ModeGame
		}
	case ModeGame:
		g.ship.Draw(screen, g.cfg)
		for bullet := range g.bullets {
			bullet.Draw(screen)
		}
		for alien := range g.aliens {
			alien.Draw(screen)
		}
		text.Draw(screen, strconv.FormatInt(g.score, 10), titleArcadeFont, g.cfg.ScreenWidth-64, 48, color.Opaque)
	case ModeOver:
		texts = []string{"GAME OVER!", "YOU SCORE IS " + strconv.FormatInt(g.score, 10), "PRESS ENTER RESTART"}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.ResetGame()
		}
	}
	for i, l := range texts {
		text.Draw(screen, l, arcadeFont, 0, 100*(i+1), color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

func (g *Game) addBullet(bullet *Bullet) {
	g.bullets[bullet] = struct{}{}
}

func (g *Game) createAliens() {
	for {
		g.aliensLock.Lock()
		alien := NewAlien(g.cfg)
		randTime := rand.Intn(2000)
		randX := rand.Intn(g.cfg.ScreenWidth)
		randY := rand.Intn(g.cfg.ScreenWidth / 3)
		alien.x = randX
		alien.y = randY

		flag := rand.Intn(100)
		f := true
		if flag > 50 {
			f = false
		}
		g.aliens[alien] = f
		fmt.Printf("NewAlien: width: %d, height: %d, x: %d, y: %d \n", alien.width, alien.height, alien.x, alien.y)
		g.aliensLock.Unlock()
		time.Sleep(time.Millisecond * time.Duration(randTime))
	}
}

func (g *Game) CheckCollision() {
	for alien := range g.aliens {
		for bullet := range g.bullets {
			if CheckCollision(bullet, alien) {
				delete(g.aliens, alien)
				delete(g.bullets, bullet)
				g.score += 1
			}
		}
	}
}

func (g *Game) CheckCollision2() {
	for alien := range g.aliens {
		if CheckCollision(alien, g.ship) {
			g.mode = ModeOver
		}
	}
}
