package result

import (
	either "github.com/tperdue321/mo/either"
	option "github.com/tperdue321/mo/option"
)

type ResultWrapper[A, B any] struct {
	r Result[A]
}


// Ok builds a Result when value is valid.
// Play: https://go.dev/play/p/PDwADdzNoyZ
func OkWrapper[A, B any](value A) ResultWrapper[A, B] {
	return ResultWrapper[A, B]{
		r: Result[A]{
			value: value,
			isErr: false,
		},
	}
}

// Err builds a Result when value is invalid.
// Play: https://go.dev/play/p/PDwADdzNoyZ
func ErrWrapper[A, B any](err error) ResultWrapper[A, B] {
	return ResultWrapper[A, B]{
		r: Result[A]{
			err:   err,
			isErr: true,
		},
	}
}

// TupleToResult convert a pair of T and error into a Result.
// Play: https://go.dev/play/p/KWjfqQDHQwa
func TupleToResultWrapper[A, B any](value A, err error) ResultWrapper[A, B] {
	if err != nil {
		return ErrWrapper[A, B](err)
	}
	return OkWrapper[A, B](value)
}

// Try returns either a Ok or Err object.
// Play: https://go.dev/play/p/ilOlQx-Mx42
func TryWrapper[A, B any](f func() (A, error)) ResultWrapper[A, B] {
	return TupleToResultWrapper[A, B](f())
}


// IsOk returns true when value is valid.
// Play: https://go.dev/play/p/sfNvBQyZfgU
func (rw ResultWrapper[A, B]) IsOk() bool {
	return rw.r.IsOk()
}

// IsError returns true when value is invalid.
// Play: https://go.dev/play/p/xkV9d464scV
func (rw ResultWrapper[A, B]) IsError() bool {
	return rw.r.IsError()
}

// Error returns error when value is invalid or nil.
// Play: https://go.dev/play/p/CSkHGTyiXJ5
func (rw ResultWrapper[A, B]) Error() error {
	return rw.Error()
}

// Get returns value even if empty
// Play: https://go.dev/play/p/8KyX3z6TuNo
func (rw ResultWrapper[A, B]) Get() A {
	return rw.Get()
}

// MustGet returns value when Result is valid or panics.
// Play: https://go.dev/play/p/8LSlndHoTAE
func (rw ResultWrapper[A, B]) MustGet() A {
	if rw.IsError() {
		panic(rw.r.Error())
	}

	return rw.Get()
}

// OrElse returns value when Result is valid or default value.
// Play: https://go.dev/play/p/MN_ULx0soi6
func (rw ResultWrapper[A, B]) OrElse(fallback A) A {
	if rw.IsError() {
		return fallback
	}

	return rw.Get()
}

// ToEither transforms a Result into an Either type.
// Play: https://go.dev/play/p/Uw1Zz6b952q
func (rw ResultWrapper[A, B]) ToEitherWrapper() either.EitherWrapper[error, A, B] {
	if rw.IsError() {
		return either.LeftEitherWrapper[error, A, B](rw.r.err)
	}

	return either.RightEitherWrapper[error, A, B](rw.Get())
}

// ToOption transforms a Result into None on failure and Some on OK
// Play: https://go.dev/play/p/Uw1Zz6b952q
func (rw ResultWrapper[A, B]) ToOptionWrapper() option.OptionWrapper[A, B] {
	if rw.IsError() {
		return option.NoneOptionWrapper[A, B]()
	}

	return option.SomeOptionWrapper[A, B](rw.Get())
}

// ForEach executes the given side-effecting function if Result is valid.
func (rw ResultWrapper[A, B]) ForEach(mapper func(value A)) {
	if rw.IsOk() {
		mapper(rw.r.value)
	}
}

// Fold is an alias for Match which executes the first function if Result is valid and second function if invalid.
// It returns a new Result.
// Play: https://go.dev/play/p/-_eFaLJ31co
func (rw ResultWrapper[A, B]) Fold(onSuccess func(value A) (B, error), onError func(err error) (B, error)) Result[B] {
	return rw.Match(onSuccess, onError)
}

// Match executes the first function if Result is valid and second function if invalid.
// It returns a new Result.
// Play: https://go.dev/play/p/-_eFaLJ31co
func (rw ResultWrapper[A, B]) Match(onSuccess func(value A) (B, error), onError func(err error) (B, error)) Result[B] {
	if rw.IsError() {
		return TupleToResult(onError(rw.Error()))
	}
	return TupleToResult(onSuccess(rw.Get()))
}

// Map executes the mapper function if Result is valid. It returns a new Result.
// Play: https://go.dev/play/p/-ndpN_b_OSc
func (rw ResultWrapper[A, B]) Map(mapper func(value A) (B, error)) Result[B] {
	if rw.IsOk() {
		return TupleToResult(mapper(rw.Get()))
	}

	return Err[B](rw.Error())
}

// FlatMap executes the mapper function if Result is valid. It returns a new Result.
// Play: https://go.dev/play/p/Ud5QjZOqg-7
func (rw ResultWrapper[A, B]) FlatMap(mapper func(value A) Result[B]) Result[B] {
	if rw.IsOk() {
		return mapper(rw.Get())
	}

	return Err[B](rw.Error())
}
