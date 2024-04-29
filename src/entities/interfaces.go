package entities

import "context"

type Entity interface {
	GetId() string
}

// Wraps writing to DB functions for entities
type EntityWriter interface {
	Write(context.Context) error
}
