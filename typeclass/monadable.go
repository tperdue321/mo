package typeclass

type Monadable[T any] interface {
	Lift(T) Monadable[T]
	Map(func(T) T) Monadable[T]
	FlatMap(func(T) Monadable[T]) Monadable[T]
}
