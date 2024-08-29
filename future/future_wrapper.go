package future

type FutureWrapper[A, B any] struct {
	f Future[A]
}

func NewFutureWrapper[A, B any](cb func(resolve func(A), reject func(error))) *FutureWrapper[A, B] {
	fw := FutureWrapper[A, B]{
		f: Future[A]{
			cb:       cb,
			cancelCb: func() {},
			done:     make(chan struct{}),
		},
	}

	fw.f.active()

	return &fw
}

// Then is called when Future is resolved. It returns a new Future.
func (fw *FutureWrapper[A, B]) Then(cb func(A) (B, error)) *Future[B] {
	fw.f.mu.Lock()
	defer fw.f.mu.Unlock()

	fw.f.next = &Future[B]{
		cb: func(resolve func(B), reject func(error)) {
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
	case <-fw.f.done:
		fw.f.next.active()
	default:
	}
	return fw.f.next
}

// Catch is called when Future is rejected. It returns a new Future.
func (f *FutureWrapper[A, B]) Catch(cb func(error) (B, error)) *Future[B] {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.next = &Future[B]{
		cb: func(resolve func(B), reject func(error)) {
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
func (f *FutureWrapper[A, B]) Finally(cb func(A, error) (A, error)) *Future[B] {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.next = &Future[T]{
		cb: func(resolve func(T), reject func(error)) {
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
