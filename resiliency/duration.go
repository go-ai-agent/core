package resiliency

import (
	"math"
	"time"
)

type ExponentialDuration struct {
	rate    float64
	current float64
	count   int
}

func NewExponentialDuration(rate float64, initialMS float64) *ExponentialDuration {
	e := new(ExponentialDuration)
	e.rate = rate
	e.current = initialMS
	return e
}

func (e *ExponentialDuration) Current() time.Duration {
	return time.Duration(math.Trunc(e.current)) * time.Millisecond
}

func (e *ExponentialDuration) Eval() time.Duration {
	if e.count != 0 {
		e.current = e.current * (1 - e.rate)
	}
	e.count++
	return e.Current() //time.Duration(math.Trunc(e.current)) * time.Millisecond
}