package physics

type Particle struct {
	X, Y   float64 // позиция центра
	VX, VY float64 // скорость
	Radius float64
	Mass   float64
}

// Step интегрирует движение частицы за dt секунд.
func (p *Particle) Step(dt float64) {
	p.X += p.VX * dt
	p.Y += p.VY * dt
}

// KineticEnergy возвращает кинетическую энергию: ½·m·|v|²
func (p *Particle) KineticEnergy() float64 {
	return 0.5 * p.Mass * (p.VX*p.VX + p.VY*p.VY)
}
