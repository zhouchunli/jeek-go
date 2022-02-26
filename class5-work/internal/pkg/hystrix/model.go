package hystrix

import (
	"sync"
	"time"
)

type Bucket struct {
	sync.RWMutex
	Total int
	Failed    int
	Timestamp time.Time
}

func NewBucket() *Bucket {
	return &Bucket{
		Timestamp: time.Now(),
	}
}

func (b *Bucket) Record(result bool) {
	b.Lock()
	defer b.Unlock()
	if !result {
		b.Failed++
	}
	b.Total++
}

