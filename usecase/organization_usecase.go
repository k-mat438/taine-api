package usecase

import (
	"context"
	"taine-api/domain"
)

type OrganizationSvc interface {
	UpsertByExternalID(ctx context.Context, externalID, name string) error
	SoftDeleteByExternalID(ctx context.Context, externalID string) error
}

type organizationSvc struct {
	orgRepository domain.OrganizationRepository
}

func NewOrganizationSvc(orgRepository domain.OrganizationRepository) OrganizationSvc {
	return &organizationSvc{orgRepository: orgRepository}
}

func (s *organizationSvc) UpsertByExternalID(ctx context.Context, externalID, name string) error {
	_, err := s.orgRepository.UpsertByExternalID(ctx, externalID, name)
	return err
}

func (s *organizationSvc) SoftDeleteByExternalID(ctx context.Context, externalID string) error {
	return s.orgRepository.SoftDeleteByExternalID(ctx, externalID)
}
