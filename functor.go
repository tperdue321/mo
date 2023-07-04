package mo

type OptionFunctor[T, U any] interface {
	Map(func (T) U) Option[U]
}

type OptionApplicative[T, U any] interface {
	OptionFunctor[T, U]

	Ap(ff Option[func (T) U], ft Option[T]) Option[U]
}

type OptionWrapper[T, U any] interface {
	OptionApplicative[T, U]

	// flatMap[A, B](fa: F[A])(f: (A) => F[B]): F[B]
	FlatMap(ft Option[T], f func(T) Option[U]) Option[U]
}

type OptionnMonad[T, U any] struct {
	o Option[T]
}

func PureOptionWrapper[T, U any](t T) OptionnMonad[T, U] {
	return OptionnMonad[T, U]{
		o: Option[T]{
			value: t,
			isPresent: true,
		},
	}
}

func (om *OptionnMonad[T, U]) Ap(ff Option[func(value T) U]) Option[U] {
	if ff.IsPresent() && om.o.IsPresent() {
		return EmptyableToOption(ff.value(om.o.value))
	}
	return None[U]()
}

func (om *OptionnMonad[T, U]) FlatMap(mapper func(value T) (Option[U])) Option[U] {
	val, present := om.o.Get()
	if present {
		return mapper(val)
	}
	return None[U]()
}

func (om *OptionnMonad[T, U]) Map(mapper func(value T) (U)) Option[U] {
	val, present := om.o.Get()
	if present {
		return EmptyableToOption(mapper(val))
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

type FutureFunctor[T, U any] struct {
	future *Future[T]
	next     *Future[U]
}

func (f *FutureFunctor[T, U]) Map(mapper func(T) (U, error)) *Future[U] {
	f.future.mu.Lock()
	defer f.future.mu.Unlock()

	f.next = &Future[U]{
		cb: func(resolve func(U), reject func(error)) {
			if f.future.result.IsError() {
				reject(f.future.result.Error())
				return
			}
			newValue, err := mapper(f.future.result.MustGet())
			if err != nil {
				reject(err)
				return
			}
			resolve(newValue)
		},
		cancelCb: func() {
			f.future.Cancel()
		},
		done: make(chan struct{}),
	}
	

	select {
	case <-f.future.done:
		f.next.active()
	default:
	}
	return f.next
}

type Functorr[A any, B any] interface {
	Option[A] | Either[B, A] | Result[A] | Future[A] | IO[A] | IOEither[A] | Task[A] | TaskEither[A]
}
