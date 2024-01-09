package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Gameplay struct {
	Face        `json:"Face", yaml:"Face"`
	HistoryPlay `json:"HistoryPlay", yaml:"HistoryPlay"`
}

type Face struct {
	Size float64 `json:"size", yaml:"Size"`
	Font string  `json:"font", yaml:"Font"`
}

type HistoryPlay struct {
	Highest int `json:"highest", yaml:"Highest"`
	You     int `json:"you", yaml:"You"`
}

type Developer struct {
	Sprite SpriteConfig
}

type SpriteConfig struct {
	Height int
	Width  int
}

// In Json
func Unmarshal(fileLoc string) (*Gameplay, error) {
	_, err := os.Stat(fileLoc)
	if os.IsNotExist(err) {
		log.Println(fmt.Errorf("error while opening config file: %s: %w\n",
			fileLoc, err).Error())
		return nil, err
	}
	cfgBytes, err := os.ReadFile(fileLoc)
	if err != nil {
		log.Println(fmt.Errorf("error while opening config file: %s: %w\n",
			fileLoc, err).Error())
		return nil, err
	}

	var cfg Gameplay
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		log.Println(fmt.Errorf("error while decofing config file: %s: %w\n",
			fileLoc, err).Error())
		return nil, err
	}

	return &cfg, nil
}
