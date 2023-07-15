package brot

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teadove/awesome-fractals/internal/utils"
)

func TestUnit_Brot_Render_Ok(t *testing.T) {
	file, err := os.ReadFile("test_image.png")
	utils.Check(err)

	image, err := png.Decode(bytes.NewReader(file))
	utils.Check(err)

	renderer, _ := New(&Input{
		ColorPalette: "Plan9",
		ColorStep:    6000,
		XPos:         -0.00275,
		YPos:         0.78912,
		Width:        1000,
		Height:       1000,
		MaxIteration: 800,
		EscapeRadius: .125689,
	})
	imageCompile := renderer.Render(nil)

	assert.Equal(t, image, imageCompile)
}
