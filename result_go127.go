//go:build go1.27
// +build go1.27

package mo

// Match executes the first function if Result is valid and second function if invalid.
// It returns a new Result.
// Play: https://go.dev/play/p/-_eFaLJ31co
func (r Result[T]) Match[U any](onSuccess func(value T) (U, error), onError func(err error) (U, error)) Result[U] {
	if r.isErr {
		return TupleToResult(onError(r.err))
	}
	return TupleToResult(onSuccess(r.value))
}

// Map executes the mapper function if Result is valid. It returns a new Result.
// Play: https://go.dev/play/p/-ndpN_b_OSc
func (r Result[T]) Map[U any](mapper func(value T) (U, error)) Result[U] {
	if !r.isErr {
		return TupleToResult(mapper(r.value))
	}

	return Err[U](r.err)
}

// MapValue executes the mapper function if Result is valid. It returns a new Result.
func (r Result[T]) MapValue[U any](mapper func(value T) U) Result[U] {
	if !r.isErr {
		return TupleToResult(mapper(r.value), nil)
	}

	return Err[U](r.err)
}

// FlatMap executes the mapper function if Result is valid. It returns a new Result.
// Play: https://go.dev/play/p/Ud5QjZOqg-7
func (r Result[T]) FlatMap[U any](mapper func(value T) Result[U]) Result[U] {
	if !r.isErr {
		return mapper(r.value)
	}

	return Err[U](r.err)
}
