package main

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
)

type Config struct {
	ScreenWidth       int        `json:"screenWidth"`
	ScreenHeight      int        `json:"screenHeight"`
	MoveStep          int        `json:"moveStep"`
	Title             string     `json:"title"`
	BgColor           color.RGBA `json:"bgColor"`
	BulletWidth       int        `json:"bulletWidth"`
	BulletHeight      int        `json:"bulletHeight"`
	BulletColor       color.RGBA `json:"bulletColor"`
	BulletSpeedFactor int        `json:"bulletSpeedFactor"`
	BulletInterval    int64      `json:"bulletInterval"`
	AlienSpeedFactor  int        `json:"alienSpeedFactor"`
}

func loadConfig() *Config {
	f, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("os.Open failed: %v\n", err)
	}

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		log.Fatalf("json.Decode failed: %v\n", err)
	}

	return &cfg
}
