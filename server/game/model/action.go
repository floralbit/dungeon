package model

type Action interface {
	Execute() bool
}
