package core_util

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"time"
)

type MemberUtil struct {
	memberAdapter port.MemberAdapter
}

func NewMemberUtil(memberAdapter port.MemberAdapter) *MemberUtil {
	return &MemberUtil{
		memberAdapter: memberAdapter,
	}
}
func (m *MemberUtil) RegisterMember(registerForm domain.RegisterForm) error {
	member := domain.Member{
		Email:     registerForm.Email,
		Password:  util.HashPassword(registerForm.Password, "test"),
		Gender:    registerForm.Gender,
		Hash:      util.GenerateHash(registerForm.Email),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.memberAdapter.RegisterMember(member)
	return nil
}
