package domain

type Nullable[T any] struct {
	Value *T
	Set   bool
}

type Speaker interface {
	Speak() string
}
