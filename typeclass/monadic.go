package typeclass

type Monadic[T, U any] interface {
	ForEach(func(T))
	Map(func(T) U) Monadable[U]
	FlatMap(func(T) Monadable[U]) Monadable[U]
	Ap(Monadable[func(T) U]) Monadable[U]
}
