package leakingbucket

import (
	"fmt"
	"net/http"
	"time"

	statemanager "github.com/epicmet/rate-limit-algorithms/server/rate-limiter/state-manager"
	"github.com/gin-gonic/gin"
)

type LeakingBucket struct {
	BucketCapcity int64
	OutflowRate   time.Duration
	Key           string
	lastLeak      time.Time
	sm            statemanager.StateManager
}

func New(keyPrefix string, capacity int64, outflowRate time.Duration, sm statemanager.StateManager) *LeakingBucket {
	key := fmt.Sprintf("%s::rate-limiter::leaking-bucket::bucket", keyPrefix)
	sm.Set(key, 0, 0)

	return &LeakingBucket{
		BucketCapcity: capacity,
		OutflowRate:   outflowRate,
		Key:           key,
		sm:            sm,
	}
}

func (lb *LeakingBucket) Allow() bool {
	currentBucketCount, err := lb.sm.GetIntValue(lb.Key)

	// Drop the request of the bucket is at full capacity
	if err != nil || currentBucketCount >= lb.BucketCapcity {
		return false
	}

	lb.sm.Incr(lb.Key)

	now := time.Now()

	if lb.lastLeak.IsZero() {
		lb.lastLeak = now
		return true
	}

	sleepFor := lb.OutflowRate - now.Sub(lb.lastLeak)
	if sleepFor < 0 {
		sleepFor = lb.OutflowRate
	}
	time.Sleep(sleepFor)

	lb.sm.Decr(lb.Key)
	lb.lastLeak = now.Add(sleepFor)
	return true
}

func (lb *LeakingBucket) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !lb.Allow() {
			c.JSON(
				http.StatusTooManyRequests,
				gin.H{},
			)
			c.Abort()
		}
	}
}
