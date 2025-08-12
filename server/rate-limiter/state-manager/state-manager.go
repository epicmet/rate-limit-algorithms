package statemanager

import (
	"log"
	"time"
)

type Config struct {
	Addr string
}

type StateManager interface {
	Set(key string, value interface{}, expireTime time.Duration) (string, error)
	GetIntValue(key string) (int64, error)
	Decr(key string) (int64, error)
}

func New(t string, c Config) StateManager {
	switch t {
	case "redis":
		{
			return newRedis(c.Addr)
		}
	default:
		{
			log.Fatalf("Cound not create unknown state manager %s", t)
			return nil
		}
	}
}
