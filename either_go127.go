//go:build go1.27
// +build go1.27

package mo

// Match executes the given function, depending of value is Left or Right, and returns result.
func (e Either[L, R]) Match[U any](onLeft func(L) Either[L, U], onRight func(R) Either[L, U]) Either[L, U] {
	if e.IsLeft() {
		return onLeft(e.left)
	} else if e.IsRight() {
		return onRight(e.right)
	}

	panic(errEitherShouldBeLeftOrRight)
}


// MapRight executes the given function, if Either is of type Right, and returns result.
func (e Either[L, R]) MapRight[U any](mapper func(R) Either[L, U]) Either[L, U] {
	if e.isLeft {
		return Left[L, U](e.left)
	} else if e.IsRight() {
		return mapper(e.right)
	}

	panic(errEitherShouldBeLeftOrRight)
}
