package httpx

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

var (
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
)

type options struct {
	SkipVerify bool
	Pooled     bool
	Log        retryablehttp.LeveledLogger

	RetryWaitMin *time.Duration
	RetryWaitMax *time.Duration
	RetryMax     *int

	RetryPolicy retryablehttp.CheckRetry
	Backoff     retryablehttp.Backoff
}

type Option func(options *options)

func WithSkipVerify(skip bool) Option {
	return func(options *options) {
		options.SkipVerify = skip
	}
}

func WithPooled(pooled bool) Option {
	return func(options *options) {
		options.Pooled = pooled
	}
}

func WithLog(log retryablehttp.LeveledLogger) Option {
	return func(options *options) {
		options.Log = log
	}
}

func WithRetryWaitMin(retryWaitMin time.Duration) Option {
	return func(options *options) {
		options.RetryWaitMin = &retryWaitMin
	}
}

func WithRetryWaitMax(retryWaitMax time.Duration) Option {
	return func(options *options) {
		options.RetryWaitMax = &retryWaitMax
	}
}

func WithRetryMax(retryMax int) Option {
	return func(options *options) {
		options.RetryMax = &retryMax
	}
}

func WithRetryPolicy(retryPolicy retryablehttp.CheckRetry) Option {
	return func(options *options) {
		options.RetryPolicy = retryPolicy
	}
}

func WithBackoff(backoff retryablehttp.Backoff) Option {
	return func(options *options) {
		options.Backoff = backoff
	}
}

func setOptDefault(options *options) {
	if options.RetryWaitMin == nil {
		options.RetryWaitMin = &defaultRetryWaitMin
	}
	if options.RetryWaitMax == nil {
		options.RetryWaitMax = &defaultRetryWaitMax
	}
	if options.RetryMax == nil {
		options.RetryMax = &defaultRetryMax
	}

	if options.RetryPolicy == nil {
		options.RetryPolicy = RetryPolicy
	}
	if options.Backoff == nil {
		options.Backoff = retryablehttp.DefaultBackoff
	}
}
