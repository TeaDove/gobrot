package cli

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/teadove/awesome-fractals/internal/palette"
	"github.com/urfave/cli/v2"
)

var (
	widthFlag  = &cli.IntFlag{Name: "width", Usage: "Rendered image width", Value: 1_000}
	heightFlag = &cli.IntFlag{Name: "height", Usage: "Rendered image height", Value: 1_000}
	stepFlag   = &cli.IntFlag{
		Name: "step",
		Usage: "Color smooth step. " +
			"Value should be greater than iteration count, " +
			"otherwise the value will be adjusted to the iteration count.",
		Value: 6_000,
	}
	xposFlag = &cli.Float64Flag{
		Name:  "xpos",
		Usage: "Point position on the real axis (defined on `x` axis)",
		Value: -0.00275,
	}
	yposFlag = &cli.Float64Flag{
		Name:  "ypos",
		Usage: "Point position on the imaginary axis (defined on `y` axis)",
		Value: 0.78912,
	}
	radiusFlag    = &cli.Float64Flag{Name: "radius", Usage: "Escape Radius", Value: .125689}
	iterationFlag = &cli.IntFlag{Name: "iteration", Value: 800, Usage: "Iteration count"}
	paletteFlage  = &cli.StringFlag{
		Name:  "palette",
		Value: "Hippi",
		Usage: strings.Join(palette.GetPaletteNames(), " | "),
	}
	imageFilenameFlag = &cli.StringFlag{
		Name:  "filename",
		Value: "fractal.png",
		Usage: "Path to save image",
	}

	maxprocsFlag = &cli.IntFlag{
		Name:  "maxprocs",
		Usage: "max amount of processes to use, by default is amount of cores in CPU munis one",
		Value: runtime.NumCPU() - 1,
	}
)

func Run() {
	captureInterrupt()

	app := &cli.App{Flags: []cli.Flag{maxprocsFlag},
		Commands: []*cli.Command{{
			Name:   "image",
			Action: drawImage,
			Flags: []cli.Flag{
				imageFilenameFlag,
				widthFlag,
				heightFlag,
				stepFlag,
				xposFlag,
				yposFlag,
				radiusFlag,
				iterationFlag,
				paletteFlage,
			}},
			{
				Name:   "video",
				Action: drawVideo,
				Flags: []cli.Flag{
					imageFilenameFlag,
					widthFlag,
					heightFlag,
					stepFlag,
					xposFlag,
					yposFlag,
					radiusFlag,
					iterationFlag,
					paletteFlage,
				},
			},
		}}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	//  done := make(chan struct{})
	//  iterations := 10
	//
	//  go func() {

	//  video, err := vidio.NewVideoWriter("video.mp4", 800, 800, &vidio.Options{})
	//  if err != nil {
	//	panic(err)
	//  }
	//
	//  for i := 0; i < 10; i++ {
	//	img := service.Render()
	//	service.EscapeRadius -= 0.01
	//
	//	buf := new(bytes.Buffer)
	//	err := jpeg.Encode(buf, img, nil)
	//	if err != nil {
	//		panic(err)
	//	}
	//	sendS3 := buf.Bytes()
	//
	//	err = video.Write(sendS3)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	// TODO add err check
	//	output, _ := os.Create(fmt.Sprintf("file_%d.png", i))
	//	//// TODO add err check
	//	_ = png.Encode(output, img)
	//
	//	err = bar.Add(1)
	//	if err != nil {
	//		println(err.Error())
	//	}
	//  }

}
