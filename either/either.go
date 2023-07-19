package either

import "fmt"

var eitherShouldBeLeftOrRight = fmt.Errorf("either should be Left or Right")
var eitherMissingLeftValue = fmt.Errorf("no such Left value")
var eitherMissingRightValue = fmt.Errorf("no such Right value")

type Eitherish[L, R any] interface {
	ForEach(func(R))
	Map(func(R) R) Either[L, R]
	FlatMap(func(R) Either[L, R]) Either[L, R]
}

// Left builds the left side of the Either struct, as opposed to the Right side.
func Left[L any, R any](value L) Either[L, R] {
	return Either[L, R]{
		isLeft: true,
		left:   value,
	}
}

// Right builds the right side of the Either struct, as opposed to the Left side.
func Right[L any, R any](value R) Either[L, R] {
	return Either[L, R]{
		isLeft: false,
		right:  value,
	}
}

// Either respresents a value of 2 possible types.
// An instance of Either is an instance of either A or B.
// Either is Right biased in order to allow for monadic functions
// In ordr to map over left biased values, call Left to generate a LeftProjection
type Either[L any, R any] struct {
	isLeft bool

	left  L
	right R
}

// IsLeft returns true if Either is an instance of Left.
func (e Either[L, R]) IsLeft() bool {
	return e.isLeft
}

// IsRight returns true if Either is an instance of Right.
func (e Either[L, R]) IsRight() bool {
	return !e.isLeft
}

// Left returns left value of a Either struct.
func (e Either[L, R]) Left() L {
	return e.left
}

// Right returns right value of a Either struct.
func (e Either[L, R]) Right() R {
	return e.right
}

// MustLeft returns left value of a Either struct or panics.
func (e Either[L, R]) MustLeft() L {
	if !e.IsLeft() {
		panic(eitherMissingLeftValue)
	}

	return e.left
}

// MustRight returns right value of a Either struct or panics.
func (e Either[L, R]) MustRight() R {
	if !e.IsRight() {
		panic(eitherMissingRightValue)
	}

	return e.right
}

// LeftOrElse returns left value of a Either struct or fallback.
func (e Either[L, R]) LeftOrElse(fallback L) L {
	if e.IsLeft() {
		return e.left
	}

	return fallback
}

// GetOrElse returns right value of a Either struct or fallback.
func (e Either[L, R]) RightOrElse(fallback R) R {
	if e.IsRight() {
		return e.right
	}

	return fallback
}

// Swap takes a Left and makes it a Right or a Right and makes it a Left
func (e Either[L, R]) Swap() Either[R, L] {
	if e.IsLeft() {
		return Right[R, L](e.left)
	}

	return Left[R, L](e.right)
}

// ForEach executes the given side-effecting function on a Right value.
// Either is Right biased, if there is a desire to map over left values, `Swap` must be called first
func (e Either[L, R]) ForEach(rightCb func(R)) {
	if e.IsRight() {
		rightCb(e.right)
	}
}

// Map executes the given function, if Either is of type Right, and returns result.
// Either is Right biased, if there is a desire to map over left values, `Swap` must be called first
func (e Either[L, R]) Map(mapper func(R) R) Either[L, R] {
	if e.isLeft {
		return Left[L, R](e.left)
	} else {
		return Right[L, R](mapper(e.right))
	}
}

// FlatMap executes the given function, if Either is of type Right, and returns result.
// Either is Right biased, if there is a desire to flatmap over left values, `Swap` must be called first
func (e Either[L, R]) FlatMap(mapper func(R) Either[L, R]) Either[L, R] {
	if e.isLeft {
		return Left[L, R](e.left)
	} else {
		return mapper(e.right)
	}
}


func FlattenEither[L1, L2, R any](e Either[L1, Either[L2, R]]) Either[L2, R] {
	return e.Right()
}
