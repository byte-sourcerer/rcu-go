package rcu

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	state := NewRCU(
		time.Duration(30)*time.Millisecond,
		func(x *int) *int {
			var newX int
			if x == nil {
				newX = 0
			} else {
				newX = *x + 1
			}
			return &newX
		},
	)

	time.Sleep(time.Duration(100) * time.Millisecond)
	if *state.Load() < 2 {
		t.Errorf("state too small")
	}
	state.Close()
}
