package mo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultOk(t *testing.T) {
	is := assert.New(t)

	is.Equal(Result[int]{value: 42, isErr: false, err: nil}, Ok(42))
}

func TestResultErr(t *testing.T) {
	is := assert.New(t)

	is.Equal(Result[int]{value: 0, isErr: true, err: assert.AnError}, Err[int](assert.AnError))
}

func TestResultTupleToResult(t *testing.T) {
	is := assert.New(t)

	is.Equal(Result[int]{value: 0, isErr: true, err: assert.AnError}, TupleToResult(42, assert.AnError))
}

func TestResultTry(t *testing.T) {
	is := assert.New(t)

	is.Equal(Result[int]{value: 42, isErr: false, err: nil}, Try(func() (int, error) {
		return 42, nil
	}))
	is.Equal(Result[int]{value: 0, isErr: true, err: assert.AnError}, Try(func() (int, error) {
		return 42, assert.AnError
	}))
}

func TestResultIsOk(t *testing.T) {
	is := assert.New(t)

	is.True(Ok(42).IsOk())
	is.False(Err[int](assert.AnError).IsOk())
}

func TestResultIsError(t *testing.T) {
	is := assert.New(t)

	is.False(Ok(42).IsError())
	is.True(Err[int](assert.AnError).IsError())
}

func TestResultError(t *testing.T) {
	is := assert.New(t)

	is.Nil(Ok(42).Error())
	is.NotNil(Err[int](assert.AnError).Error())
	is.Equal(assert.AnError, Err[int](assert.AnError).Error())
}

func TestResultGet(t *testing.T) {
	is := assert.New(t)

	v1, err1 := Ok(42).Get()
	v2, err2 := Err[int](assert.AnError).Get()

	is.Equal(42, v1)
	is.Nil(err1)
	is.Error(assert.AnError, err1)

	is.Equal(0, v2)
	is.NotNil(err2)
	is.Error(assert.AnError, err2)
}

func TestResultMustGet(t *testing.T) {
	is := assert.New(t)

	is.NotPanics(func() {
		Ok(42).MustGet()
	})
	is.Panics(func() {
		Err[int](assert.AnError).MustGet()
	})

	is.Equal(42, Ok(42).MustGet())
}

func TestResultOrElse(t *testing.T) {
	is := assert.New(t)

	is.Equal(42, Ok(42).OrElse(21))
	is.Equal(21, Err[int](assert.AnError).OrElse(21))
}

func TestResultOrEmpty(t *testing.T) {
	is := assert.New(t)

	is.Equal(42, Ok(42).OrEmpty())
	is.Equal(0, Err[int](assert.AnError).OrEmpty())
}

func TestResultToEither(t *testing.T) {
	is := assert.New(t)

	right, ok1 := Ok(42).ToEither().Right()
	left, ok2 := Err[int](assert.AnError).ToEither().Left()

	is.Equal(42, right)
	is.True(ok1)
	is.Equal(assert.AnError, left)
	is.True(ok2)
}

func TestResultForEach(t *testing.T) {
	is := assert.New(t)

	Err[int](assert.AnError).ForEach(func(i int) {
		is.Fail("should not enter here")
	})

	Ok(42).ForEach(func(i int) {
		is.Equal(42, i)
	})
}

func TestResultMatch(t *testing.T) {
	is := assert.New(t)

	opt1 := Ok(21).Match(
		func(i int) (int, error) {
			is.Equal(21, i)
			return i * 2, nil
		},
		func(err error) (int, error) {
			is.Fail("should not enter here")
			return 0, err
		},
	)
	opt2 := Err[int](assert.AnError).Match(
		func(i int) (int, error) {
			is.Fail("should not enter here")
			return i * 2, nil
		},
		func(err error) (int, error) {
			is.Equal(assert.AnError, err)
			return 0, err
		},
	)

	is.Equal(Result[int]{value: 42, isErr: false, err: nil}, opt1)
	is.Equal(Result[int]{value: 0, isErr: true, err: assert.AnError}, opt2)
}

func TestResultMap(t *testing.T) {
	is := assert.New(t)

	opt1 := Ok(21).Map(func(i int) (int, error) {
		return i * 2, nil
	})
	opt2 := Err[int](assert.AnError).Map(func(i int) (int, error) {
		is.Fail("should not be called")
		return 42, nil
	})

	is.Equal(Result[int]{value: 42, isErr: false, err: nil}, opt1)
	is.Equal(Result[int]{value: 0, isErr: true, err: assert.AnError}, opt2)
}

func TestResultMapErr(t *testing.T) {
	is := assert.New(t)

	opt1 := Ok(21).MapErr(func(err error) (int, error) {
		is.Fail("should not be called")
		return 42, nil
	})
	opt2 := Err[int](assert.AnError).MapErr(func(err error) (int, error) {
		return 42, nil
	})

	is.Equal(Result[int]{value: 21, isErr: false, err: nil}, opt1)
	is.Equal(Result[int]{value: 42, isErr: false, err: nil}, opt2)
}

func TestResultFlatMap(t *testing.T) {
	is := assert.New(t)

	opt1 := Ok(21).FlatMap(func(i int) Result[int] {
		return Ok(42)
	})
	opt2 := Err[int](assert.AnError).FlatMap(func(i int) Result[int] {
		is.Fail("should not be called")
		return Ok(42)
	})

	is.Equal(Result[int]{value: 42, isErr: false, err: nil}, opt1)
	is.Equal(Result[int]{value: 0, isErr: true, err: assert.AnError}, opt2)
}
