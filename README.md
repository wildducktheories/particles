# NAME

particles - a library and utility for simulating a simple particle universe

# DESCRIPTION

**particles** implements a simulation of a toy-universe.

This universe contains particles which have 2 quarks. Quarks exist in either an
up state (up-quarks) or a down state (down-quarks). Particles
with one up and and one down quark are called 'calm' particles. Particles
with two up-quarks are called 'happy' particles. Particles with two down-quarks
would be called 'sad' particles, but such particles have never been observed.
In nature, the ratio of 'calm' particles to 'happy' particles is 2:1.

Only one quark of a particle is directly observable in a given observation and
the only observable quarks are up-quarks. An observation can be "confirmed" to
determine the state of the unobserved quark. Once this happens the associated
particle disappears and it is replaced by a particle with the same quark-state
somewhere else close by - the replacement particle has no other discernible
relationship to the replaced particle. This behaviour has annoyed physicists
for a long time, because they haven't been able to isolate a pool of pure 'happy'
particles to understand how they behave in an environment that unconstrained
by neighbouring 'calm' particles.

Particles can be pooled. Pools can be sampled by random sampling processes.
Sampling processes choose a particle or quark at random and then emit an observation
which encodes details of the particle and one related up-quark (sample-particle)
or the quark and its related particle (sample-quark) in an observation.

The observation also records the probability that the observation is of a
'happy' particle. The sample-particle process generates observations
that have a 1/3 chance of identifying a 'happy' particle. The sample-quark
process generates observations that have a 1/2 chance of identifying a 'happy'
particle.

Different observations of the same particle generated by the same or different
processes can be matched. Matching can occur in one of two ways - matching
by particle or matching by quark.

Matching must be done by sensitive equipment that does not reveal the
identity of either the quark or particle being matched but, if matching
is successful, does emit a new observation for the matched particle and
one of its up-quarks.

# QUESTIONS

**Q1.** Suppose you match two observations of the same particle, one from the
sample-quark process, one from the sample-particle process, each asserting a
different probability, 1/2 and 1/3 respectively, that the observed
particle is a 'happy' particle.

Given the contradictory assertions made by the matching observations,
which observation, if any, do you believe?

**Q2.** Now you repeat the experiment in Q1 many times with matching observations
of many different particles. What is the probability that the observations emitted
by the matching equipment are observations of 'happy' particles?

# SIMULATOR

The simulator command (called particles) generates a simulated pool of particles
and then observes that pool with two different sampling processes. The pools are
comprised of 'calm' and 'happy' particles in the natural occuring ratio of 2:1.
The  observations generated by these processes are matched with a matching
process. The resulting stream of matched observations is confirmed to count the
number of particles and the number of 'happy' particles. This number is reported on
stderr.

The options that can be passed to the command are:

    Usage of particles:
      -match-type string
            Type of match: particle or quark (default "particle")
      -max-matches int
            Maximum number of matches (default 1000)
      -pool-size int
            Size of the particle pool. (default 10000)
      -process-A string
            Sampling process A: particle or quark (default "quark")
      -process-B string
            Sampling process B: particle or quark (default "particle")
      -verbose
            Be verbose about statistics

The expected probabilities and resulting matching process for each combination of parameters are:

    process-A   process-B   match-type  probability resulting-matching-process
    ----------------------------------------------- -----------------------------------------
    particle    particle    quark       1/5         U1u2 U3d dU4 U3d dU4 u1U2 U3d dU4 U3d dU4
    particle    particle    particle    1/3         U1u2 U3d dU4 u1U2 U3d dU4  # also, sample-particle
    particle    quark       quark       1/3         U1u2 U3d dU4 u1U2 U3d dU4
    quark       particle    particle    1/2         U1u2 U3d dU4 u1U2
    quark       quark       quark       1/2         U1u2 U3d dU4 u1U2 # also, sample-quark
    quark       quark       particle    2/3         U1u2 u1U2 U3d U1u2 u1U2 dU4

    U - observed up-quark
    u - unobserved up-quark
    d - unobserved down-quark
    u1U2, U1u2 - two different observations of a single happy particle
    dU4, U3d - two different observations of two different calm particles
    xy, ..., xy - two different observations of the same kind of the same particle

# NOTES

This simulation was inspired in part by a [variant](https://blog.jonseymour.net/the-boy-girl-paradox-with-a-twist) of the Boy-Girl paradox and ultimately by the
well-known [Boy or Girl paradox](https://en.wikipedia.org/wiki/Boy_or_Girl_paradox).
The analogy between those problems and this system is:

* mother/family <=> particle
* girl <=> up-quark
* boy  <=> down-quark
* boy-girl family <=> calm particle
* girl-girl family <=> happy particle

The amusing thing about the linked problem and the simulated universe is that you
one can have two sampling processes which accurately describe the world as they
have sampled it yet when the outputs of two different processes make conflicting
probabilistic statements about the state of a given object, it isn't
clear which (if either) is to be believed. It turns out that the means by which
observations from different processes are matched is extremely relevant to the
characteristation of the probability of observations emerging from the
matching process - in this example, the probability can range from 20% to 67%
depending on what combination of sampling and matching processes are used.