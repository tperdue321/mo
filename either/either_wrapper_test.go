package either

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEitherWrapperLeft(t *testing.T) {
	is := assert.New(t)

	left := LeftEitherWrapper[int, string, bool](42)
	is.Equal(EitherWrapper[int, string, bool]{e: Either[int, string]{left: 42, right: "", isLeft: true}}, left)
}

func TestEitherWrapperRight(t *testing.T) {
	is := assert.New(t)

	right := RightEitherWrapper[int, string, bool]("foo")
	is.Equal(EitherWrapper[int, string, bool]{e: Either[int, string]{left: 0, right: "foo", isLeft: false}}, right)
}

func TestEitherWrapperIsLeftOrRight(t *testing.T) {
	is := assert.New(t)

	left := LeftEitherWrapper[int, string, bool](42)
	right := RightEitherWrapper[int, string, bool]("foo")

	is.True(left.IsLeft())
	is.False(left.IsRight())
	is.False(right.IsLeft())
	is.True(right.IsRight())
}

func TestEitherWrapperLeftOrRight(t *testing.T) {
	is := assert.New(t)

	left := LeftEitherWrapper[int, string, bool](42)
	right := RightEitherWrapper[int, string, bool]("foo")

	result1 := left.Left()
	result2 := left.Right()
	result3 := right.Left()
	result4 := right.Right()

	is.Equal(42, result1)
	is.Equal("", result2)
	is.Equal(0, result3)
	is.Equal("foo", result4)
}

func TestEitherWrapperMustLeftOrRight(t *testing.T) {
	is := assert.New(t)

	left := LeftEitherWrapper[int, string, bool](42)
	right := RightEitherWrapper[int, string, bool]("foo")

	is.NotPanics(func() {
		is.Equal(42, left.MustLeft())
	})
	is.Panics(func() {
		left.MustRight()
	})
	is.Panics(func() {
		right.MustLeft()
	})
	is.NotPanics(func() {
		is.Equal("foo", right.MustRight())
	})
}

func TestEitherWrapperGetOrElse(t *testing.T) {
	is := assert.New(t)

	left := LeftEitherWrapper[int, string, bool](42)
	right := RightEitherWrapper[int, string, bool]("foo")

	leftFallthrough := 21
	rightFallthrough := "bar"

	is.Equal(42, left.LeftOrElse(leftFallthrough))
	is.Equal(21, right.LeftOrElse(leftFallthrough))
	is.Equal("bar", left.RightOrElse(rightFallthrough))
	is.Equal("foo", right.RightOrElse(rightFallthrough))
}

func TestEitherWrapperSwap(t *testing.T) {
	is := assert.New(t)

	left := LeftEitherWrapper[int, string, bool](42)
	right := RightEitherWrapper[int, string, bool]("foo")

	is.Equal(EitherWrapper[string, int, bool]{e: Either[string, int]{left: "", right: 42, isLeft: false}}, left.Swap())
	is.Equal(EitherWrapper[string, int, bool]{e: Either[string, int]{left: "foo", right: 0, isLeft: true}}, right.Swap())
}

func TestEitherWrapperForEach(t *testing.T) {
	is := assert.New(t)

	LeftEitherWrapper[int, string, bool](42).ForEach(
		func(b string) {
			is.Fail("should not enter here")
		},
	)

	RightEitherWrapper[int, string, bool]("foobar").ForEach(
		func(b string) {
			is.Equal("foobar", b)
		},
	)
}

func TestEitherWrapperMap(t *testing.T) {
	is := assert.New(t)

	e1 := LeftEitherWrapper[int, string, bool](42).Map(
		func(b string) bool {
			is.Fail("should not enter here")
			return true
		},
	)

	e2 := RightEitherWrapper[int, string, bool]("foo").Map(testFunc)
	e3 := RightEitherWrapper[int, string, bool]("foobar").Map(testFunc)

	is.Equal(Either[int, bool]{left: 42, right: false, isLeft: true}, e1)
	is.Equal(Either[int, bool]{left: 0, right: true, isLeft: false}, e2)
	is.Equal(Either[int, bool]{left: 0, right: false, isLeft: false}, e3)
}

func TestEitherWrapperFlatMap(t *testing.T) {
	is := assert.New(t)

	e1 := LeftEitherWrapper[int, string, bool](42).FlatMap(
		func(right string) Either[int, bool] {
			is.Fail("should not enter here")
			return Right[int, bool](testFunc(right))
		},
	)

	e2 := RightEitherWrapper[int, string, bool]("foo").FlatMap(
		func(right string) Either[int, bool] {
			switch result := testFunc(right); result {
			case true:
				return Right[int, bool](result)
			default:
				return Left[int, bool](-1)
			}
			
		},
	)

	e3 := RightEitherWrapper[int, string, bool]("foobar").FlatMap(
		func(right string) Either[int, bool] {
			switch result := testFunc(right); result {
			case true:
				return Right[int, bool](result)
			default:
				return Left[int, bool](-1)
			}
			
		},
	)

	is.Equal(Either[int, bool]{left: 42, right: false, isLeft: true}, e1)
	is.Equal(Either[int, bool]{left: 0, right: true, isLeft: false}, e2)
	is.Equal(Either[int, bool]{left: -1, right: false, isLeft: true}, e3)
}

func TestEitherWrapperApply(t *testing.T) {
	is := assert.New(t)

	eitherFuncLeft := Left[int, func(string) bool](42)
	e1 := LeftEitherWrapper[int, string, bool](42).Apply(eitherFuncLeft)

	eitherFuncRight := Right[int, func(string) bool](testFunc)
	e2 := RightEitherWrapper[int, string, bool]("foo").Apply(eitherFuncRight)

	is.Equal(Either[int, bool]{left: 42, right: false, isLeft: true}, e1)
	is.Equal(Either[int, bool]{left: 0, right: true, isLeft: false}, e2)
}

func TestEitherWrapperMonadicLaws(t *testing.T) {
	x := "foo"

	right := RightEitherWrapper[int, string, string](x)

	f := func(right string) Either[int, string] {
		switch result := testFunc(right); result {
		case true:
			return Right[int, string](right)
		default:
			return Left[int, string](-1)
		}
	}

	t.Run("Left Identity", func (t *testing.T) {
		is := assert.New(t)

		is.Equal(
			right.FlatMap(f),
			f(x),
		)

	})

	t.Run("Right Identity", func (t *testing.T) {
		x := "foo"
		is := assert.New(t)

		is.Equal(
			right.FlatMap(func (right string) Either[int, string] {
				return Right[int, string](right)
			}),
			Right[int, string](x),
		)
	})

	t.Run("Associativity", func (t *testing.T) {
		is := assert.New(t)

		g := func(right string) Either[int, bool] {
			switch result := testFunc(right); result {
			case true:
				return Right[int, bool](result)
			default:
				return Left[int, bool](-1)
			}
		}

		associateEither := RightEitherWrapper[int, string, bool]("foo")
		fa := func (x string) Either[int, bool] {
			return WrapEither[int, string, bool](f(x)).FlatMap(g)
		}

		// proves either.flatMap(f).flatMap(g) == either.flatMap(x => f(x).flatMap(g))
		is.Equal(
			WrapEither[int, string, bool](right.FlatMap(f)).FlatMap(g),
			associateEither.FlatMap(fa),
		)

	})
}


func testFunc(right string) bool {
	return right == "foo"
}
