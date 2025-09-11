package executor

import "context"

type ctxKey struct{}

var executorOptsKey = ctxKey{}

func WithExecutorOpts(ctx context.Context, opts ...ExecutorOpts) context.Context {
	if ctxOpts, ok := ctx.Value(executorOptsKey).([]ExecutorOpts); ok {
		return context.WithValue(ctx, executorOptsKey, append(ctxOpts, opts...))
	}

	return context.WithValue(ctx, executorOptsKey, opts)
}

func GetExecutorOpts(ctx context.Context) []ExecutorOpts {
	if v := ctx.Value(executorOptsKey); v != nil {
		if opts, ok := v.([]ExecutorOpts); ok {
			return opts
		}
	}

	return nil
}
