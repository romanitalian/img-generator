package colors

import (
	"image/color"
	"strconv"
)

type rgb struct {
	red   uint8
	green uint8
	blue  uint8
}

func ToRGBA(h string) (color.RGBA, error) {
	rgb, err := hex2RGB(h)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{R: rgb.red, G: rgb.green, B: rgb.blue, A: 255}, nil
}

func hex2RGB(hex string) (rgb, error) {
	values, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return rgb{}, err
	}

	return rgb{
		red:   uint8(values >> 16),
		green: uint8((values >> 8) & 0xFF),
		blue:  uint8(values & 0xFF),
	}, nil
}
