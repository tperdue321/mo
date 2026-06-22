//go:build !go1.22
// +build !go1.22

package mo

// Match executes the given function, depending of value is Left or Right, and returns result.
func (e Either[L, R]) Match(onLeft func(L) Either[L, R], onRight func(R) Either[L, R]) Either[L, R] {
	if e.IsLeft() {
		return onLeft(e.left)
	} else if e.IsRight() {
		return onRight(e.right)
	}

	panic(errEitherShouldBeLeftOrRight)
}

// MapRight executes the given function, if Either is of type Right, and returns result.
func (e Either[L, R]) MapRight(mapper func(R) Either[L, R]) Either[L, R] {
	if e.isLeft {
		return Left[L, R](e.left)
	} else if e.IsRight() {
		return mapper(e.right)
	}

	panic(errEitherShouldBeLeftOrRight)
}
