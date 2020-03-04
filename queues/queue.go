package queues

import (
	tos "github.com/vanhtuan0409/scheduler"
)

type Queue interface {
	Name() string
	Enqueue(*tos.Task)
}
