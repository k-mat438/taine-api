package usecase

import (
	"context"
	"taine-api/domain"
	"taine-api/models"
)

type OrganizationSvc interface {
	UpsertByExternalID(ctx context.Context, externalID, name string) error
	UpsertByExternalIDWithCreator(ctx context.Context, externalID, name, creatorSubID string) error
	SoftDeleteByExternalID(ctx context.Context, externalID string) error
}

type organizationSvc struct {
	orgRepository        domain.OrganizationRepository
	userRepository       domain.UserRepository
	membershipRepository domain.MembershipRepository
}

func NewOrganizationSvc(
	orgRepository domain.OrganizationRepository,
	userRepository domain.UserRepository,
	membershipRepository domain.MembershipRepository,
) OrganizationSvc {
	return &organizationSvc{
		orgRepository:        orgRepository,
		userRepository:       userRepository,
		membershipRepository: membershipRepository,
	}
}

func (s *organizationSvc) UpsertByExternalID(ctx context.Context, externalID, name string) error {
	_, err := s.orgRepository.UpsertByExternalID(ctx, externalID, name)
	return err
}

func (s *organizationSvc) UpsertByExternalIDWithCreator(ctx context.Context, externalID, name, creatorSubID string) error {
	// 組織を作成/更新
	org, err := s.orgRepository.UpsertByExternalID(ctx, externalID, name)
	if err != nil {
		return err
	}

	// 作成者のユーザーを取得
	if creatorSubID != "" {
		user, err := s.userRepository.GetUserBySubID(ctx, creatorSubID)
		if err != nil {
			return err
		}
		if user != nil {
			// 作成者をownerとして追加
			_, err = s.membershipRepository.UpsertByUserAndOrg(ctx, user.ID, org.ID, models.RoleOwner)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *organizationSvc) SoftDeleteByExternalID(ctx context.Context, externalID string) error {
	return s.orgRepository.SoftDeleteByExternalID(ctx, externalID)
}
