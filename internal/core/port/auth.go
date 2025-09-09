package port

import "Badminton-Hub/internal/core/domain"

// type MiddlewareService interface {
// 	Authentication
// 	MiddlewareAuthUtil
// 	MiddlewareCallback
// }

// type Authentication interface {
// 	Authenticate(token string) (int, domain.AuthResponse)
// }

// type MiddlewareAuthUtil interface {
// 	// Encryption() EncryptionUtil
// 	GenBearerToken(hashBody domain.HashAuth) (domain.BearerToken, error)
// 	ValidateBearerToken(token domain.BearerToken) (domain.AuthBody, error)
// }

// type MiddlewareCallback interface {
// 	GoogleLoginCallback(state, code string) (int, domain.ResponseGoogleLoginCallback)
// 	GoogleRegisterCallback(state, code string) (int, domain.ResponseGoogleRegisterCallback)
// }

//---New Flow---//
type AuthenticationSystem interface {
	MiddlewareService
	AuthenticationService
}

type AuthenticationService interface {
	Login(loginInfo domain.LoginInfo) (int, domain.RespLogin)
	Register(registerInfo domain.RegisterInfo) (int, domain.RespRegister)
}
type MiddlewareService interface {
	AuthenticateUtil
	MiddlewareUtil
}

type AuthenticateUtil interface {
	Authenticate(authInfo domain.AuthInfo) (int, domain.RespAuth)
}

type MiddlewareUtil interface {
	GenBearerToken(hashBody domain.HashAuth) (domain.BearerToken, error)
	ValidateBearerToken(token domain.BearerToken) (domain.AuthBody, error)
}
