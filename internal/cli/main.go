package cli

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/teadove/awesome-fractals/internal/brot"
)

var service brot.Service

func init() {
	service.WG = &sync.WaitGroup{}

	flag.Float64Var(
		&service.ColorStep,
		"step",
		6000,
		"Color smooth step. Value should be greater than iteration count, otherwise the value will be adjusted to the iteration count.",
	)
	flag.IntVar(&service.Width, "width", 1000, "Rendered image width")
	flag.IntVar(&service.Height, "height", 1000, "Rendered image height")
	flag.Float64Var(
		&service.XPos,
		"xpos",
		-0.00275,
		"Point position on the real axis (defined on `x` axis)",
	)
	flag.Float64Var(
		&service.YPos,
		"ypos",
		0.78912,
		"Point position on the imaginary axis (defined on `y` axis)",
	)
	flag.Float64Var(&service.EscapeRadius, "radius", .125689, "Escape Radius")
	flag.IntVar(&service.MaxIteration, "iteration", 800, "Iteration count")
	flag.StringVar(
		&service.ColorPalette,
		"palette",
		"Hippi",
		"Hippi | Plan9 | AfternoonBlue | SummerBeach | Biochimist | Fiesta",
	)
	flag.StringVar(
		&service.OutputFile,
		"file",
		"mandelbrot.png",
		"The rendered mandelbrot image filname",
	)
	flag.Parse()
}

func Run() {
	done := make(chan struct{})
	ticker := time.NewTicker(time.Millisecond * 100)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Print(".")
			case <-done:
				ticker.Stop()
				fmt.Printf("\n\nMandelbrot set rendered into `%s`\n", service.OutputFile)
			}
		}
	}()

	service.Run(done)
}
