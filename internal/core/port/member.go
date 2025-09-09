package port

import (
	"Badminton-Hub/internal/core/domain"
	"context"
)

type MemberRepo interface {
	SaveMember(ctx context.Context, member domain.Member) error
	FindEmailMember(ctx context.Context, email string) (domain.Member, error)
	GetMemberByUserID(ctx context.Context, userID string) (domain.Member, error)
	UpdateMember(ctx context.Context, userID string, request domain.RequestUpdateProfile) error
}
