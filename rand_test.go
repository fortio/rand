package rand_test

import (
	"fmt"
	"math"
	"testing"

	"fortio.org/rand"
)

func ExampleRand_Random3() {
	// Example of how it can be used with a vector type (like in fortio.org/ray).
	type Vec3 struct {
		x, y, z float64
	}

	NewVec3 := func(x, y, z float64) Vec3 {
		return Vec3{x, y, z}
	}
	RandomVec3 := func(r rand.Rand) Vec3 {
		return NewVec3(r.Random3())
	}

	r := rand.NewRand(42)
	v := RandomVec3(r)

	fmt.Printf("%#v", v)
	// Output:
	// rand_test.Vec3{x:0.7680643711325947, y:0.5560374848919416, z:0.6664016849143646}
}

type interval struct {
	Start, End float64
}

func (iv interval) Contains(v float64) bool {
	return v >= iv.Start && v < iv.End
}

type vec3 struct {
	x, y, z float64
}

func newVec3(x, y, z float64) vec3 {
	return vec3{x: x, y: y, z: z}
}

func (v vec3) Components() [3]float64 {
	return [3]float64{v.x, v.y, v.z}
}

func Length(v vec3) float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func RandForTests() rand.Rand {
	return rand.NewRand(0)
}

// TestRandom just... exercises the Random function
// and that values are ... different.
func TestRandom(t *testing.T) {
	const samples = 10
	results := map[vec3]struct{}{}
	expected := interval{Start: 0.0, End: 1.0}
	r := RandForTests()
	for range samples {
		v := newVec3(r.Random3())
		// Check each component is in [0,1)
		c := v.Components()
		for i := range 3 {
			if !expected.Contains(c[i]) {
				t.Errorf("Random() component %d = %v, want in [0,1)", i, c[i])
			}
		}
		// Collect unique samples
		results[v] = struct{}{}
	}
	if len(results) != samples {
		t.Errorf("Random() produced %d unique samples, want %d", len(results), samples)
	}
}

// TestRandomUnitVectorCorrectness verifies that all three RandomUnitVector variants
// produce vectors of unit length.
func TestRandomUnitVectorCorrectness(t *testing.T) {
	r := RandForTests()
	const samples = 100
	const tolerance = 1e-9

	for i := range samples {
		v := newVec3(r.RandomUnitVector())
		length := Length(v)

		if math.Abs(length-1.0) > tolerance {
			t.Errorf("sample %d: Length() = %.15f, want 1.0 (diff: %.15e)",
				i, length, length-1.0)
		}
	}
}

