package img

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/freetype/truetype"
	"github.com/romanitalian/img-generate/v2/pkg/logger"
	"golang.org/x/image/font"
)

const (
	// Font paths
	fontDir     = "assets/fonts"
	defaultFont = "wqy-zenhei.ttf"

	// Font settings
	dpiDefault float64 = 72
)

var (
	// Font cache
	fontCache *truetype.Font
	log       = logger.Get()
)

// loadFont loads font from file and caches it
func loadFont(ctx context.Context) (*truetype.Font, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if fontCache != nil {
			log.Debug().Msg("using cached font")
			return fontCache, nil
		}

		// Try to load from assets directory first
		fontPath := filepath.Join(fontDir, defaultFont)
		log.Debug().Str("path", fontPath).Msg("trying to load font from assets")

		fontBytes, err := os.ReadFile(fontPath)
		if err != nil {
			log.Warn().Err(err).Str("path", fontPath).Msg("failed to load font from assets, trying system fonts")

			// Try to load from system fonts
			systemFonts := []string{
				"/System/Library/Fonts/Supplemental/Arial Unicode.ttf", // macOS
				"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",      // Linux
				"C:\\Windows\\Fonts\\arial.ttf",                        // Windows
			}

			for _, path := range systemFonts {
				if _, err := os.Stat(path); err == nil {
					log.Debug().Str("path", path).Msg("trying to load system font")
					fontBytes, err = os.ReadFile(path)
					if err == nil {
						log.Info().Str("path", path).Msg("successfully loaded system font")
						break
					}
				}
			}

			if fontBytes == nil {
				err := fmt.Errorf("failed to load font: %v", err)
				log.Error().Err(err).Msg("all font loading attempts failed")
				return nil, err
			}
		} else {
			log.Info().Str("path", fontPath).Msg("successfully loaded font from assets")
		}

		fontCache, err = truetype.Parse(fontBytes)
		if err != nil {
			err := fmt.Errorf("failed to parse font: %v", err)
			log.Error().Err(err).Msg("failed to parse font data")
			return nil, err
		}

		return fontCache, nil
	}
}

// createFontFace creates a new font face with given size
func createFontFace(ctx context.Context, size float64) (font.Face, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		fnt, err := loadFont(ctx)
		if err != nil {
			return nil, err
		}

		log.Debug().Float64("size", size).Msg("creating new font face")
		return truetype.NewFace(fnt, &truetype.Options{
			Size:    size,
			DPI:     dpiDefault,
			Hinting: font.HintingNone,
		}), nil
	}
}
