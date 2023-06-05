package mo

// trait Functor[F[_]] {
// 	def map[A, B](fa: F[A])(f: A => B): F[B]
//   }
// func GetX[P interface { Point | Rect | Elli; GetX() int }] (p P) int {
// 	return p.GetX()
// }

// func Map[T, U any](option Option[T], mapper func(value T) (U, bool)) Option[U] {
// 	val, present := option.Get()
// 	if present {
// 		return TupleToOption(mapper(val))
// 	}
// 	return None[U]()
// }

type OptionFunctor[T, U any] struct {
	option Option[T]
}

func (f OptionFunctor[T, U]) Map(mapper func(value T) (U, bool)) Option[U] {
	val, present := f.option.Get()
	if present {
		return TupleToOption(mapper(val))
	}
	return None[U]()
}

type EitherFunctor[L, R1, R2 any] struct {
	either Either[L, R1]
}

func (f EitherFunctor[L, R1, R2]) Map(mapper func(value R1) (R2, bool), defaultLeft L) Either[L, R2] {
	left := Left[L, R2](defaultLeft)
	if r1, ok := f.either.Right(); ok {
		if r2, success := mapper(r1); success {
			return Right[L, R2](r2)
		}
	}
	return left
}

type ResultFunctor[T, U any] struct {
	result Result[T]
}

func (f ResultFunctor[T, U]) Map(mapper func(value T) (U, error)) Result[U] {
	if f.result.IsOk() {
		return TupleToResult[U](mapper(f.result.MustGet()))
	}
	return Err[U](f.result.Error())
}

// type FutureFunctor[T, U any] struct {
// 	future Future[T]
// }

// func (f FutureFunctor[T, U]) Map(mapper func(value T) (U)) Future[U] {
// 	val, err := f.future.Collect()
	
// 	NewFuture[U](mapper(val))
// }
// 	if f.result.IsOk() {
// 		return TupleToResult[U](mapper(f.result.MustGet()))
// 	}
// 	return Err[U](f.result.Error())
// }

// type Functorr[A any, B any] interface {
// 	Option[A] | Either[B, A] | Result[A] | Future[A] | IO[A] | IOEither[A] | Task[A] | TaskEither[A]
// }
