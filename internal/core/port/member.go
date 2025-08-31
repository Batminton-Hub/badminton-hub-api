package port

import (
	"Badminton-Hub/internal/core/domain"
	"context"
)

type MemberUtil interface {
	RegisterMember(registerForm domain.RegisterForm) (int, domain.ResponseRegisterMember)
	Login(loginForm domain.LoginForm) (int, domain.ResponseLogin)
}

type MemberRepo interface {
	RegisterMember(ctx context.Context, member domain.Member) error
	FindEmailMember(ctx context.Context, loginForm domain.LoginForm) (domain.Member, error)
}
