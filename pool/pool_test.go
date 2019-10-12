package pool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testRunnable struct {
	wasStarted bool
	wasStopped bool
}

func (t *testRunnable) Start() {
	t.wasStarted = true
}

func (t *testRunnable) Stop() {
	t.wasStopped = true
}

func TestPool_Add(t *testing.T) {
	assert := assert.New(t)

	p := Pool{}

	r1 := testRunnable{}

	assert.Equal(0, p.Len())

	p.Add(&r1)
	assert.Equal(1, p.Len())
}

func TestPool_StartStop(t *testing.T) {
	assert := assert.New(t)

	p := Pool{}

	r1 := testRunnable{}

	assert.Equal(false, r1.wasStarted)
	assert.Equal(false, r1.wasStopped)

	p.Add(&r1)
	assert.Equal(false, r1.wasStarted)
	assert.Equal(false, r1.wasStopped)

	p.StartAll()
	time.Sleep(100 * time.Millisecond) // how to do this in a nicer way without sleep?
	assert.Equal(true, r1.wasStarted)
	assert.Equal(false, r1.wasStopped)

	p.StopAll()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(true, r1.wasStarted)
	assert.Equal(true, r1.wasStopped)
}
