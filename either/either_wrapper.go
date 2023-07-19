package either

import "github.com/tperdue321/mo/option"

type EitherWrapper[L, R, T any] struct {
	e Either[L, R]
}

func RightEitherWrapper[L, R, T any](right R) EitherWrapper[L, R, T] {
	return EitherWrapper[L, R, T]{
		e: Right[L, R](right),
	}
}

func LeftEitherWrapper[L, R, T any](left L) EitherWrapper[L, R, T] {
	return EitherWrapper[L, R, T]{
		e: Left[L, R](left),
	}
}

func WrapEither[L, R, T any](either Either[L, R]) EitherWrapper[L, R, T] {
	return EitherWrapper[L, R, T]{
		e: either,
	}
}

func (ew EitherWrapper[L, R, T]) IsLeft() bool {
	return ew.e.IsLeft()
}

func (ew EitherWrapper[L, R, T]) IsRight() bool {
	return ew.e.IsRight()
}

func (ew EitherWrapper[L, R, T]) Left() L {
	return ew.e.Left()
}

func (ew EitherWrapper[L, R, T]) Right() R {
	return ew.e.Right()
}

func (ew EitherWrapper[L, R, T]) MustLeft() L {
	if !ew.e.IsLeft() {
		panic(eitherMissingLeftValue)
	}

	return ew.e.Left()
}

func (ew EitherWrapper[L, R, T]) MustRight() R {
	if ew.e.IsLeft() {
		panic(eitherMissingRightValue)
	}

	return ew.e.Right()
}

// LeftOrElse returns left value of a Either struct or fallback.
func (ew EitherWrapper[L, R, T]) LeftOrElse(fallback L) L {
	if ew.IsLeft() {
		return ew.e.left
	}

	return fallback
}

// GetOrElse returns right value of a Either struct or fallback.
func (ew EitherWrapper[L, R, T]) RightOrElse(fallback R) R {
	if ew.e.IsRight() {
		return ew.e.right
	}

	return fallback
}

// Swap takes a Left and makes it a Right or a Right and makes it a Left
func (ew EitherWrapper[L, R, T]) Swap() EitherWrapper[R, L, T] {
	if ew.e.IsLeft() {
		return RightEitherWrapper[R, L, T](ew.e.left)
	}

	return LeftEitherWrapper[R, L, T](ew.e.right)
}

// ForEach executes the given side-effecting function on a Right value.
// Either is Right biased, if there is a desire to map over left values, `Swap` must be called first
func (ew EitherWrapper[L, R, T]) ForEach(rightCb func(R)) {
	if ew.e.IsRight() {
		rightCb(ew.e.right)
	}
}

// ForEach executes the given side-effecting function on a Right value.
func (ew EitherWrapper[L, R, T]) Apply(ff Either[L, func(value R) T]) Either[L, T] {
	if ff.IsRight() {
		return ew.Map(ff.Right())
	}
	return Left[L, T](ff.Left())
}

func (ew EitherWrapper[L, R, T]) FlatMap(mapper func(value R) Either[L, T]) Either[L, T] {
	if ew.e.IsRight() {
		return mapper(ew.Right())
	}
	return Left[L, T]((ew.Left()))
}

func (ew EitherWrapper[L, R, T]) Map(mapper func(value R) T) Either[L, T] {
	if ew.e.IsRight() {
		return Right[L, T](mapper(ew.Right()))
	}
	return Left[L, T]((ew.Left()))
}

func (ew EitherWrapper[L, R, T]) ToOptionWrapper() option.OptionWrapper[R, T] {
	if ew.IsRight() {
		return option.SomeOptionWrapper[R, T](ew.Right())
	}

	return option.NoneOptionWrapper[R, T]()
} 

func FlattenEitherWrapper[L1, L2, R, T1, T2 any](e EitherWrapper[L1, EitherWrapper[L2, R, T2], T1]) EitherWrapper[L2, R, T2] {
	return e.Right()
}
