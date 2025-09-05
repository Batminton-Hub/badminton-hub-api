package port

import (
	"Badminton-Hub/internal/core/domain"
	"context"
)

type MemberUtilGroup interface {
	MemberUtil
	ProfileUtil
}

type MemberUtil interface {
	RegisterMember(registerForm domain.RegisterForm) (int, domain.ResponseRegisterMember)
	Login(loginForm domain.LoginForm) (int, domain.ResponseLogin)
	GoogleRegister(responseGoogle any) (int, domain.ResponseRegisterMember)
	GoogleLogin(responseGoogle any) (int, domain.ResponseLogin)
}

type ProfileUtil interface {
	GetProfile(userID string) (int, domain.ResponseGetProfile)
}

type RedirectUtil interface {
	GoogleLogin() (int, domain.ResponseRedirectGoogleLogin)
	GoogleRegister() (int, domain.ResponseRedirectGoogleRegister)
}

type MemberRepo interface {
	SaveMember(ctx context.Context, member domain.Member) error
	FindEmailMember(ctx context.Context, email string) (domain.Member, error)
	GetMemberByUserID(ctx context.Context, userID string) (domain.Member, error)
}
