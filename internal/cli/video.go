package cli

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"runtime"

	"github.com/icza/mjpeg"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/teadove/awesome-fractals/internal/brot"
	"github.com/teadove/awesome-fractals/internal/palette"

	"github.com/urfave/cli/v2"
)

func drawVideo(cCtx *cli.Context) error {
	runtime.GOMAXPROCS(cCtx.Int(maxprocsFlag.Name))

	colors, ok := palette.ColorPalettes[cCtx.String(paletteFlage.Name)]
	if !ok {
		return cli.Exit("Palette not found", 1)
	}

	buildInput := brot.BuilderInput{
		Colors:       colors,
		ColorStep:    cCtx.Int(stepFlag.Name),
		XPos:         cCtx.Float64(xposFlag.Name),
		YPos:         cCtx.Float64(yposFlag.Name),
		Width:        cCtx.Int(widthFlag.Name),
		Height:       cCtx.Int(heightFlag.Name),
		MaxIteration: cCtx.Int(iterationFlag.Name),
		EscapeRadius: cCtx.Float64(radiusFlag.Name),
	}

	iterations := cCtx.Int(framesFlag.Name)
	input, _ := brot.NewInputFromBuilderInput(&buildInput)
	bar := progressbar.Default(int64(iterations))

	aw, err := mjpeg.New(cCtx.String(videoFilenameFlag.Name), int32(cCtx.Int(widthFlag.Name)), int32(cCtx.Int(heightFlag.Name)), int32(cCtx.Int(fpsFlag.Name)))
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error creating video file: %s", err.Error()), 1)
	}

	for i := 0; i < iterations; i++ {
		img := brot.Render(input)
		input.XPos -= 0.01

		buf := &bytes.Buffer{}
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Error encoding jpeg: %s", err.Error()), 1)
		}
		err = aw.AddFrame(buf.Bytes())
		if err != nil {
			return cli.Exit(fmt.Sprintf("Error adding frame: %s", err.Error()), 1)
		}

		err := bar.Add(1)
		if err != nil {
			log.Error().Err(err).Str("status", "bar.add.error")
		}
	}
	err = aw.Close()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error closing video file: %s", err.Error()), 1)
	}
	return nil
}
