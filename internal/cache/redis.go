package cache

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"gym/internal/constants"
	"log"
	"os"
	"time"
)

// Cache interface sets what a cache needs to do. In this case we are implementing this methods using Redis, but the implementation could change.
type Cache interface {
	Ping() (string, error)
	Exists(id string) (bool, error)
	Set(key, value string, expr time.Duration) error
	Get(key string) (string, error)
	Drop(key string) error
	Increment(key string) error
	SetMul(key ...string) error
	GetIfExists(key string, defaultValue string) (string, error)
	FlushAll() error
	GetAllKeys(pattern string) ([]string, error)
}

type universalClient struct {
	client redis.UniversalClient
	opts   redis.UniversalOptions
}

func NewClient(opts redis.UniversalOptions) Cache {
	client := &universalClient{
		client: redis.NewUniversalClient(&opts),
		opts:   opts,
	}

	// docker compose sync
	retryCount := constants.RetryDelayCycles
	for i := 1; i < retryCount+1; i++ {
		_, err := client.Ping()
		if err != nil {
			log.Printf("Error pinging cache: " + err.Error())
			if retryCount == i {
				errMsg := fmt.Sprintf("Not able to establish connection to cache %s after %d tries", opts.Addrs, i)
				panic(errMsg)
			}
			log.Printf(fmt.Sprintf("Could not connect to cache. Waiting %d seconds. %d retries left...", i, retryCount-i))
			// circuit breaker logic -> 1 2 3 4 5 6 7 8 ..
			time.Sleep(time.Duration(1*i+1) * time.Second)
			// try connection again
			client = &universalClient{
				client: redis.NewUniversalClient(&opts),
				opts:   opts,
			}
		} else {
			break
		}
	}
	log.Printf(constants.Blue + "CACHE CONNECTION CREATED at: " + os.Getenv("REDIS_ADDRESS") + constants.Reset)
	return client
}

func (r *universalClient) Exists(key string) (bool, error) {
	e, err := r.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return e == 1, nil
}

func (r *universalClient) Ping() (string, error) {
	return r.client.Ping().Result()
}

func (r *universalClient) Set(key, value string, expr time.Duration) error {
	return r.client.Set(key, value, expr).Err()
}

func (r *universalClient) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *universalClient) Drop(key string) error {
	return r.client.Del(key).Err()
}

func (r *universalClient) Increment(key string) error {
	return r.client.Incr(key).Err()
}

func (r *universalClient) SetMul(key ...string) error {
	for x := 0; x < len(key); x = x + 2 {
		key[x] = (key[x])
	}
	return r.client.MSet(key).Err()
}

func (r *universalClient) GetIfExists(key string, defaultValue string) (string, error) {
	exist, err := r.Exists(key)
	if err != nil {
		return "", err
	}
	if !exist {
		if err := r.Set(key, defaultValue, 0); err != nil {
			return "", err
		}
	}
	value, err := r.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *universalClient) FlushAll() error {
	return r.client.FlushAll().Err()
}

// GetAllKeys get all key that matches to pattern, this pattern possible values: *<name>* -> word contain, * -> all keys
func (r *universalClient) GetAllKeys(pattern string) ([]string, error) {
	return r.client.Keys(pattern).Result()
}

type RedisConfig struct {
	Password      string   `json:"password"`
	Master        string   `json:"master"`
	Addresses     []string `json:"addresses"`
	UseTLS        bool     `json:"useTls"`
	SkipSSLVerify bool     `json:"skipVerify"`
}
