package app_handlers

type AppHandler[T any, U any] interface {
	Handle(request T) (*U, error)
}
