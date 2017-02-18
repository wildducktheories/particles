package particles

import (
	"math/rand"
)

type flags uint8

const (
	flags_up0 flags = 1 << iota
	flags_up1
	flags_read

	flags_up flags = flags_up0 | flags_up1
)

// A particle consists of two quarks. It has an id used for tracing
// and a slotId used for space reclamation purposes. It has some
// flags to simplfy certain tests.
//
// A particle is never directly observed - the most we know about
// a particle is what we can observe via observations and what we
// learn once an observation is read. Once this occurs the
// underlying particle ceases to be directly observable.
type particle struct {
	quarks [2]*quark // the two quarks of the particle

	id     Id    // for tracing, never re-used
	slotId int   // for reclaiming space used by unobservable particles
	pool   *pool // the pool from which the particle was sampled.
	flags  flags // bit 0 = quark[0].up, bit 1 = quark[1].up
}

// flip a proverbial coin.
func flip() bool {
	if rand.Float64() < 0.50 {
		return false
	} else {
		return true
	}
}

// newParticle returns a new particle with the given id and slot id.
func newParticle(id Id, slotId int, states [2]bool, pool *pool) *particle {

	p := &particle{
		id:     id,
		slotId: slotId,
		pool:   pool,
	}

	for i, _ := range p.quarks {
		q := &quark{}
		p.quarks[i] = q
		q.up = states[i]
		q.index = i
		q.particle = p
		if q.up {
			p.flags |= 1 << uint(i)
		}
	}

	return p
}

// read the particle and render unobservable in the future
func (p *particle) read() bool {
	if p.flags&flags_read == 0 {
		p.flags |= flags_read
		p.pool.returnToPool(p)
	}
	return p.flags&flags_up == flags_up
}

// answer if the particle has been read yet
func (p *particle) isRead() bool {
	return p.flags&flags_read != 0
}
