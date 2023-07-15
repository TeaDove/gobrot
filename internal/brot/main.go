package brot

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/teadove/awesome-fractals/internal/palette"
)

type Renderer struct {
	xpos, ypos    float64
	width, height int
	maxIteration  int
	EscapeRadius  float64
	colors        []color.RGBA
}

type Input struct {
	ColorPalette  string
	ColorStep     int
	XPos          float64
	YPos          float64
	Width, Height int
	MaxIteration  int
	EscapeRadius  float64
}

func New(input *Input) (*Renderer, int) {
	r := Renderer{
		xpos:         input.XPos,
		ypos:         input.YPos,
		width:        input.Width,
		height:       input.Height,
		maxIteration: input.MaxIteration,
		EscapeRadius: input.EscapeRadius,
	}

	if input.ColorStep < input.MaxIteration {
		input.ColorStep = input.MaxIteration
	}
	r.colors = input.InterpolateColors()

	return &r, r.height
}

func (s *Input) InterpolateColors() []color.RGBA {
	colors := palette.ColorPalettes[s.ColorPalette]

	var factor float64
	var steps []float64
	var cols []uint32
	var interpolated []uint32
	var interpolatedColors []color.RGBA

	factor = 1.0 / float64(s.ColorStep)
	for index, col := range colors {
		if col.Step == 0.0 && index != 0 {
			stepRatio := float64(index+1) / float64(len(colors))
			step := float64(int(stepRatio*100)) / 100 // truncate to 2 decimal precision
			steps = append(steps, step)
		} else {
			steps = append(steps, col.Step)
		}
		r, g, b, a := col.Color.RGBA()
		r /= 0xff
		g /= 0xff
		b /= 0xff
		a /= 0xff
		uintColor := r<<24 | g<<16 | b<<8 | a
		cols = append(cols, uintColor)
	}

	var min, max, minColor, maxColor float64
	for i := 0.0; i <= 1; i += factor {
		for j := 0; j < len(colors)-1; j++ {
			if !(i >= steps[j] && i < steps[j+1]) {
				continue
			}
			min = steps[j]
			max = steps[j+1]
			minColor = float64(cols[j])
			maxColor = float64(cols[j+1])
			uintColor := cosineInterpolation(
				maxColor,
				minColor,
				(i-min)/(max-min),
			)
			interpolated = append(interpolated, uint32(uintColor))
		}
	}

	for _, pixelValue := range interpolated {
		r := pixelValue >> 24 & 0xff
		g := pixelValue >> 16 & 0xff
		b := pixelValue >> 8 & 0xff
		a := 0xff

		interpolatedColors = append(
			interpolatedColors,
			color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)},
		)
	}

	return interpolatedColors
}

func (s *Renderer) Render(done chan struct{}) *image.RGBA {
	ratio := float64(s.height) / float64(s.width)
	xmin, xmax := s.xpos-s.EscapeRadius/2.0, math.Abs(s.xpos+s.EscapeRadius/2.0)
	ymin, ymax := s.ypos-s.EscapeRadius*ratio/2.0, math.Abs(s.ypos+s.EscapeRadius*ratio/2.0)

	rgbaImageComplied := image.NewRGBA(
		image.Rectangle{Min: image.Point{}, Max: image.Point{X: s.width, Y: s.height}},
	)

	wg := sync.WaitGroup{}
	for iy := 0; iy < s.height; iy++ {
		wg.Add(1)
		go func(iy int) {
			defer wg.Done()
			defer func() {
				if done != nil {
					done <- struct{}{}
				}
			}()

			for ix := 0; ix < s.width; ix++ {
				x := xmin + (xmax-xmin)*float64(ix)/float64(s.width-1)
				y := ymin + (ymax-ymin)*float64(iy)/float64(s.width-1)
				norm, it := mandelIteration(x, y, s.maxIteration)
				iteration := float64(s.maxIteration-it) + math.Log(norm)

				if int(math.Abs(iteration)) >= len(s.colors)-1 {
					continue
				}
				color1 := s.colors[int(math.Abs(iteration))]
				color2 := s.colors[int(math.Abs(iteration))+1]
				compiledColor := linearInterpolation(
					rgbaToUint(color1),
					rgbaToUint(color2),
					uint32(iteration),
				)

				rgbaImageComplied.Set(ix, iy, uint32ToRgba(compiledColor))
			}
		}(iy)
	}

	wg.Wait()
	return rgbaImageComplied
}

func cosineInterpolation(c1, c2, mu float64) float64 {
	mu2 := (1 - math.Cos(mu*math.Pi)) / 2.0
	return c1*(1-mu2) + c2*mu2
}

func linearInterpolation(c1, c2, mu uint32) uint32 {
	return c1*(1-mu) + c2*mu
}

func mandelIteration(cx, cy float64, maxIter int) (float64, int) {
	x, y, xx, yy := 0.0, 0.0, 0.0, 0.0

	for i := 0; i < maxIter; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + cx
		y = 2*xy + cy
	}

	logZn := (x*x + y*y) / 2
	return logZn, maxIter
}

func rgbaToUint(color color.RGBA) uint32 {
	r, g, b, a := color.RGBA()
	r /= 0xff
	g /= 0xff
	b /= 0xff
	a /= 0xff
	return r<<24 | g<<16 | b<<8 | a
}

func uint32ToRgba(col uint32) color.RGBA {
	r := col >> 24 & 0xff
	g := col >> 16 & 0xff
	b := col >> 8 & 0xff
	a := 0xff
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}
