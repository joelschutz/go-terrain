package pkg

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Particle struct {
	position, velocity mgl32.Vec2

	volume, sediment, density float32
	Stopped                   bool
}

func (p Particle) GetPosition() mgl32.Vec2 {
	return p.position
}

func (p Particle) X() float32 {
	return p.position.X()
}

func (p Particle) Y() float32 {
	return p.position.Y()
}
func (p *Particle) Drop(hm *HeightMap, dt, friction, depositionRate, evaporationRate float32) {
	hmSize := hm.Bounds().Dx()
	x, y := int(p.position.X()), int(p.position.Y())
	if x < hmSize-1 && y < hmSize-1 && x > 0 && y > 0 {
		n := CalculateNormal(x, y, *hm, 20)

		p.velocity = p.velocity.Add(mgl32.Vec2{float32(n.X()), float32(n.Y())}.Mul(dt / (p.volume * p.density)))
		p.position = p.position.Add(p.velocity.Mul(dt))
		p.velocity = p.velocity.Mul(1 - dt*friction)
		if p.velocity.Len() < .1 {
			p.Stopped = true
			return
		}
		nX := int(p.position.X())
		nY := int(p.position.Y())
		if nX > hmSize-1 || nY > hmSize-1 || nX < 0 || nY < 0 {
			p.Stopped = true
			return
		}
		hDiff := (*hm)[x][y] - (*hm)[nX][nY]
		cEq := p.volume * p.velocity.Len() * float32(hDiff)
		if cEq < 0 {
			cEq = 0
		}

		cDiff := cEq - p.sediment
		depo := dt * depositionRate * cDiff
		p.sediment += depo
		hm.Set(x, y, hm.Get(x, y)-float64(depo*p.volume))
		p.volume *= 1 - dt*evaporationRate
		if p.volume < .1 {
			p.Stopped = true
			return
		}
	} else {
		p.Stopped = true
	}

}

func NewParticle(pos, vel mgl32.Vec2, vol, sed, den float32) *Particle {
	return &Particle{
		position: pos,
		velocity: vel,
		volume:   vol,
		sediment: sed,
		density:  den,
	}
}
