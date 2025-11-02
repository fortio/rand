// Package rand is a wrapper over stdlib math/rand/v2 PCG,
// the fastest available in stdlib meant to have 1 instance per goroutine
// for maximum performance.
// It provides convenience methods to either have different rng each time (seed 0)
// or a specific repeatable sequence.
// It also provides methods useful for vector graphics, ray tracing and
// simulations (e.g. random unit 3d vectors).
package rand

import (
	"math"
	"math/rand/v2"
)

// Rand wraps a random number generator, is meant to be embedded in other structs and
// reused but not shared across goroutines. It's ok to copy by value as the underlying
// rng state is a pointer.
type Rand struct {
	rng *rand.Rand
}

// NewRand generates a new Rand with the given seed. If seed is 0, a random seed is used.
// It's meant for the single goroutine or pre-multiple goroutine (main thread) case.
func NewRand(seed uint64) Rand {
	return NewRandIdx(0, seed)
}

// NewRandIdx creates a new Rand using the given index and seed.
// idx can be used to create different Rand instances for different goroutines.
// If seed is 0, a random seed is used and index is ignored.
//
//nolint:gosec // not crypto use.
func NewRandIdx(idx int, seed uint64) Rand {
	seed1 := uint64(idx)
	seed2 := seed
	if seed == 0 {
		seed1 = rand.Uint64()
		seed2 = rand.Uint64()
	}
	return newRandSeeds(seed1, seed2)
}

//nolint:gosec // not crypto use.
func newRandSeeds(seed1, seed2 uint64) Rand {
	return Rand{rng: rand.New(rand.NewPCG(seed1, seed2))}
}

// Forward methods to underlying rng

func (r Rand) Float64() float64 {
	return r.rng.Float64()
}

func (r Rand) NormFloat64() float64 {
	return r.rng.NormFloat64()
}

func (r Rand) IntN(n int) int {
	return r.rng.IntN(n)
}

func (r Rand) Uint64() uint64 {
	return r.rng.Uint64()
}

// Random3 generates a random vector of 3 components in [0,1).
func (r Rand) Random3() (float64, float64, float64) {
	return r.rng.Float64(), r.rng.Float64(), r.rng.Float64()
}

// RandomInRange generates a random value in the range [start,end).
func (r Rand) RandomInRange(start, end float64) float64 {
	l := end - start
	return start + l*r.rng.Float64()
}

// RandomUnitVector generates a random unit vector using normal distribution.
// It is the fastest of the three methods tested (versus rejection or angle methods)
// and produces uniformly distributed points on the unit sphere.
// Being both correct and most efficient, this is now the only method for generating
// random unit vectors provided (compared to the original tray 3 methods).
func (r Rand) RandomUnitVector() (float64, float64, float64) {
	for {
		x, y, z := r.rng.NormFloat64(), r.rng.NormFloat64(), r.rng.NormFloat64()
		radius := math.Sqrt(x*x + y*y + z*z)
		if radius > 1e-24 {
			return x / radius, y / radius, z / radius
		}
	}
}

// SampleDisc returns a random point (x,y) within a disc of the given radius
// from random source (and currently implemented via rejection sampling).
func (r Rand) SampleDisc(radius float64) (x, y float64) {
	for {
		x = 2*r.rng.Float64() - 1.0
		y = 2*r.rng.Float64() - 1.0
		if x*x+y*y <= 1 {
			break
		}
	}
	return radius * x, radius * y
}

// SampleDiscAngle returns a random point (x,y) within a disc of given radius.
// Angle method.
func (r Rand) SampleDiscAngle(radius float64) (x, y float64) {
	theta := 2.0 * math.Pi * r.rng.Float64()
	rad := radius * math.Sqrt(r.rng.Float64())
	x = rad * math.Cos(theta)
	y = rad * math.Sin(theta)
	return x, y
}
