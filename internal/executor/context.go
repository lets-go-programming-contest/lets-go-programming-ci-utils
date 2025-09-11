package executor

import "context"

type ctxKey struct{}

var executorOptsKey = ctxKey{}

func WithExecutorOpts(ctx context.Context, opts ...Opts) context.Context {
	if ctxOpts, ok := ctx.Value(executorOptsKey).([]Opts); ok {
		return context.WithValue(ctx, executorOptsKey, append(ctxOpts, opts...))
	}

	return context.WithValue(ctx, executorOptsKey, opts)
}

func GetExecutorOpts(ctx context.Context) []Opts {
	if v := ctx.Value(executorOptsKey); v != nil {
		if opts, ok := v.([]Opts); ok {
			return opts
		}
	}

	return nil
}
