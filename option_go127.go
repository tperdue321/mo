//go:build go1.27
// +build go1.27

package mo

// Map executes the mapper function if value is present or returns None if absent.
// Play: https://go.dev/play/p/mvfP3pcP_eJ
func (o Option[T]) Map[U any](mapper func(value T) (U, bool)) Option[U] {
if o.isPresent {
		return TupleToOption(mapper(o.value))
	}

	return None[U]()
}

// FlatMap executes the mapper function if value is present or returns None if absent.
// Play: https://go.dev/play/p/OXO-zJx6n5r
func (o Option[T]) FlatMap[U any](mapper func(value T) Option[U]) Option[U] {
	if o.isPresent {
		return mapper(o.value)
	}

	return None[U]()
}

// MapValue executes the mapper function if value is present or returns None if absent.
func (o Option[T]) MapValue[U any](mapper func(value T) U) Option[U] {
	if o.isPresent {
		return Some(mapper(o.value))
	}

	return None[U]()
}

// Match executes the first function if value is present and second function if absent.
// It returns a new Option.
// Play: https://go.dev/play/p/1V6st3LDJsM
func (o Option[T]) Match[U any](onValue func(value T) (U, bool), onNone func() (U, bool)) Option[U] {
	if o.isPresent {
		return TupleToOption(onValue(o.value))
	}
	return TupleToOption(onNone())
}

