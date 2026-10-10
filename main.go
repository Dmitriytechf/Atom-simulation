package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"atoms/internal/render"
	"atoms/internal/world"
)

const (
	screenWidth  = 800
	screenHeight = 600

	physicsDT = 1.0 / 120.0 // фиксированный шаг

	particleCount  = 800
	particleRadius = 2.0
	particleMass   = 1.0
	sigmaV         = 90.0 // пикселей/секунду — "температура" газа
)

type Game struct {
	world       *world.World
	accumulator float64
}

func NewGame() *Game {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	w := world.New(screenWidth, screenHeight, rng)

	w.SpawnParticles(world.SpawnOptions{
		Count:    particleCount,
		Radius:   particleRadius,
		Mass:     particleMass,
		SigmaV:   sigmaV,
		MarginPx: 0,
	})

	return &Game{world: w}
}

// Игровой цикл
func (g *Game) Update() error {
	g.accumulator += 1.0 / float64(ebiten.TPS())
	for g.accumulator >= physicsDT {
		g.world.Step(physicsDT)
		g.accumulator -= physicsDT
	}

	ebiten.SetWindowTitle(fmt.Sprintf(
		"Atoms — N=%d, E=%.1f",
		len(g.world.Particles), g.world.TotalKineticEnergy(),
	))
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})
	render.Draw(screen, g.world)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Atoms Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(120) // чаще тикаем физикой

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
