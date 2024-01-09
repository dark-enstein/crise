package config

import (
	"log"
	"strings"
)

var (
	FLAG_GRACEFULTERM = "GracefulTermination"
	FLAG_PLAYGAME     = "PlayGame"
)

var Features = map[string]FFlags{
	FLAG_GRACEFULTERM: FlagDisabled,
	FLAG_PLAYGAME:     FlagDisabled,
}

func IsValidFeature(s string) bool {
	splice := strings.Split(s, ",")
	for i := 0; i < len(splice); i++ {
		if _, ok := Features[splice[i]]; !ok {
			log.Printf("invalid feature flag: %s\n", splice[i])
			return false
		}
	}
	return true
}

type FFlags int

const (
	FlagDisabled FFlags = iota
	FlagEnabled
)

func ResolveFeatureFlags(s string) map[string]FFlags {
	splice := strings.Split(s, ",")
	for i := 0; i < len(splice); i++ {
		Features[splice[i]] = FlagEnabled
		log.Printf("Enabled featuure flag: %s\n", splice[i])
	}
	return Features
}
