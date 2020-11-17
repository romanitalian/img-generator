package img

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/romanitalian/img-generate/pkg/colors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"strconv"
)

const (
	imgColorDefault = "E5E5E5"
	msgColorDefault = "AAAAAA"
	imgWDefault     = 300
	imgHDefault     = 300

	fontSizeDefault         = 0
	dpiDefault      float64 = 72
	fontfileDefault         = "wqy-zenhei.ttf"
)

type Label struct {
	Text     string
	FontSize int
	Color    string
}

type Img struct {
	Width  int
	Height int
	Color  string
	Label
}

func GenerateFavicon() (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	m := image.NewRGBA(image.Rect(0, 0, 16, 16))
	clr := color.RGBA{B: 0, A: 0}
	draw.Draw(m, m.Bounds(), &image.Uniform{C: clr}, image.ZP, draw.Src)

	var img image.Image = m
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		return nil, err
	}

	return buffer, nil
}

// Generate - return picture in bytes (actually in bytes.Buffer for write it in http response).
// params []string - слайс из ширины, высоты и т.д ( /ШИРИНА/ВЫСОТА/ЦВЕТ/ТЕКСТ/ЦВЕТ_ТЕКСТА/РАЗМЕР_ШРИФТА).
func Generate(urlPart []string) (*bytes.Buffer, error) {
	var (
		err      error
		imgW     = imgWDefault
		imgH     = imgHDefault
		imgColor = imgColorDefault
		msg      = ""
		msgColor = msgColorDefault
		fontSize = fontSizeDefault
	)
	for i, val := range urlPart {
		switch i {
		case 1:
			if val != "" {
				imgW, err = strconv.Atoi(val)
				if err != nil {
					return nil, err
				}
			}
		case 2:
			if val != "" {
				imgH, err = strconv.Atoi(val)
				if err != nil {
					return nil, err
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
				return nil, err
			}
		}
	}
	// Соберём структуру Текста
	label := Label{Text: msg, FontSize: fontSize, Color: msgColor}
	// Соберём структуру Картинки с нужными полями - высота, ширина, цвет и текст
	img := Img{Width: imgW, Height: imgH, Color: imgColor, Label: label}

	// Сгенерим нашу картинку с текстом
	return img.generate()
}


// generate - соберёт картинку по нужным размерам, цветом и текстом.
func (i Img) generate() (*bytes.Buffer, error) {
	// Если есть размеры и нет требований по Тексту - соберём Текст по умолчанию.
	if ((i.Width > 0 || i.Height > 0) && i.Text == "") || i.Text == "" {
		i.Text = fmt.Sprintf("%d x %d", i.Width, i.Height)
	}
	// Если нет требований по размеру шрифта - подберём его исходя из размеров картинки.
	if i.FontSize == 0 {
		i.FontSize = i.Width / 10
		if i.Height < i.Width {
			i.FontSize = i.Height / 5
		}
	}
	// Переведём цвет из строки в color.RGBA.
	clr, err := colors.ToRGBA(i.Color)
	if err != nil {
		return nil, err
	}

	// Создадим in-memory картинку с нужными размерами.
	m := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))
	// Отрисуем картинку:
	// - по размерам (Bounds)
	// - и с цветом (Uniform - обёртка над color.Color c Image функциями)
	// - исходя из точки (Point), как базовой картинки
	// - заполним цветом нашу Uniform (draw.Src)
	draw.Draw(m, m.Bounds(), image.NewUniform(clr), image.Point{}, draw.Src)
	// Добавим текст в картинку.
	if err = i.drawLabel(m); err != nil {
		return nil, err
	}
	var im image.Image = m
	// Выделим память под нашу данные (байты картинки).
	buffer := &bytes.Buffer{}
	// Закодируем картинку в нашу аллоцированную память.
	err = jpeg.Encode(buffer, im, nil)

	return buffer, err
}

// drawLabel - добавит текст на картинку.
func (i *Img) drawLabel(m *image.RGBA) error {
	// Разберём цвет текста из строки в RGBA.
	clr, err := colors.ToRGBA(i.Label.Color)
	if err != nil {
		return err
	}
	// Получим шрифт (должен работать и с латиницей и с кириллицей).
	fontBytes, err := ioutil.ReadFile(fontfileDefault)
	if err != nil {
		return err
	}
	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	// Подготовим Drawer для отрисовки текста на картинке.
	d := &font.Drawer{
		Dst: m,
		Src: image.NewUniform(clr),
		Face: truetype.NewFace(fnt, &truetype.Options{
			Size:    float64(i.FontSize),
			DPI:     dpiDefault,
			Hinting: font.HintingNone,
		}),
	}
	// Зададим базовую линию.
	d.Dot = fixed.Point26_6{
		X: (fixed.I(i.Width) - d.MeasureString(i.Text)) / 2,
		Y: fixed.I((i.Height+i.FontSize)/2 - 12),
	}
	// Непосредственно отрисовка текста в нашу RGBA картинку.
	d.DrawString(i.Text)

	return nil
}
