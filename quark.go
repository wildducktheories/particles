package particles

// A quark is a component of a particle and is either 'up' or 'down'.
type quark struct {
	up       bool      // the state of the quark. 'up' or 'down'
	particle *particle // the particle this quark is a component of
	index    int       // the index of the quark within the particle
}