// TestRandomUnitVectorDistribution checks that the generated vectors are
// uniformly distributed over the unit sphere by testing:
// 1. Mean of components approaches zero.
// 2. Standard deviation of each component approaches expected value.
// 3. Points cover all octants of the sphere.
func TestRandomUnitVectorDistribution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping distribution test in short mode")
	}
	const samples = 100000
	r := RandForTests()

	// Track statistics
	var sumX, sumY, sumZ float64
	var sumX2, sumY2, sumZ2 float64
	octantCounts := make([]int, 8)

	for range samples {
		v := newVec3(r.RandomUnitVector())
		components := v.Components()
		x, y, z := components[0], components[1], components[2]

		// Accumulate for mean and variance
		sumX += x
		sumY += y
		sumZ += z
		sumX2 += x * x
		sumY2 += y * y
		sumZ2 += z * z

		// Count octant
		octant := 0
		if x > 0 {
			octant |= 1
		}
		if y > 0 {
			octant |= 2
		}
		if z > 0 {
			octant |= 4
		}
		octantCounts[octant]++
	}

	// Check means are near zero
	meanX := sumX / samples
	meanY := sumY / samples
	meanZ := sumZ / samples

	// For uniform distribution on sphere, mean should be (0,0,0)
	// With 100k samples, standard error ≈ 1/sqrt(100000) ≈ 0.003
	// We use 5 sigma threshold for robustness: 5 * 0.003 ≈ 0.015
	const meanTolerance = 0.015
	if math.Abs(meanX) > meanTolerance {
		t.Errorf("mean X = %.6f, want ≈0 (within %.6f)", meanX, meanTolerance)
	}
	if math.Abs(meanY) > meanTolerance {
		t.Errorf("mean Y = %.6f, want ≈0 (within %.6f)", meanY, meanTolerance)
	}
	if math.Abs(meanZ) > meanTolerance {
		t.Errorf("mean Z = %.6f, want ≈0 (within %.6f)", meanZ, meanTolerance)
	}

	// Check variance for each component
	// For uniform distribution on unit sphere, variance of each component ≈ 1/3
	varX := sumX2/samples - meanX*meanX
	varY := sumY2/samples - meanY*meanY
	varZ := sumZ2/samples - meanZ*meanZ

	expectedVar := 1.0 / 3.0
	const varTolerance = 0.01 // Allow 1% deviation

	if math.Abs(varX-expectedVar) > varTolerance {
		t.Errorf("variance X = %.6f, want ≈%.6f (within %.6f)",
			varX, expectedVar, varTolerance)
	}
	if math.Abs(varY-expectedVar) > varTolerance {
		t.Errorf("variance Y = %.6f, want ≈%.6f (within %.6f)",
			varY, expectedVar, varTolerance)
	}
	if math.Abs(varZ-expectedVar) > varTolerance {
		t.Errorf("variance Z = %.6f, want ≈%.6f (within %.6f)",
			varZ, expectedVar, varTolerance)
	}

	// Check octant distribution
	// Each octant should contain approximately samples/8 points
	expectedPerOctant := samples / 8
	// Allow 15% deviation from expected
	octantTolerance := float64(expectedPerOctant) * 0.15

	for octant, count := range octantCounts {
		diff := math.Abs(float64(count) - float64(expectedPerOctant))
		if diff > octantTolerance {
			t.Errorf("octant %d: count = %d, want ≈%d (within %.0f)",
				octant, count, expectedPerOctant, octantTolerance)
		}
	}
}

func TestRandFloat64(t *testing.T) {
	r := RandForTests()
	const samples = 1000
	for range samples {
		v := r.Float64()
		if v < 0 || v >= 1 {
			t.Errorf("Float64() = %v, want in [0,1)", v)
		}
	}
}

func TestSampleDisc(t *testing.T) {
	r := RandForTests()
	radius := 2.5
	const samples = 1000

	for range samples {
		x, y := r.SampleDisc(radius)

		// Check that point is within disc
		dist := math.Sqrt(x*x + y*y)
		if dist > radius {
			t.Errorf("SampleDisc(%v) = (%v, %v), distance %v exceeds radius",
				radius, x, y, dist)
		}
	}
}

func TestSampleDiscAngle(t *testing.T) {
	r := RandForTests()
	radius := 3.0
	const samples = 1000

	for range samples {
		x, y := r.SampleDiscAngle(radius)

		// Check that point is within disc
		dist := math.Sqrt(x*x + y*y)
		if dist > radius {
			t.Errorf("SampleDiscAngle(%v) = (%v, %v), distance %v exceeds radius",
				radius, x, y, dist)
		}
	}
}

func TestSampleDiscMethods(t *testing.T) {
	// Compare both methods produce valid results
	r := RandForTests()
	radius := 1.0
	const samples = 100

	for range samples {
		x1, y1 := r.SampleDisc(radius)
		x2, y2 := r.SampleDiscAngle(radius)

		// Both should be within disc
		dist1 := math.Sqrt(x1*x1 + y1*y1)
		dist2 := math.Sqrt(x2*x2 + y2*y2)

		if dist1 > radius {
			t.Errorf("SampleDisc produced point outside disc: (%v, %v), dist=%v",
				x1, y1, dist1)
		}
		if dist2 > radius {
			t.Errorf("SampleDiscAngle produced point outside disc: (%v, %v), dist=%v",
				x2, y2, dist2)
		}
	}
}

func BenchmarkRandomUnitVectorNorm(b *testing.B) {
	r := RandForTests()
	for range b.N {
		_, _, _ = r.RandomUnitVector()
	}
}
