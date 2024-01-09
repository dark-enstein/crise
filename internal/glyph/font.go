package glyph

import (
	_ "embed"
)

var (
	//go:embed resource/mplus-1p-regular.ttf
	MPlus1pRegular_ttf []byte

	//go:embed resource/montserrat/Montserrat-Bold.otf
	MontserratBold []byte

	//go:embed resource/montserrat/Montserrat-Light.otf
	MontserratLight []byte

	//go:embed resource/Roboto/Roboto-Bold.ttf
	RobotoBold []byte

	//go:embed resource/Roboto/Roboto-Light.ttf
	RobotoLight []byte

	//go:embed resource/Titillium_Web/TitilliumWeb-Bold.ttf
	TitilliumWebBold []byte

	//go:embed resource/Titillium_Web/TitilliumWeb-Light.ttf
	TitilliumWebLight []byte

	//go:embed resource/Ubuntu/Ubuntu-Bold.ttf
	UbuntuBold []byte

	//go:embed resource/Ubuntu/Ubuntu-Light.ttf
	UbuntuLight []byte
)
