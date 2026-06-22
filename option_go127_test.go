//go:build go1.27
// +build go1.27

package mo

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionWrapperMonadicLaws(t *testing.T) {
	x := 42

	some := Some(x)

	f := func(value int) Option[string] {
		return Some(strconv.Itoa(value))
	}

	t.Run("Left Identity", func(t *testing.T) {
		is := assert.New(t)

		is.Equal(
			some.FlatMap(f),
			f(x),
		)

	})

	t.Run("Right Identity", func(t *testing.T) {
		y := "foo"
		is := assert.New(t)

		someStrToStr := Some[string](y)
		is.Equal(
			someStrToStr.FlatMap(func(value string) Option[string] {
				return Some(value)
			}),
			Some(y),
		)
	})

	t.Run("Associativity", func(t *testing.T) {
		is := assert.New(t)

		g := func(value string) Option[bool] {
			switch result := value == "foo"; result {
			case true:
				return Some[bool](result)
			default:
				return None[bool]()
			}
		}

		associateSome := Some[int](x)
		fa := func(value int) Option[bool] {
			return f(value).FlatMap(g)
		}

		// proves option.flatMap(f).flatMap(g) == option.flatMap(x => f(x).flatMap(g))
		is.Equal(
			(some.FlatMap(f)).FlatMap(g),
			associateSome.FlatMap(fa),
		)

	})
}
