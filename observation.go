package particles

//import "log"

// An observation of a particle and one of its quarks.
//
// A given observation can only observe one quark of a particle at
// a time and then only if it is an up quark. Down quarks cannot
// be directly observed. The existence or non-existence of down
// quark can only be confirmed after an observation is resolved.
//
// It is associated with a probability which is a number
// between 0 and 1 that the other quark is also an up quark.
type observation struct {
	particle *particle // the observed particle
	quark    *quark    // the observed quark
	process  *process  // the process that generated the observation
}

// Probability returns that probability that Read() will return true.
func (o *observation) Probability() float64 {
	if !o.particle.isRead() {
		return o.process.probability
	} else if o.particle.read() {
		return 1.0
	} else {
		return 0
	}
}

// Read reads the state of the particle's other quark and returns true if it in the up state.
func (o *observation) Read() bool {
	return o.particle.read()
}

// Matches two Observations on their particles (only)
func ParticleMatcher(a Observation, b Observation) Observation {
	aO := a.(*observation)
	bO := b.(*observation)

	if aO.particle == bO.particle {
		p := matchingProcess(aO.process, bO.process, false)

		cO := &observation{
			particle: aO.particle,
			quark:    aO.quark,
			process:  p,
		}
		if flip() {
			cO.quark = bO.quark
		}
		return cO
	} else {
		return nil
	}
}

// Matches two Observations on their quarks (only)
func QuarkMatcher(a Observation, b Observation) Observation {
	aO := a.(*observation)
	bO := b.(*observation)

	if aO.quark == bO.quark {
		p := matchingProcess(aO.process, bO.process, true)

		cO := &observation{
			particle: aO.particle,
			quark:    aO.quark,
			process:  p,
		}
		return cO
	} else {
		return nil
	}
}
