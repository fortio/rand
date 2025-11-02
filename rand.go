package rand

import (
	"math"
	"math/rand/v2"
)

// Rand wraps a random number generator, is meant to be embedded in other structs and
// reused during rendering but not shared across goroutines.
type Rand struct {
	rng *rand.Rand
}

// NewRand generates a new (scene) Rand with the given seed. If seed is 0, a random seed is used.
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

func (r Rand) Float64() float64 {
	return r.rng.Float64()
}

// Random3 generates a random vector of 3 components in [0,1).
func Random3(r Rand) (float64, float64, float64) {
	return r.rng.Float64(), r.rng.Float64(), r.rng.Float64()
}

// RandomInRange generates a random vector with each component in the Interval
// excluding the end.
func RandomInRange(r Rand, start, end float64) float64 {
	l := end - start
	return start + l*r.rng.Float64()
}

// RandomUnitVector generates a random unit vector using normal distribution.
// It is the fastest of the three methods provided here and produces uniformly
// distributed points on the unit sphere. Being both correct and most efficient,
// this is the preferred method for generating random unit vectors and thus gets
// the default name.
func RandomUnitVector(r Rand) (float64, float64, float64) {
	for {
		x, y, z := r.rng.NormFloat64(), r.rng.NormFloat64(), r.rng.NormFloat64()
		radius := math.Sqrt(x*x + y*y + z*z)
		if radius > 1e-24 {
			return x / radius, y / radius, z / radius
		}
	}
}

// SampleDisc returns a random point (x,y) within a disc of radius r
// using the provided random source (and currently implemented via rejection sampling).
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

// SampleDiscAngle returns a random point (x,y) within a disc of radius r.
// Angle method.
func (r Rand) SampleDiscAngle(radius float64) (x, y float64) {
	theta := 2.0 * math.Pi * r.rng.Float64()
	rad := radius * math.Sqrt(r.rng.Float64())
	x = rad * math.Cos(theta)
	y = rad * math.Sin(theta)
	return x, y
}
