package cli

import (
	"flag"
	"image/png"
	"os"
	"runtime/pprof"

	"github.com/teadove/awesome-fractals/internal/brot"
)

var service brot.Renderer

func init() {
	var input brot.Input

	flag.Float64Var(
		&input.ColorStep,
		"step",
		6000,
		"Color smooth step. Value should be greater than iteration count, otherwise the value will be adjusted to the iteration count.",
	)
	flag.IntVar(&input.Width, "width", 1000, "Rendered image width")
	flag.IntVar(&input.Height, "height", 1000, "Rendered image height")
	flag.Float64Var(
		&input.XPos,
		"xpos",
		-0.00275,
		"Point position on the real axis (defined on `x` axis)",
	)
	flag.Float64Var(
		&input.YPos,
		"ypos",
		0.78912,
		"Point position on the imaginary axis (defined on `y` axis)",
	)
	flag.Float64Var(&input.EscapeRadius, "radius", .125689, "Escape Radius")
	flag.IntVar(&input.MaxIteration, "iteration", 800, "Iteration count")
	flag.StringVar(
		&input.ColorPalette,
		"palette",
		"Hippi",
		"Hippi | Plan9 | AfternoonBlue | SummerBeach | Biochimist | Fiesta",
	)
	//  flag.StringVar(
	//	&service.OutputFile,
	//	"file",
	//	"mandelbrot.png",
	//	"The rendered mandelbrot image filname",
	//  )
	flag.Parse()
	service = *brot.New(&input)
}

func Run() {
	file, err := os.OpenFile("main.prof", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	err = pprof.StartCPUProfile(file)
	if err != nil {
		panic(err)
	}

	//  done := make(chan struct{})
	//  iterations := 2000
	//
	//  go func() {
	//	bar := progressbar.Default(int64(iterations))
	//	for i := 0; i <= iterations; i++ {
	//		<-done
	//		err := bar.Add(1)
	//		if err != nil {
	//			println(err.Error())
	//		}
	//	}
	//  }()

	image := service.Render()

	//// TODO add err check
	output, _ := os.Create("file.png")
	//// TODO add err check
	_ = png.Encode(output, image)

	pprof.StopCPUProfile()
}
