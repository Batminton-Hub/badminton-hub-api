package port

import "Badminton-Hub/internal/core/domain"

type MemberUtil interface {
	RegisterMember(registerForm domain.RegisterForm) error
}

type MemberAdapter interface {
	RegisterMember(member domain.Member) error
}
