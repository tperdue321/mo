package option

// Pointer returns a pointer to a variable holding the supplied T constant
func Pointer[T any](x T) *T {
	return &x
}

type OptionWrapper[T, U any] struct {
	o Option[T]
}

// Some builds an Option when value is present.
// Play: https://go.dev/play/p/iqz2n9n0tDM
func SomeOptionWrapper[T, U any](value T) OptionWrapper[T, U] {
	return OptionWrapper[T, U]{
		o: Option[T]{
			isPresent: true,
			value:     value,
		},
	}
}

// Some builds an Option when value is present.
// Play: https://go.dev/play/p/iqz2n9n0tDM
func NoneOptionWrapper[T, U any]() OptionWrapper[T, U] {
	return OptionWrapper[T, U]{
		o: Option[T]{
			isPresent: false,
		},
	}
}



// PointerToOption builds a Some Option when value is not nil, or None.
// Play: https://go.dev/play/p/yPVMj4DUb-I
func PointerToOptionWrapper[T, U any](value *T) OptionWrapper[T, U] {
	if value == nil {
		return NoneOptionWrapper[T, U]()
	}

	return SomeOptionWrapper[T, U](*value)
}

func WrapOption[T, U any](option Option[T]) OptionWrapper[T, U] {
	return OptionWrapper[T, U]{
		o: option,
	}
}

func (om *OptionWrapper[T, U]) Ap(ff Option[func(value T) U]) Option[U] {
	if ff.IsPresent() {
		return om.Map(ff.value)
	}
	return None[U]()
}

func (om *OptionWrapper[T, U]) FlatMap(mapper func(value T) Option[U]) Option[U] {
	if value, ok := om.o.Get(); ok {
		return mapper(value)
	}
	return None[U]()
}

func (om *OptionWrapper[T, U]) Map(mapper func(value T) (U)) Option[U] {
	if value, ok := om.o.Get(); ok {
		return PointerToOption(Pointer(mapper(value)))
	}
	return None[U]()
}
