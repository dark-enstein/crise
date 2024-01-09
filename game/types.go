package game

import "github.com/dark-enstein/crise/internal/utils"

type Phase int

const (
	// WELCOME phase defines the welcome screen, and what happens during this time period
	WELCOME Phase = iota
	SCORE
	START
	HELP
	QUIT
	DEBUG
	INIT
)

var (
	PHASES = map[Phase]P_CONTENT{
		WELCOME: P_CONTENT{
			"Welcome to Tetris",
		},
		SCORE: {},
		START: P_CONTENT{
			`"Score: %d\n`,
		},
		HELP:  {},
		QUIT:  {},
		DEBUG: {},
	}
	SPRITE_W = 291 / utils.SPRITE_WIDTH
	SPRITE_H = 494 / utils.SPRITE_HEIGHT
)

type P_CONTENT struct {
	Text string
}
