package result

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultWrapperErr(t *testing.T) {
	is := assert.New(t)

	err := ErrWrapper[string, bool](assert.AnError)
	is.Equal(ResultWrapper[string, bool]{r: Result[string]{err: assert.AnError, value: "", isErr: true}}, err)
}

func TestResultWrapperOk(t *testing.T) {
	is := assert.New(t)

	ok := OkWrapper[string, bool]("foo")
	is.Equal(ResultWrapper[string, bool]{r: Result[string]{err: nil, value: "foo", isErr: false}}, ok)
}

func TestResultWrapperIsErrOrOk(t *testing.T) {
	is := assert.New(t)

	err := ErrWrapper[string, bool](assert.AnError)
	ok := OkWrapper[string, bool]("foo")

	is.True(err.IsError())
	is.False(err.IsOk())
	is.False(ok.IsError())
	is.True(ok.IsOk())
}

func TestResultWrapperErrOrOk(t *testing.T) {
	is := assert.New(t)

	err := ErrWrapper[string, bool](assert.AnError)
	ok := OkWrapper[string, bool]("foo")

	result1 := err.Error()
	result2 := err.Get()
	result3 := ok.Error()
	result4 := ok.Get()

	is.Equal(assert.AnError, result1)
	is.Equal("", result2)
	is.Equal(nil, result3)
	is.Equal("foo", result4)
}

func TestResultWrapperMustErrOrOk(t *testing.T) {
	is := assert.New(t)

	err := ErrWrapper[string, bool](assert.AnError)
	ok := OkWrapper[string, bool]("foo")

	is.Panics(func() {
		err.MustGet()
	})
	is.NotPanics(func() {
		is.Equal("foo", ok.MustGet())
	})
}

func TestResultWrapperGetOrElse(t *testing.T) {
	is := assert.New(t)

	err := ErrWrapper[string, bool](assert.AnError)
	ok := OkWrapper[string, bool]("foo")

	rightFallthrough := "bar"

	is.Equal("bar", err.GetOrElse(rightFallthrough))
	is.Equal("foo", ok.GetOrElse(rightFallthrough))
}

func TestResultWrapperForEach(t *testing.T) {
	is := assert.New(t)

	ErrWrapper[string, bool](assert.AnError).ForEach(
		func(b string) {
			is.Fail("should not enter here")
		},
	)

	OkWrapper[string, bool]("foobar").ForEach(
		func(b string) {
			is.Equal("foobar", b)
		},
	)
}

func TestResultWrapperMap(t *testing.T) {
	is := assert.New(t)

	err := ErrWrapper[string, bool](assert.AnError).Map(
		func(b string) (bool, error) {
			is.Fail("should not enter here")
			return true, nil
		},
	)
	testFunc := func (value string) (bool, error) {
		return value == "foo", nil
	}

	ok1 := OkWrapper[string, bool]("foo").Map(testFunc)
	ok2 := OkWrapper[string, bool]("foobar").Map(testFunc)

	is.Equal(Result[bool]{err: assert.AnError, value: false, isErr: true}, err)
	is.Equal(Result[bool]{err: nil, value: true, isErr: false}, ok1)
	is.Equal(Result[bool]{err: nil, value: false, isErr: false}, ok2)
}

func TestResultWrapperFlatMap(t *testing.T) {
	is := assert.New(t)

	testFunc := func (value string) bool {
		return value == "foo"
	}

	err := ErrWrapper[string, bool](assert.AnError).FlatMap(
		func(value string) (Result[bool]) {
			is.Fail("should not enter here")
			return Ok[bool](testFunc(value))
		},
	)

	ok1 := OkWrapper[string, bool]("foo").FlatMap(
		func(value string) Result[bool] {
			switch result := testFunc(value); result {
			case true:
				return Ok[bool](result)
			default:
				return Err[bool](assert.AnError)
			}
			
		},
	)

	ok2 := OkWrapper[string, bool]("foobar").FlatMap(
		func(value string) Result[bool] {
			switch result := testFunc(value); result {
			case true:
				return Ok[bool](result)
			default:
				return Err[bool](assert.AnError)
			}
			
		},
	)

	is.Equal(Result[bool]{err: assert.AnError, value: false, isErr: true}, err)
	is.Equal(Result[bool]{err: nil, value: true, isErr: false}, ok1)
	is.Equal(Result[bool]{err: assert.AnError, value: false, isErr: true}, ok2)
}

func TestResultWrapperApply(t *testing.T) {
	is := assert.New(t)

	f := func(value string) (int, error) {
		return strconv.Atoi(value)
	}


	resultFuncErr := Err[func(string) (int, error)](assert.AnError)
	e1 := ErrWrapper[string, int](assert.AnError).Apply(resultFuncErr)

	resultFuncOk := Ok[func(string) (int, error)](f)
	e2 := OkWrapper[string, int]("42").Apply(resultFuncOk)


	is.Equal(Result[int]{err: assert.AnError, value: 0, isErr: true}, e1)
	is.Equal(Result[int]{err: nil, value: 42, isErr: false}, e2)
}

func TestResultWrapperMonadicLaws(t *testing.T) {
	x := "42"

	ok := OkWrapper[string, int](x)

	f := func(value string) Result[int] {
		return TupleToResult(strconv.Atoi(value))
	}

	t.Run("Err Identity", func (t *testing.T) {
		is := assert.New(t)

		is.Equal(
			ok.FlatMap(f),
			f(x),
		)

	})

	t.Run("Ok Identity", func (t *testing.T) {
		y := 42
		is := assert.New(t)

		is.Equal(
			ok.FlatMap(func (value string) Result[int] {
				return TupleToResult(strconv.Atoi(value))
			}),
			Ok[int](y),
		)
	})

	t.Run("Associativity", func (t *testing.T) {
		is := assert.New(t)

		g := func(value int) Result[bool] {
			return Ok(value == 42)
		}

		associateResult := OkWrapper[string, bool](x)
		fa := func (x string) Result[bool] {
			return WrapResult[int, bool](f(x)).FlatMap(g)
		}

		// proves ok.flatMap(f).flatMap(g) == ok.flatMap(x => f(x).flatMap(g))
		is.Equal(
			WrapResult[int, bool](ok.FlatMap(f)).FlatMap(g),
			associateResult.FlatMap(fa),
		)

	})
}
