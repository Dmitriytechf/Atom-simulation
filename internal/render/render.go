package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"atoms/internal/world"
)

var particleColor = color.RGBA{R: 120, G: 220, B: 255, A: 255}

// Draw рисует все частицы мира на экране
func Draw(screen *ebiten.Image, w *world.World) {
	for i := range w.Particles {
		p := &w.Particles[i]
		vector.FillCircle(
			screen,
			float32(p.X),
			float32(p.Y),
			float32(p.Radius),
			particleColor,
			true, // antialias
		)
	}
}
