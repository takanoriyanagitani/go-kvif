package kvif

import (
	"context"
)

func Curry[T, U, V any](f func(T, U) V) func(T) func(U) V {
	return func(t T) func(U) V {
		return func(u U) V {
			return f(t, u)
		}
	}
}

func CurryCtx[T, U any](f func(context.Context, T) (U, error)) func(context.Context) func(T) (U, error) {
	return func(ctx context.Context) func(T) (U, error) {
		return func(t T) (U, error) {
			return f(ctx, t)
		}
	}
}
