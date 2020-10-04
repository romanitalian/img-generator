package img

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/romanitalian/pixel.local/img-generator/pkg/colors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"strconv"
)

var (
	imgW = 300
	imgH = 300

	imgColor = "E5E5E5"
	msgColor = "AAAAAA"

	fontSize         = 0
	dpi      float64 = 72
	fontfile         = "wqy-zenhei.ttf" // TODO вынести в конфиг
	hinting          = "none"
)

// TODO добавить конструктор

func addLabel(img *image.RGBA, imgW, imgH int, msg string, msgFontSize int, msgColor colors.Hex) {
	h := font.HintingNone
	switch hinting {
	case "full":
		h = font.HintingFull
	}
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	rgb, err := colors.Hex2RGB(msgColor)
	if err != nil {
		log.Println(err)
		return
	}
	clr := color.Color(color.RGBA{R: rgb.Red, G: rgb.Blue, B: rgb.Green, A: 255})
	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(clr),
		Face: truetype.NewFace(fnt, &truetype.Options{
			Size:    float64(msgFontSize),
			DPI:     dpi,
			Hinting: h,
		}),
	}

	y := imgH/2 + msgFontSize/2 - 12
	d.Dot = fixed.Point26_6{
		X: (fixed.I(imgW) - d.MeasureString(msg)) / 2,
		Y: fixed.I(y),
	}
	d.DrawString(msg)
}

func GenerateFavicon() *bytes.Buffer {
	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	blue := color.RGBA{B: 0, A: 0}
	draw.Draw(m, m.Bounds(), &image.Uniform{C: blue}, image.ZP, draw.Src)

	var img image.Image = m
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	return buffer
}

// TODO принимать интерфейс пакета. Заполнять входные значения извне.
func Generate(urlPart []string) *bytes.Buffer {
	var err error

	msg := ""
	// TODO парсить urlPart в отдельном методе, в этом же методе заполнять структуру и возвращать её
	for i, val := range urlPart {
		switch i {
		case 1:
			if val != "" {
				imgW, err = strconv.Atoi(val)
				if err != nil {
					log.Println("Can not parse 'imgW', err: ", err)
					return nil
				}
			}
		case 2:
			if val != "" {
				imgH, err = strconv.Atoi(val)
				if err != nil {
					log.Println("Can not parse 'imgH', err: ", err)
					return nil
				}
			}
		case 3:
			if val != "" {
				imgColor = val
			}
		case 4:
			if val != "" {
				msg = val
			}
		case 5:
			if val != "" {
				msgColor = val
			}
		case 6:
			fontSize, err = strconv.Atoi(val)
			if err != nil {
				log.Println("Can not parse 'fontSize', err: ", err)
				return nil
			}
		}
	}
	if ((imgW > 0 || imgH > 0) && msg == "") || msg == "" {
		msg = fmt.Sprintf("%d x %d", imgW, imgH)
	}

	if fontSize == 0 {
		fontSize = imgW / 9
		if imgH < imgW {
			fontSize = imgH / 5
		}
	}

	hx := colors.Hex(imgColor)
	rgb, err := hx.ToRGB()
	if err != nil {
		log.Println("Can not parse 'imgColor', err: ", err)
		return nil
	}

	m := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	imgRgba := color.RGBA{R: rgb.Red, G: rgb.Green, B: rgb.Blue, A: 10}
	draw.Draw(m, m.Bounds(), &image.Uniform{C: imgRgba}, image.ZP, draw.Src)

	addLabel(m, imgW, imgH, msg, fontSize, colors.Hex(msgColor))

	var img image.Image = m
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	return buffer
}
