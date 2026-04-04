package utils

import (
	"context"
	"errors"
	"sync"
)

type call[T any] struct {
	err   error
	value T
	done  chan struct{}
}

type SingleFlight[T any] struct {
	mutex sync.Mutex
	calls map[string]*call[T]
}

func NewSingleFlight[T any]() *SingleFlight[T] {
	return &SingleFlight[T]{
		calls: make(map[string]*call[T]),
	}
}

// Do объединяет параллельные запросы с одним и тем же ключом,
// чтобы только один поток выполнил expensive action, а остальные
// дождались его результата. Экземпляр SingleFlight хранит один
// конкретный тип T, поэтому наружу не требуется приведение типов.
func (s *SingleFlight[T]) Do(ctx context.Context, key string, action func(context.Context) (T, error)) (T, error) {
	var zero T

	s.mutex.Lock()
	if call, found := s.calls[key]; found {
		s.mutex.Unlock()
		return s.wait(ctx, call)
	}

	call := &call[T]{
		done: make(chan struct{}),
	}

	s.calls[key] = call
	s.mutex.Unlock()

	go func() {
		defer func() {
			if v := recover(); v != nil {
				call.err = errors.New("error from single flight")
				call.value = zero
			}

			close(call.done)

			s.mutex.Lock()
			delete(s.calls, key)
			s.mutex.Unlock()
		}()

		call.value, call.err = action(ctx)
	}()

	return s.wait(ctx, call)
}

func (s *SingleFlight[T]) wait(ctx context.Context, call *call[T]) (T, error) {
	var zero T

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case <-call.done:
		if call.err != nil {
			return zero, call.err
		}

		return call.value, nil
	}
}
