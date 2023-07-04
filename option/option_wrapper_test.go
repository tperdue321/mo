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

func TestOptionWrapperAp(t *testing.T) {
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

		is.Equal(mappedValue, baseOptionWrapper.Ap(someMapper))
	})
	
	t.Run("Map to None", func(t *testing.T) {

		noneMapper := None[func(int) string]()

		baseOptionWrapper := SomeOptionWrapper[int, string](42)

		is.Equal(baseOptionWrapper.Ap(noneMapper), None[string]())
	})
}
