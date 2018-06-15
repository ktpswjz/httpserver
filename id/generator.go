package id

const (
	maxUint64 = ^uint64(0) - 1
)

type Generator interface {
	Max() uint64
	New() uint64
	NewAsString() string
}

func NewNormal() Generator {
	return &innerUint64{id: 0}
}

func NewTime() Generator {
	return &innerTime{
		id:  0,
		max: 0,
	}
}
