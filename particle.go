package particles

import (
	"math/rand"
)

// flags is the type used to represent particle flags
type flags uint8

const (
	flags_up0 flags = 1 << iota
	flags_up1
	flags_confirmed

	flags_up flags = flags_up0 | flags_up1
)

// A quark has a state and an 'identity' for matching purposes.
type quark struct {
	up bool // the state of the quark. 'up' or 'down'
}

// A particle consists of two quarks. It has an id used for tracing
// and a slotId used for space reclamation purposes. It has some
// flags to simplfy certain tests.
//
// A particle is never directly observed - the most we know about
// a particle is what we can observe it via observations and what we
// learn once an observation is confirmed. Once confirmation
// occurs the particle is disappeared which means that no future
// observations relating to this particle can be made.
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
		if q.up {
			p.flags |= 1 << uint(i)
		}
	}

	return p
}

// confirm if the particle is a 'happy' particle. calling this
// method causes the particle to disappear, meaning that no
// further observations can be made.
func (p *particle) confirm() bool {
	if p.flags&flags_confirmed == 0 {
		p.flags |= flags_confirmed
		p.pool.disappear(p)
	}
	return p.flags&flags_up == flags_up
}

// isConfirmed answers true if the confirm() method has been called yet.
func (p *particle) isConfirmed() bool {
	return p.flags&flags_confirmed != 0
}
