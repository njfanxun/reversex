package lib

import (
	r "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

var randOnce sync.Once

func SeedMathRand() {
	randOnce.Do(func() {
		n, err := r.Int(r.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			rand.Seed(time.Now().UTC().UnixNano())
		} else {
			rand.Seed(n.Int64())
		}
	})
}
