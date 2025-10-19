package usecase

import (
	"context"
	"errors"
	"taine-api/domain"

	"github.com/google/uuid"
)

type WishSvc interface {
	CreateWish(ctx context.Context, organizationID uuid.UUID, title, note string, orderNo int) (*domain.Wish, error)
	CreateWishByOrganizationExternalID(ctx context.Context, externalID, title, note string, orderNo int) (*domain.Wish, error)
	GetWish(ctx context.Context, id uuid.UUID) (*domain.Wish, error)
	GetWishesByOrganization(ctx context.Context, organizationID uuid.UUID) ([]*domain.Wish, error)
	GetWishesByOrganizationExternalID(ctx context.Context, externalID string) ([]*domain.Wish, error)
	UpdateWish(ctx context.Context, id uuid.UUID, title, note string, orderNo int) (*domain.Wish, error)
	DeleteWish(ctx context.Context, id uuid.UUID) error
	SoftDeleteWish(ctx context.Context, id uuid.UUID) error
	RestoreWish(ctx context.Context, id uuid.UUID) error
	UpdateWishOrder(ctx context.Context, id uuid.UUID, orderNo int) error
}

type wishSvc struct {
	wishRepository domain.WishRepository
	orgRepository  domain.OrganizationRepository
}

func NewWishSvc(
	wishRepository domain.WishRepository,
	orgRepository domain.OrganizationRepository,
) WishSvc {
	return &wishSvc{
		wishRepository: wishRepository,
		orgRepository:  orgRepository,
	}
}

func (s *wishSvc) CreateWish(ctx context.Context, organizationID uuid.UUID, title, note string, orderNo int) (*domain.Wish, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	// 組織の存在確認
	org, err := s.orgRepository.FindByExternalID(ctx, organizationID.String())
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errors.New("organization not found")
	}

	wish := &domain.Wish{
		OrganizationID: organizationID,
		Title:          title,
		Note:           note,
		OrderNo:        orderNo,
	}

	return s.wishRepository.Create(ctx, wish)
}

func (s *wishSvc) CreateWishByOrganizationExternalID(ctx context.Context, externalID, title, note string, orderNo int) (*domain.Wish, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	// external_idから組織を取得
	org, err := s.orgRepository.FindByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}

	wish := &domain.Wish{
		OrganizationID: org.ID,
		Title:          title,
		Note:           note,
		OrderNo:        orderNo,
	}

	return s.wishRepository.Create(ctx, wish)
}

func (s *wishSvc) GetWish(ctx context.Context, id uuid.UUID) (*domain.Wish, error) {
	wish, err := s.wishRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if wish == nil {
		return nil, errors.New("wish not found")
	}
	return wish, nil
}

func (s *wishSvc) GetWishesByOrganization(ctx context.Context, organizationID uuid.UUID) ([]*domain.Wish, error) {
	return s.wishRepository.FindByOrganizationID(ctx, organizationID)
}

func (s *wishSvc) GetWishesByOrganizationExternalID(ctx context.Context, externalID string) ([]*domain.Wish, error) {
	// external_idから組織を取得
	org, err := s.orgRepository.FindByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errors.New("organization not found")
	}

	// 組織のIDでWishを取得
	return s.wishRepository.FindByOrganizationID(ctx, org.ID)
}

func (s *wishSvc) UpdateWish(ctx context.Context, id uuid.UUID, title, note string, orderNo int) (*domain.Wish, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	// 既存のWishを取得
	wish, err := s.wishRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if wish == nil {
		return nil, errors.New("wish not found")
	}

	// 更新
	wish.Title = title
	wish.Note = note
	wish.OrderNo = orderNo

	return s.wishRepository.Update(ctx, wish)
}

func (s *wishSvc) DeleteWish(ctx context.Context, id uuid.UUID) error {
	// 存在確認
	wish, err := s.wishRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if wish == nil {
		return errors.New("wish not found")
	}

	return s.wishRepository.Delete(ctx, id)
}

func (s *wishSvc) SoftDeleteWish(ctx context.Context, id uuid.UUID) error {
	// 存在確認
	wish, err := s.wishRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if wish == nil {
		return errors.New("wish not found")
	}
	if wish.DeletedAt != nil {
		return errors.New("wish is already deleted")
	}

	return s.wishRepository.SoftDelete(ctx, id)
}

func (s *wishSvc) RestoreWish(ctx context.Context, id uuid.UUID) error {
	// 存在確認（削除済みも含めて検索する必要があるため、別途実装が必要かもしれません）
	wish, err := s.wishRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if wish == nil {
		return errors.New("wish not found")
	}
	if wish.DeletedAt == nil {
		return errors.New("wish is not deleted")
	}

	return s.wishRepository.Restore(ctx, id)
}

func (s *wishSvc) UpdateWishOrder(ctx context.Context, id uuid.UUID, orderNo int) error {
	// 存在確認
	wish, err := s.wishRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if wish == nil {
		return errors.New("wish not found")
	}

	return s.wishRepository.UpdateOrder(ctx, id, orderNo)
}
