package usecase

import (
	"context"
	"taine-api/domain"
)

type MembershipSvc interface {
	UpsertByExternalIDs(ctx context.Context, userSubID, orgExternalID, role string) error
	SoftDeleteByExternalIDs(ctx context.Context, userSubID, orgExternalID string) error
}

type membershipSvc struct {
	membershipRepository domain.MembershipRepository
	userRepository       domain.UserRepository
	orgRepository        domain.OrganizationRepository
}

func NewMembershipSvc(membershipRepository domain.MembershipRepository, userRepository domain.UserRepository, orgRepository domain.OrganizationRepository) MembershipSvc {
	return &membershipSvc{
		membershipRepository: membershipRepository,
		userRepository:       userRepository,
		orgRepository:        orgRepository,
	}
}

func (s *membershipSvc) UpsertByExternalIDs(ctx context.Context, userSubID, orgExternalID, role string) error {
	// FindUserBySubID
	user, err := s.userRepository.GetUserBySubID(ctx, userSubID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	// FindOrgByExternalID
	org, err := s.orgRepository.FindByExternalID(ctx, orgExternalID)
	if err != nil {
		return err
	}
	if org == nil {
		return domain.ErrOrganizationNotFound
	}

	// Upsert membership
	_, err = s.membershipRepository.UpsertByUserAndOrg(ctx, user.ID, org.ID, role)
	return err
}

func (s *membershipSvc) SoftDeleteByExternalIDs(ctx context.Context, userSubID, orgExternalID string) error {
	// FindUserBySubID
	user, err := s.userRepository.GetUserBySubID(ctx, userSubID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	// FindOrgByExternalID
	org, err := s.orgRepository.FindByExternalID(ctx, orgExternalID)
	if err != nil {
		return err
	}
	if org == nil {
		return domain.ErrOrganizationNotFound
	}

	// Soft delete membership
	return s.membershipRepository.SoftDeleteByUserAndOrg(ctx, user.ID, org.ID)
}
