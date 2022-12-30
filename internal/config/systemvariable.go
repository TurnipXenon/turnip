package config

import "time"

type genericSystemVariable[T any] struct {
	value      T
	expiration time.Time
}

func (s *genericSystemVariable[T]) Set(newValue T) {
	//TODO implement me
	panic("implement me")
}

type SystemVariable[T any] interface {
	Update()
	Get() T
	Set(newValue T)
}

// Update checks if expired and updates the valeu
func (s *genericSystemVariable[T]) Update() {
	// todo
}

func (s *genericSystemVariable[T]) Get() T {
	s.Update()
	return s.value
}

func NewGenericSystemVariable[T any](value T) SystemVariable[T] {
	return &genericSystemVariable[T]{
		value: value,
		// todo: time
	}
}
