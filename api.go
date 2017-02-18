package particles

type Id uint64

// A Pool is a pool of particles.
//
// Particles can never be observed directly only by sampling which produces an Observation.
type Pool interface {
	// SampleQuark randomly selects an up quark from all the quarks in the pool and
	// returns an observation that describes the quark and the particle to which is bound.
	//
	// The probability of an Observation returned by this method is 1/2.
	SampleQuark() Observation

	// SampleParticle randomly selects a particle from all the and returns an Observation
	// that describes the particle and one of its up quarks.
	//
	// The probability of an Observation returned by this method is 1/3.
	SampleParticle() Observation

	// Release the resources associated with the pool.
	Close()
}

// An Observation describes a particle and one its up quarks.
//
// The probability of an observation is the probability both quarks of the
// particle are up quarks or, alternatively, that the Read() method will return
// true when called. Its value always between 0 and 1.0 until the Read()
// method is called at which point it becomes 0 or 1.
type Observation interface {
	Probability() float64
	Confirm() bool
}

// A Matcher matches two observations and generates a new observation or returns nil
type Matcher func(Observation, Observation) Observation

type MatcherProcess func(<-chan Observation, <-chan Observation, chan<- Observation)

// An Observer generates a single observation.
type Observer func() Observation

// A process that generates observations.
type ObserverProcess func(<-chan struct{}, chan<- Observation)

// ProcessFromMatcher returns a function that can match Observations read from two channels
// and generates a new Observation on third channel.
func ProcessFromMatcher(m Matcher) MatcherProcess {
	return func(inA <-chan Observation, inB <-chan Observation, outC chan<- Observation) {
		defer close(outC)
		for obsA := range inA {
			obsB := <-inB
			if obsB == nil {
				return
			}
			obsC := m(obsA, obsB)
			if obsC != nil {
				outC <- obsC
			}
		}
	}
}

// ProcessFromObserver returns a function that calls an observer function to generate
// new Observations until the input channel is closed.
func ProcessFromObserver(o Observer) ObserverProcess {
	return func(in <-chan struct{}, out chan<- Observation) {
		defer close(out)
		for {
			select {
			case <-in:
				return
			default:
				out <- o()
			}
		}
	}
}

var Verbose bool
