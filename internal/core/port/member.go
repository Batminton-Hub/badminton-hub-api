package port

import (
	"Badminton-Hub/internal/core/domain"
	"context"
)

type MemberService interface {
	ProfileUtil
}

type ProfileUtil interface {
	GetProfile(userID domain.ReqGetProfile) (int, domain.RespGetProfile)
	UpdateProfile(userID domain.ReqGetProfile, request domain.ReqUpdateProfile) (int, domain.RespUpdateProfile)
}

type MemberRepo interface {
	SaveMember(ctx context.Context, member domain.Member) domain.ErrInfo
	FindEmailMember(ctx context.Context, email string) (domain.Member, domain.ErrInfo)
	GetMemberByUserID(ctx context.Context, userID string) (domain.Member, domain.ErrInfo)
	UpdateMember(ctx context.Context, userID string, request domain.ReqUpdateProfile) domain.ErrInfo
}
