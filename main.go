package main

import (
	"context"
	"fmt"
	"github.com/dark-enstein/crise/config"
	"github.com/dark-enstein/crise/game"
	"github.com/dark-enstein/crise/internal/glyph"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
)

var cfg = &config.Gameplay{}
var Debug = "false"
var FeatureFlags = ""
var PlayGame = ""

var Fonts = map[string][]byte{
	"RobotoLight":        glyph.RobotoLight,
	"MPlus1pRegular_ttf": glyph.MPlus1pRegular_ttf,
	"RobotoBold":         glyph.RobotoBold,
	"MontserratBold":     glyph.MontserratBold,
	"MontserratLight":    glyph.MontserratLight,
	"TitilliumWebBold":   glyph.TitilliumWebBold,
	"TitilliumWebLight":  glyph.TitilliumWebLight,
	"UbuntuBold":         glyph.UbuntuBold,
	"UbuntuLight":        glyph.UbuntuLight,
}

//type GameManager struct {
//	log
//}

func init() {
	var err error
	configFile := "config.json"
	cfg, err = config.Unmarshal(configFile)
	if err != nil {
		log.Println("error parsing tetris config file:", configFile)
		os.Exit(1)
	}
	if _, ok := Fonts[cfg.Font]; !ok {
		log.Printf("font type received not valid: %s\n", cfg.Font)
		os.Exit(1)
	}
	tt, err := opentype.Parse(Fonts[cfg.Font])
	if err != nil {
		log.Println("error parsing ttf font file")
		os.Exit(1)
	}
	game.DefaultFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    cfg.Size,
		DPI:     game.DEFAULT_DPI,
		Hinting: font.HintingVertical,
	})
	fmt.Printf("Metrics. LineHeight: %d, All: %#v \n", int(game.DefaultFont.Metrics().Height), game.DefaultFont.Metrics())
}

func main() {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Crise: a tetris incarnate")
	ff := map[string]config.FFlags{}
	mode := 0
	if Debug == "true" {
		mode = game.EXEC_DEBUG
		log.Println("Running in debug mode")
		if len(FeatureFlags) < 1 {
			log.Println("No feature flags enabled")
		} else {
			if !config.IsValidFeature(FeatureFlags) {
				log.Println("Feature flag invalid")
				os.Exit(1)
			}
			ff = config.ResolveFeatureFlags(FeatureFlags)
		}
	} else {
		mode = game.EXEC_PRODUCTION
		log.Println("Running in production mode")
	}

	if mode != game.EXEC_PRODUCTION && len(ff) > 0 && ff[config.FLAG_GRACEFULTERM] > 0 {
		var gE errgroup.Group
		gE.Go(func() error {
			return ebiten.RunGame(game.NewGameWithContext(&ctx, cancel, game.INIT, cfg, mode, ff))
		})
		select {
		default:
			log.Println("Game exited successfully")
		case <-ctx.Done():
			err := context.Cause(ctx)
			if err != nil {
				log.Printf("Halted game due to: %s\n", err.Error())
			} else {
				log.Println("Game exited successfully, without errors.")
			}
		}
		if err := gE.Wait(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := ebiten.RunGame(game.NewGame(cancel, game.INIT, cfg, mode, ff)); err != nil {
			log.Fatal(err)
		}
	}
}
