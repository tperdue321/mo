package option

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomeOptionWrapper(t *testing.T) {
	is := assert.New(t)

	optionWrapper := OptionWrapper[int, string]{
		o: Option[int]{
			value: 42,
			isPresent: true,
		},
	}

	is.Equal(
		optionWrapper,
		SomeOptionWrapper[int,string](42),
	)
}

func TestNoneOptionWrapper(t *testing.T) {
	is := assert.New(t)

	is.Equal(
		OptionWrapper[int,string]{o: Option[int]{isPresent: false}},
		NoneOptionWrapper[int,string](),
	)
}

func TestPointerToOptionWrapper(t *testing.T) {
	is := assert.New(t)
	
	t.Run("pointer to Some", func(t *testing.T) {
		optionWrapper := OptionWrapper[int, string]{
			o: Option[int]{
				value: 42,
				isPresent: true,
			},
		}

		is.Equal(
			optionWrapper,
			PointerToOptionWrapper[int,string](Pointer(42)),
		)
	})

	t.Run("pointer to None", func(t *testing.T) {
		optionWrapper := OptionWrapper[int, string]{
			o: Option[int]{
				isPresent: false,
			},
		}

		is.Equal(
			optionWrapper,
			PointerToOptionWrapper[int,string](nil),
		)
	})
}

func Test_OptionWrapperMap(t *testing.T) {
	is := assert.New(t)

	mappedValue := Option[string]{
		value: "42",
		isPresent: true,
	}

	mapper := func(value int) string {
		return strconv.Itoa(value)
	}

	baseOptionWrapper := SomeOptionWrapper[int, string](42)

	is.Equal(baseOptionWrapper.Map(mapper), mappedValue)

}

func Test_OptionWrapperFlatMap(t *testing.T) {
	is := assert.New(t)

	t.Run("Map to Some", func(t *testing.T) {
		mappedValue := Option[string]{
			value: "42",
			isPresent: true,
		}

		mapper := func(value int) Option[string] {
			return Some(strconv.Itoa(value))
		}

		baseOptionWrapper := SomeOptionWrapper[int, string](42)

		is.Equal(baseOptionWrapper.FlatMap(mapper), mappedValue)
	})
	
	t.Run("Map to None", func(t *testing.T) {

		mapper := func(value int) Option[string] {
			return None[string]()
		}

		baseOptionWrapper := SomeOptionWrapper[int, string](42)

		is.Equal(baseOptionWrapper.FlatMap(mapper), None[string]())
	})
}

func TestOptionWrapperApply(t *testing.T) {
	is := assert.New(t)

	t.Run("Map to Some", func(t *testing.T) {
		mappedValue := Option[string]{
			value: "42",
			isPresent: true,
		}

		someMapper := Some[func(int) string](
			func(value int) string {
				return strconv.Itoa(value)
			},
		)

		baseOptionWrapper := SomeOptionWrapper[int, string](42)

		is.Equal(mappedValue, baseOptionWrapper.Apply(someMapper))
	})
	
	t.Run("Map to None", func(t *testing.T) {

		noneMapper := None[func(int) string]()

		baseOptionWrapper := SomeOptionWrapper[int, string](42)

		is.Equal(baseOptionWrapper.Apply(noneMapper), None[string]())
	})
}

func TestOptionWrapperMonadicLaws(t *testing.T) {
	x := 42

	some := SomeOptionWrapper[int, string](x)

	f := func(value int) Option[string] {
		return Some(strconv.Itoa(value))
	}

	t.Run("Left Identity", func (t *testing.T) {
		is := assert.New(t)

		is.Equal(
			some.FlatMap(f),
			f(x),
		)

	})

	t.Run("Right Identity", func (t *testing.T) {
		y := "foo"
		is := assert.New(t)


		someStrToStr := SomeOptionWrapper[string, string](y)
		is.Equal(
			someStrToStr.FlatMap(func (value string) Option[string] {
				return Some(value)
			}),
			Some(y),
		)
	})

	t.Run("Associativity", func (t *testing.T) {
		is := assert.New(t)

		g := func(value string) Option[bool] {
			switch result := value == "foo"; result {
			case true:
				return Some[bool](result)
			default:
				return None[bool]()
			}
		}

		associateSome := SomeOptionWrapper[int, bool](x)
		fa := func (value int) Option[bool] {
			return WrapOption[string, bool](f(value)).FlatMap(g)
		}

		// proves option.flatMap(f).flatMap(g) == option.flatMap(x => f(x).flatMap(g))
		is.Equal(
			WrapOption[string, bool](some.FlatMap(f)).FlatMap(g),
			associateSome.FlatMap(fa),
		)

	})
}
