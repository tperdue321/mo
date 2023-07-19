package either

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEitherLeft(t *testing.T) {
	is := assert.New(t)

	left := Left[int, bool](42)
	is.Equal(Either[int, bool]{left: 42, right: false, isLeft: true}, left)
}

func TestEitherRight(t *testing.T) {
	is := assert.New(t)

	right := Right[int, bool](true)
	is.Equal(Either[int, bool]{left: 0, right: true, isLeft: false}, right)
}

func TestEitherIsLeftOrRight(t *testing.T) {
	is := assert.New(t)

	left := Left[int, bool](42)
	right := Right[int, bool](true)

	is.True(left.IsLeft())
	is.False(left.IsRight())
	is.False(right.IsLeft())
	is.True(right.IsRight())
}

func TestEitherLeftOrRight(t *testing.T) {
	is := assert.New(t)

	left := Left[int, bool](42)
	right := Right[int, bool](true)

	result1 := left.Left()
	result2 := left.Right()
	result3 := right.Left()
	result4 := right.Right()

	is.Equal(42, result1)
	is.Equal(false, result2)
	is.Equal(0, result3)
	is.Equal(true, result4)
}

func TestEitherMustLeftOrRight(t *testing.T) {
	is := assert.New(t)

	left := Left[int, bool](42)
	right := Right[int, bool](true)

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
		is.Equal(true, right.MustRight())
	})
}

func TestEitherGetOrElse(t *testing.T) {
	is := assert.New(t)

	left := Left[int, string](42)
	right := Right[int, string]("foobar")

	is.Equal(42, left.LeftOrElse(21))
	is.Equal(21, right.LeftOrElse(21))
	is.Equal("baz", left.RightOrElse("baz"))
	is.Equal("foobar", right.RightOrElse("baz"))
}

func TestEitherSwap(t *testing.T) {
	is := assert.New(t)

	left := Left[int, string](42)
	right := Right[int, string]("foobar")

	is.Equal(Either[string, int]{left: "", right: 42, isLeft: false}, left.Swap())
	is.Equal(Either[string, int]{left: "foobar", right: 0, isLeft: true}, right.Swap())
}

func TestEitherForEach(t *testing.T) {
	is := assert.New(t)

	Left[int, string](42).ForEach(
		func(b string) {
			is.Fail("should not enter here")
		},
	)

	Right[int, string]("foobar").ForEach(
		func(b string) {
			is.Equal("foobar", b)
		},
	)
}

func TestEitherMap(t *testing.T) {
	is := assert.New(t)

	e1 := Left[int, string](42).Map(
		func(b string) string {
			is.Fail("should not enter here")
			return "This is never reached"
		},
	)

	e2 := Right[int, string]("foobar").Map(
		func(b string) string {
			return b + "baz"
		},
	)

	is.Equal(Either[int, string]{left: 42, right: "", isLeft: true}, e1)
	is.Equal(Either[int, string]{left: 0, right: "foobarbaz", isLeft: false}, e2)
}

func TestEitherFlatMap(t *testing.T) {
	is := assert.New(t)

	e1 := Left[int, string](42).FlatMap(
		func(b string) Either[int, string] {
			is.Fail("should not enter here")
			return Right[int, string]("This is never reached")
		},
	)

	e2 := Right[int, string]("foobar").FlatMap(
		func(b string) Either[int, string] {
			return Right[int, string](b + "baz")
		},
	)

	is.Equal(Either[int, string]{left: 42, right: "", isLeft: true}, e1)
	is.Equal(Either[int, string]{left: 0, right: "foobarbaz", isLeft: false}, e2)
}
func TestEitherMonadicLaws(t *testing.T) {
	x := "foo"

	right := Right[int, string](x)

	f := func(right string) Either[int, string] {
		return Right[int, string](right)
	}

	g := func(right string) Either[int, string] {
		return Right[int, string](right + "bar")
	}

	t.Run("Left Identity", func (t *testing.T) {
		is := assert.New(t)

		is.Equal(f(x), right.FlatMap(f))

	})

	t.Run("Right Identity", func (t *testing.T) {
		x := "foo"
		is := assert.New(t)

		is.Equal(Right[int, string](x), right.FlatMap(func (right string) Either[int, string] {
			return Right[int, string](right)
		}))
	})

	t.Run("Associativity", func (t *testing.T) {
		is := assert.New(t)


		// proves either.flatMap(f).flatMap(g) == either.flatMap(x => f(x).flatMap(g))
		is.Equal(
			right.FlatMap(f).FlatMap(g),
			Right[int, string](x).FlatMap(func(x string) Either[int, string] {
				return f(x).FlatMap(g)
			}),
		)
	})
}
