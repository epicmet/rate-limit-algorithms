package tokenbucket

import (
	"fmt"
	"net/http"
	"time"

	statemanager "github.com/epicmet/rate-limit-algorithms/server/rate-limiter/state-manager"
	"github.com/gin-gonic/gin"
)

type TokenBucket struct {
	BucketSize int
	RefillRate time.Duration
	Key        string
	sm         statemanager.StateManager
}

func New(keyPrefix string, capacity int, refillRate time.Duration, sm statemanager.StateManager) *TokenBucket {
	key := fmt.Sprintf("rate-limiter::token-bucket::bucket::%s", keyPrefix)

	ticker := time.NewTicker(refillRate)

	sm.Set(key, capacity, 0)

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				{
					sm.Set(key, capacity, 0)
				}
			}
		}
	}()

	return &TokenBucket{
		BucketSize: capacity,
		RefillRate: refillRate,
		Key:        key,
		sm:         sm,
	}
}

func (tb TokenBucket) Allow() bool {
	bucketCounter, err := tb.sm.GetIntValue(tb.Key)
	if err != nil || bucketCounter <= 0 {
		return false
	}
	if _, err := tb.sm.Decr(tb.Key); err != nil {
		fmt.Println("Couldn't Decr the bucket counter. Error: ", err.Error())
		return false
	}

	return true
}

func (tb TokenBucket) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !tb.Allow() {
			c.JSON(
				http.StatusTooManyRequests,
				gin.H{},
			)
			c.Abort()
		}
	}
}
