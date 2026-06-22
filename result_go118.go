//go:build !go1.22
// +build !go1.22

package mo

// Match executes the first function if Result is valid and second function if invalid.
// It returns a new Result.
// Play: https://go.dev/play/p/-_eFaLJ31co
func (r Result[T]) Match(onSuccess func(value T) (T, error), onError func(err error) (T, error)) Result[T] {
	if r.isErr {
		return TupleToResult(onError(r.err))
	}
	return TupleToResult(onSuccess(r.value))
}

// Map executes the mapper function if Result is valid. It returns a new Result.
// Play: https://go.dev/play/p/-ndpN_b_OSc
func (r Result[T]) Map(mapper func(value T) (T, error)) Result[T] {
	if !r.isErr {
		return TupleToResult(mapper(r.value))
	}

	return Err[T](r.err)
}

// MapValue executes the mapper function if Result is valid. It returns a new Result.
func (r Result[T]) MapValue(mapper func(value T) T) Result[T] {
	if !r.isErr {
		return TupleToResult(mapper(r.value), nil)
	}

	return Err[T](r.err)
}

// FlatMap executes the mapper function if Result is valid. It returns a new Result.
// Play: https://go.dev/play/p/Ud5QjZOqg-7
func (r Result[T]) FlatMap(mapper func(value T) Result[T]) Result[T] {
	if !r.isErr {
		return mapper(r.value)
	}

	return Err[T](r.err)
}
