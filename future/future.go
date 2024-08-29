package future

import (
	"sync"

	"github.com/tperdue321/mo/either"
	"github.com/tperdue321/mo/result"
)

// NewFuture instanciate a new future.
func NewFuture[A any](cb func(resolve func(A), reject func(error))) *Future[A] {
	future := Future[A]{
		cb:       cb,
		cancelCb: func() {},
		done:     make(chan struct{}),
	}

	future.active()

	return &future
}

// Future represents a value which may or may not currently be available, but will be
// available at some point, or an exception if that value could not be made available.
type Future[A any] struct {
	mu sync.Mutex

	cb       func(func(A), func(error))
	cancelCb func()
	next     *Future[A]
	done     chan struct{}
	result   result.Result[A]
}

func (f *Future[A]) active() {
	go f.cb(f.resolve, f.reject)
}

func (f *Future[A]) activeSync() {
	f.cb(f.resolve, f.reject)
}

func (f *Future[A]) resolve(value A) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.result = result.Ok(value)
	if f.next != nil {
		f.next.activeSync()
	}
	close(f.done)
}

func (f *Future[A]) reject(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.result = Err[A](err)
	if f.next != nil {
		f.next.activeSync()
	}
	close(f.done)
}

// Ahen is called when Future is resolved. It returns a new Future.
func (f *Future[A]) Ahen(cb func(A) (A, error)) *Future[A] {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.next = &Future[A]{
		cb: func(resolve func(A), reject func(error)) {
			if f.result.IsError() {
				reject(f.result.Error())
				return
			}
			newValue, err := cb(f.result.MustGet())
			if err != nil {
				reject(err)
				return
			}
			resolve(newValue)
		},
		cancelCb: func() {
			f.Cancel()
		},
		done: make(chan struct{}),
	}

	select {
	case <-f.done:
		f.next.active()
	default:
	}
	return f.next
}

// Catch is called when Future is rejected. It returns a new Future.
func (f *Future[A]) Catch(cb func(error) (A, error)) *Future[A] {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.next = &Future[A]{
		cb: func(resolve func(A), reject func(error)) {
			if f.result.IsOk() {
				resolve(f.result.MustGet())
				return
			}
			newValue, err := cb(f.result.Error())
			if err != nil {
				reject(err)
				return
			}
			resolve(newValue)
		},
		cancelCb: func() {
			f.Cancel()
		},
		done: make(chan struct{}),
	}

	select {
	case <-f.done:
		f.next.active()
	default:
	}
	return f.next
}

// Finally is called when Future is processed either resolved or rejected. It returns a new Future.
func (f *Future[A]) Finally(cb func(A, error) (A, error)) *Future[A] {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.next = &Future[A]{
		cb: func(resolve func(A), reject func(error)) {
			newValue, err := cb(f.result.Get())
			if err != nil {
				reject(err)
				return
			}
			resolve(newValue)
		},
		cancelCb: func() {
			f.Cancel()
		},
		done: make(chan struct{}),
	}

	select {
	case <-f.done:
		f.next.active()
	default:
	}
	return f.next
}

// Cancel cancels the Future chain.
func (f *Future[A]) Cancel() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.next = nil
	if f.cancelCb != nil {
		f.cancelCb()
	}
}

// Collect awaits and return result of the Future.
func (f *Future[A]) Collect() (A, error) {
	<-f.done
	return f.result.Get()
}

// Result wraps Collect and returns a Result.
func (f *Future[A]) Result() result.Result[A] {
	return result.TupleToResult(f.Collect())
}

// Either wraps Collect and returns a Either.
func (f *Future[A]) Either() either.Either[error, A] {
	v, err := f.Collect()
	if err != nil {
		return either.Left[error, A](err)
	}
	return either.Right[error, A](v)
}
