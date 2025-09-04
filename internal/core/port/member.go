package port

import (
	"Badminton-Hub/internal/core/domain"
	"context"
)

type MemberUtil interface {
	RegisterMember(registerForm domain.RegisterForm) (int, domain.ResponseRegisterMember)
	Login(loginForm domain.LoginForm) (int, domain.ResponseLogin)
	GoogleRegister(responseGoogle any) (int, domain.ResponseRegisterMember)
	GoogleLogin(responseGoogle any) (int, domain.ResponseLogin)
}

type RedirectUtil interface {
	GoogleLogin() (int, domain.ResponseRedirectGoogleLogin)
	GoogleRegister() (int, domain.ResponseRedirectGoogleRegister)
}

type MemberRepo interface {
	RegisterMember(ctx context.Context, member domain.Member) error
	FindEmailMember(ctx context.Context, email string) (domain.Member, error)
}
