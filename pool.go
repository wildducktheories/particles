package particles

import (
	//"log"
	"math/rand"
)

// A pool represents a known number of particles with evenly
// distributed number of quarks. The pool also contains 'sad'
// particles, but these particles are not observable by the
// sampling methods.
type pool struct {
	particles []*particle
	next      Id
	sync      chan func()
	returned  chan *particle
}

// Return a new Pool of the specified size.
func NewPool(size int) Pool {
	particles := make([]*particle, size, size)
	pool := &pool{
		particles: particles,
		next:      Id(size),
		returned:  make(chan *particle),
		sync:      make(chan func()),
	}
	for i, _ := range particles {
		particles[i] = newParticle(Id(i), i, [2]bool{flip(), flip()}, pool)
	}
	go func() {
		for {
			f := <-pool.sync
			if f == nil {
				return
			}
			f()
		}
	}()
	return pool
}

// SampleQuark samples all the up quarks in the pool and
// returns an observation consisting of itself and the
// parent particle.
//
// Observations returned by this method have a 1/2
// probability of describing a happy particle.
func (p *pool) SampleQuark() Observation {
	for {
		px := rand.Int31n(int32(len(p.particles) * 2))
		qx := px & 1
		px >>= 1
		po := p.particles[px]
		qo := po.quarks[qx]
		if qo.up {
			//log.Printf("sample quark:  %d[%d]", po.id, qx)
			return &observation{
				particle: po,
				quark:    qo,
				process:  quarkProcess,
			}
		}
	}
}

// SampleParticle samples all the particles in the pool
// selects one of the particle's up-quarks and returns
// an Observation describing both.
//
// Observations returned by this method have a 1/3
// probability of describing a happy particle.
func (p *pool) SampleParticle() Observation {
	for {
		px := rand.Int31n(int32(len(p.particles)))
		po := p.particles[px]

		var q *quark

		switch po.flags & flags_up {
		case flags_up:
			if flip() {
				q = po.quarks[1]
			} else {
				q = po.quarks[0]
			}
		case flags_up1:
			q = po.quarks[1]
		case flags_up0:
			q = po.quarks[0]
		default:
			continue
		}
		//log.Printf("sample particle:  %d", po.id)

		return &observation{
			particle: po,
			quark:    q,
			process:  particleProcess,
		}
	}
}

// Close shuts down the async goroutine.
func (p *pool) Close() {
	close(p.returned)
	close(p.sync)
}

// disappear replaces an existing particle with a new particle
// with matching quark-states. Existing observations can never match
// the new particle.
func (p *pool) disappear(po *particle) {
	done := make(chan struct{})
	p.sync <- func() {
		if p.particles[po.slotId] == po {
			// replace the particle with one that has the same quark states
			// so that we don't change the composition of the pool
			p.particles[po.slotId] = newParticle(Id(p.next), po.slotId, [2]bool{po.quarks[0].up, po.quarks[1].up}, p)
			p.next++
		}
		done <- struct{}{}
	}
	<-done
}
