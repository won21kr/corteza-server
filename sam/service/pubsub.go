package service

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/config"
)

type PubSub struct {
	ctx context.Context

	redisServerAddr string

	healthCheckInterval time.Duration
}

func (ps PubSub) New(config *config.PubSub, ctx context.Context) (*PubSub, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &PubSub{
		ctx:                 ctx,
		redisServerAddr:     config.RedisAddr,
		healthCheckInterval: time.Minute,
	}, nil
}

func (ps *PubSub) With(ctx context.Context) *PubSub {
	return &PubSub{
		ctx:                 ctx,
		healthCheckInterval: ps.healthCheckInterval,
	}
}

func (ps *PubSub) dial() (redis.Conn, error) {
	readTimeout := redis.DialReadTimeout(ps.healthCheckInterval + 10*time.Second)
	writeTimeout := redis.DialWriteTimeout(10 * time.Second)
	return redis.Dial("tcp", ps.redisServerAddr, readTimeout, writeTimeout)
}

func (ps *PubSub) Subscribe(onStart func() error, onMessage func(channel string, payload []byte) error, channels ...string) error {
	if len(channels) == 0 {
		return errors.New("Need to subscribe at least to one channel")
	}

	// main redis connection
	conn, err := ps.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	// pubsub object
	psc := redis.PubSubConn{Conn: conn}
	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				switch n.Count {
				case len(channels):
					// Notify application when all channels are subscribed.
					if err := onStart(); err != nil {
						done <- err
						return
					}
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(ps.healthCheckInterval)
	defer ticker.Stop()

	cleanup := func(err error) error {
		psc.Unsubscribe()
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := psc.Ping(""); err != nil {
				return cleanup(err)
			}
		case <-ps.ctx.Done():
			return cleanup(ps.ctx.Err())
		case err := <-done:
			return err
		}
	}
}

func (ps *PubSub) Publish(channel, payload string) error {
	// main redis connection
	conn, err := ps.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	// publish payload on channel
	_, err = conn.Do("PUBLISH", channel, payload)
	return err
}
