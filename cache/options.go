package cache

import "time"

type Option func(o *Options)

type Options struct {
	Expiration      time.Duration
	CleanupInterval time.Duration

	RedisDB          int
	RedisHost        string
	RedisUser        string
	RedisPassword    string
	RedisMaxidle     int
	RedisMaxactive   int
	RedisIdleTimeout time.Duration

	Endpoints []string
}

func ApplyOptions(opts ...Option) *Options {
	o := &Options{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func ApplyOptionsWithDefault(defaultOptions *Options, opts ...Option) *Options {
	returnedOptions := &Options{}
	*returnedOptions = *defaultOptions

	for _, opt := range opts {
		opt(returnedOptions)
	}

	return returnedOptions
}

// WithExpiration allows to specify an expiration time when setting a value.
func WithExpiration(expiration time.Duration) Option {
	return func(o *Options) {
		o.Expiration = expiration
	}
}

// WithCleanupInterval allows to specify a cleanup interval.
func WithCleanupInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.CleanupInterval = interval
	}
}

// WithRedisDB allows to specify a redis db.
func WithRedisDB(db int) Option {
	return func(o *Options) {
		o.RedisDB = db
	}
}

// WithRedisHost allows to specify a redis host.
func WithRedisHost(host string) Option {
	return func(o *Options) {
		o.RedisHost = host
	}
}

// WithRedisUser allows to specify a redis user.
func WithRedisUser(user string) Option {
	return func(o *Options) {
		o.RedisUser = user
	}
}

// WithRedisPassword allows to specify a redis password.
func WithRedisPassword(password string) Option {
	return func(o *Options) {
		o.RedisPassword = password
	}
}

// WithEndpoints allows to specify a list of endpoints.
func WithEndpoints(endpoints []string) Option {
	return func(o *Options) {
		o.Endpoints = endpoints
	}
}

// WithRedisMaxidle allows to specify a redis max idle.
func WithRedisMaxidle(maxidle int) Option {
	return func(o *Options) {
		o.RedisMaxidle = maxidle
	}
}

// WithRedisMaxactive allows to specify a redis max active.
func WithRedisMaxactive(maxactive int) Option {
	return func(o *Options) {
		o.RedisMaxactive = maxactive
	}
}

// WithRedisIdleTimeout allows to specify a redis idle timeout.
func WithRedisIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *Options) {
		o.RedisIdleTimeout = idleTimeout
	}
}
