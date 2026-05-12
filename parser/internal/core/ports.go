package core

import "context"

type DB interface {
	add(context.Context) (int, error)
}

type Parser interface {
	Parse()
}
