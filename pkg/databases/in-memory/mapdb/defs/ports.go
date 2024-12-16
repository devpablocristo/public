package defs

type Repository interface {
	GetDb() map[string]any
}
