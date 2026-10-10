package world

import (
	"math/rand"

	"atoms/internal/physics"
)

// World — прямоугольная область с частицами и упругими стенками.
type World struct {
	Width, Height float64
	Particles     []physics.Particle
	Restitution   float64 // 1.0 = абсолютно упругий удар, <1 = с потерями
	Rand          *rand.Rand
}

type SpawnOptions struct {
	Count    int
	Radius   float64
	Mass     float64
	SigmaV   float64 // стандартное отклонение каждой компоненты скорости
	MarginPx float64 // отступ от краев мира
}

// New создаёт пустой мир заданного размера.
func New(width, height float64, rng *rand.Rand) *World {
	if rng == nil {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	return &World{
		Width:       width,
		Height:      height,
		Restitution: 1.0,
		Rand:        rng,
	}
}

// Step продвигает симуляцию на dt секунд
func (w *World) Step(dt float64) {
	for i := range w.Particles {
		p := &w.Particles[i]
		p.Step(dt)
		w.resolveWallCollision(p)
	}
}

// SpawnParticles добавляет в мир Count новых частиц со случайными
// позициями и скоростями из нормального распределения (Максвелл).
// Существующие частицы не трогает — можно вызывать многократно.
func (w *World) SpawnParticles(opts SpawnOptions) {
	// Область, в которой разрешено рождаться центрам частиц.
	// Отступаем MarginPx от стенки и ещё Radius, чтобы круг не вылез.
	minX := opts.MarginPx + opts.Radius
	minY := opts.MarginPx + opts.Radius
	maxX := w.Width - opts.MarginPx - opts.Radius
	maxY := w.Height - opts.MarginPx - opts.Radius

	if maxX <= minX || maxY <= minY {
		return // мир слишком маленький — молча ничего не делаем
	}

	for i := 0; i < opts.Count; i++ {
		p := physics.Particle{
			X:      minX + w.Rand.Float64()*(maxX-minX),
			Y:      minY + w.Rand.Float64()*(maxY-minY),
			VX:     w.Rand.NormFloat64() * opts.SigmaV,
			VY:     w.Rand.NormFloat64() * opts.SigmaV,
			Radius: opts.Radius,
			Mass:   opts.Mass,
		}
		w.Particles = append(w.Particles, p)
	}
}

// resolveWallCollision отражает частицу от стенки, если она её пересекла.
// Коэффициент Restitution задаёт упругость
func (w *World) resolveWallCollision(p *physics.Particle) {
	// Левая стенка
	if p.X-p.Radius < 0 {
		p.X = p.Radius
		p.VX = -p.VX * w.Restitution
	}
	// Правая стенка
	if p.X+p.Radius > w.Width {
		p.X = w.Width - p.Radius
		p.VX = -p.VX * w.Restitution
	}
	// Верхняя стенка
	if p.Y-p.Radius < 0 {
		p.Y = p.Radius
		p.VY = -p.VY * w.Restitution
	}
	// Нижняя стенка
	if p.Y+p.Radius > w.Height {
		p.Y = w.Height - p.Radius
		p.VY = -p.VY * w.Restitution
	}
}

// TotalKineticEnergy — сумма кинетических энергий всех частиц.
func (w *World) TotalKineticEnergy() float64 {
	var sum float64
	for i := range w.Particles {
		sum += w.Particles[i].KineticEnergy()
	}
	return sum
}
