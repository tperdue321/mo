package result

import (
	either "github.com/tperdue321/mo/either"
	option "github.com/tperdue321/mo/option"
)

// Ok builds a Result when value is valid.
// Play: https://go.dev/play/p/PDwADdzNoyZ
func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		isErr: false,
	}
}

// Err builds a Result when value is invalid.
// Play: https://go.dev/play/p/PDwADdzNoyZ
func Err[T any](err error) Result[T] {
	return Result[T]{
		err:   err,
		isErr: true,
	}
}

// TupleToResult convert a pair of T and error into a Result.
// Play: https://go.dev/play/p/KWjfqQDHQwa
func TupleToResult[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// Try returns either a Ok or Err object.
// Play: https://go.dev/play/p/ilOlQx-Mx42
func Try[T any](f func() (T, error)) Result[T] {
	return TupleToResult(f())
}

// Result respresent a result of an action having one
// of the following output: success or failure.
// An instance of Result is an instance of either Ok or Err.
// It could be compared to `Either[error, T]`.
type Result[T any] struct {
	isErr bool
	value T
	err   error
}

// IsOk returns true when value is valid.
// Play: https://go.dev/play/p/sfNvBQyZfgU
func (r Result[T]) IsOk() bool {
	return !r.isErr
}

// IsError returns true when value is invalid.
// Play: https://go.dev/play/p/xkV9d464scV
func (r Result[T]) IsError() bool {
	return r.isErr
}

// Error returns error when value is invalid or nil.
// Play: https://go.dev/play/p/CSkHGTyiXJ5
func (r Result[T]) Error() error {
	return r.err
}

// Get returns value even if empty
// Play: https://go.dev/play/p/8KyX3z6TuNo
func (r Result[T]) Get() T {
	return r.value
}

// MustGet returns value when Result is valid or panics.
// Play: https://go.dev/play/p/8LSlndHoTAE
func (r Result[T]) MustGet() T {
	if r.IsError() {
		panic(r.err)
	}

	return r.value
}

// OrElse returns value when Result is valid or default value.
// Play: https://go.dev/play/p/MN_ULx0soi6
func (r Result[T]) OrElse(fallback T) T {
	if r.IsOk() {
		return r.Get()
	}

	return fallback
}

// ToEither transforms a Result into an Either type.
// Play: https://go.dev/play/p/Uw1Zz6b952q
func (r Result[T]) ToEither() either.Either[error, T] {
	if r.IsOk() {
		return either.Right[error, T](r.value)
	}

	return either.Left[error, T](r.err)
}

// ToOption transforms a Result into None on failure and Some on OK
// Play: https://go.dev/play/p/Uw1Zz6b952q
func (r Result[T]) ToOption() option.Option[T] {
	if r.IsOk() {
		return option.Some[T](r.value)
	}

	return option.None[T]()
}

// ForEach executes the given side-effecting function if Result is valid.
func (r Result[T]) ForEach(mapper func(value T)) {
	if r.IsOk() {
		mapper(r.Get())
	}
}

// Fold is an alias for Match which executes the first function if Result is valid and second function if invalid.
// It returns a new Result.
// Play: https://go.dev/play/p/-_eFaLJ31co
func (r Result[T]) Fold(onSuccess func(value T) (T, error), onError func(err error) (T, error)) Result[T] {
	return r.Match(onSuccess, onError)
}

// Match executes the first function if Result is valid and second function if invalid.
// It returns a new Result.
// Play: https://go.dev/play/p/-_eFaLJ31co
func (r Result[T]) Match(onSuccess func(value T) (T, error), onError func(err error) (T, error)) Result[T] {
	if r.IsOk() {
		return TupleToResult(onSuccess(r.Get()))
	}

	return TupleToResult(onError(r.Error()))
}

// Map executes the mapper function if Result is valid. It returns a new Result.
// Play: https://go.dev/play/p/-ndpN_b_OSc
func (r Result[T]) Map(mapper func(value T) (T, error)) Result[T] {
	if r.IsOk() {
		return TupleToResult(mapper(r.Get()))
	}

	return Err[T](r.Error())
}

// MapErr executes the mapper function if Result is invalid. It returns a new Result.
// Play: https://go.dev/play/p/WraZixg9GGf
func (r Result[T]) MapErr(mapper func(error) (T, error)) Result[T] {
	if r.IsError() {
		return TupleToResult(mapper(r.Error()))
	}

	return Ok(r.Get())
}

// FlatMap executes the mapper function if Result is valid. It returns a new Result.
// Play: https://go.dev/play/p/Ud5QjZOqg-7
func (r Result[T]) FlatMap(mapper func(value T) Result[T]) Result[T] {
	if r.IsOk() {
		return mapper(r.value)
	}

	return Err[T](r.Error())
}
