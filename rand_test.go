package rand_test

import (
	"fmt"

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
