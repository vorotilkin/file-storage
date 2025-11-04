package interfaces

import "context"

type Hooker interface {
	OnStart(context.Context) error
	OnStop(context.Context) error
}
