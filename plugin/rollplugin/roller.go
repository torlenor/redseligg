package rollplugin

import "math/rand"

type roller struct {
	randomizer *rand.Rand
}

func newRoller(seed int64) roller {
	return roller{
		randomizer: rand.New(rand.NewSource(99)),
	}
}

func (r roller) random(max int) int {
	return r.randomizer.Intn(max + 1)
}
