package port

type MiddlewareUtil interface {
	Authenticate(token string) error
}
