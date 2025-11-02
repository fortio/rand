package rand_test

import (
	"fmt"
	"math"
	"testing"

	"fortio.org/rand"
)

func ExampleRand_Vec3() {
	// Example of how it can be used with a vector type (like in fortio.org/ray).
	type Vec3 struct {
		x, y, z float64
	}

	NewVec3 := func(x, y, z float64) Vec3 {
		return Vec3{x, y, z}
	}
	RandomVec3 := func(r rand.Rand) Vec3 {
		return NewVec3(r.Vec3())
	}

	r := rand.New(42)
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
	return rand.New(0)
}

// TestRandom just... exercises the Random function
// and that values are ... different.
func TestRandom(t *testing.T) {
	const samples = 10
	results := map[vec3]struct{}{}
	expected := interval{Start: 0.0, End: 1.0}
	r := RandForTests()
	for range samples {
		v := newVec3(r.Vec3())
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

// TestRandomUnitVectorCorrectness verifies that RandomUnitVector
// produce vectors of unit length.
func TestRandomUnitVectorCorrectness(t *testing.T) {
	r := RandForTests()
	const samples = 100
	const tolerance = 1e-9

	for i := range samples {
		v := newVec3(r.UnitVector())
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
		v := newVec3(r.UnitVector())
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
		x, y := r.InDisc(radius)

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
		x, y := r.InDiscAngle(radius)

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
		x1, y1 := r.InDisc(radius)
		x2, y2 := r.InDiscAngle(radius)

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

func TestNormFloat64(t *testing.T) {
	r := RandForTests()
	const samples = 1000
	var sum, sumSquares float64

	for range samples {
		v := r.NormFloat64()
		sum += v
		sumSquares += v * v
	}

	// Normal distribution should have mean ≈ 0 and stddev ≈ 1
	mean := sum / samples
	variance := sumSquares/samples - mean*mean
	stddev := math.Sqrt(variance)

	// With 1000 samples, allow reasonable tolerance
	const meanTolerance = 0.1
	const stddevTolerance = 0.2

	if math.Abs(mean) > meanTolerance {
		t.Errorf("NormFloat64() mean = %.6f, want ≈0 (within %.6f)", mean, meanTolerance)
	}
	if math.Abs(stddev-1.0) > stddevTolerance {
		t.Errorf("NormFloat64() stddev = %.6f, want ≈1 (within %.6f)", stddev, stddevTolerance)
	}
}

func TestIntN(t *testing.T) {
	r := RandForTests()
	const samples = 1000

	// Test with n=10
	n := 10
	counts := make([]int, n)
	for range samples {
		v := r.IntN(n)
		if v < 0 || v >= n {
			t.Errorf("IntN(%d) = %v, want in [0,%d)", n, v, n)
		}
		counts[v]++
	}

	// Check that all values were hit at least once
	for i, count := range counts {
		if count == 0 {
			t.Errorf("IntN(%d) never produced value %d in %d samples", n, i, samples)
		}
	}
}

func TestUint64(t *testing.T) {
	r := RandForTests()
	const samples = 100
	results := make(map[uint64]struct{})

	for range samples {
		v := r.Uint64()
		results[v] = struct{}{}
	}

	// With 100 samples of uint64, we should get all unique values
	// (collision probability is negligible)
	if len(results) != samples {
		t.Errorf("Uint64() produced %d unique values, want %d", len(results), samples)
	}
}

func TestRandomInRange(t *testing.T) {
	r := RandForTests()
	const samples = 1000

	tests := []struct {
		start, end float64
	}{
		{0, 1},
		{-1, 1},
		{10, 20},
		{-5.5, -2.3},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("range_%.1f_to_%.1f", tt.start, tt.end), func(t *testing.T) {
			for range samples {
				v := r.Float64Range(tt.start, tt.end)
				if v < tt.start || v >= tt.end {
					t.Errorf("RandomInRange(%v, %v) = %v, want in [%v,%v)",
						tt.start, tt.end, v, tt.start, tt.end)
				}
			}
		})
	}
}

func BenchmarkRandomUnitVectorNorm(b *testing.B) {
	r := RandForTests()
	for range b.N {
		_, _, _ = r.UnitVector()
	}
}
