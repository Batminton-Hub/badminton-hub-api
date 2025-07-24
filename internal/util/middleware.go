package core_util

type MiddlewareUtil struct {
}

func NewMiddlewareUtil() *MiddlewareUtil {
	return &MiddlewareUtil{}
}

func (m *MiddlewareUtil) Authenticate(token string) error {
	// Implement authentication logic here
	// For example, check if the token exists in the database
	return nil // Return nil if authentication is successful
}
