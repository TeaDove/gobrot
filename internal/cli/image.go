package cli

import (
	"image/png"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"

	"github.com/teadove/awesome-fractals/internal/brot"
	"github.com/teadove/awesome-fractals/internal/palette"
	"github.com/urfave/cli/v2"
)

func drawImage(cCtx *cli.Context) error {
	_, ok := palette.ColorPalettes[cCtx.String(paletteFlage.Name)]
	if !ok {
		return cli.Exit("Palette not found", 1)
	}

	input := brot.Input{
		ColorPalette: cCtx.String(paletteFlage.Name),
		ColorStep:    cCtx.Int(stepFlag.Name),
		XPos:         cCtx.Float64(xposFlag.Name),
		YPos:         cCtx.Float64(yposFlag.Name),
		Width:        cCtx.Int(widthFlag.Name),
		Height:       cCtx.Int(heightFlag.Name),
		MaxIteration: cCtx.Int(iterationFlag.Name),
		EscapeRadius: cCtx.Float64(radiusFlag.Name),
	}

	service, iterations := brot.New(&input)
	done := make(chan struct{}, iterations)
	go func() {
		bar := progressbar.Default(int64(iterations))
		for i := 0; i < iterations; i++ {
			<-done
			err := bar.Add(1)
			if err != nil {
				log.Error().Err(err).Str("status", "bar.add.error")
			}
		}
	}()
	img := service.Render(done)

	output, err := os.Create(cCtx.String(imageFilenameFlag.Name))
	if err != nil {
		return err
	}
	err = png.Encode(output, img)
	if err != nil {
		return err
	}

	return nil
}
