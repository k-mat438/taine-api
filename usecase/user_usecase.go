package usecase

import (
	"context"
	"taine-api/domain"

	"github.com/google/uuid"
)

type AuthClaims struct {
	SubID     string
	Name      string
	AvatarURL string
}

type UserService interface {
	SyncMe(ctx context.Context, claims *AuthClaims) (*domain.User, error)
	SoftDeleteBySubID(ctx context.Context, subID string) error
}

type UserUsecase interface {
	UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserBySubID(ctx context.Context, subID string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) UserUsecase {
	return &userUsecase{userRepository: userRepository}
}

func (u *userUsecase) UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return u.userRepository.UpsertUser(ctx, user)
}

func (u *userUsecase) GetUserBySubID(ctx context.Context, subID string) (*domain.User, error) {
	return u.userRepository.GetUserBySubID(ctx, subID)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.userRepository.GetUserByID(ctx, id)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.userRepository.DeleteUser(ctx, id)
}

// UserService implementation
type userService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) SyncMe(ctx context.Context, claims *AuthClaims) (*domain.User, error) {
	user := &domain.User{
		SubID:     claims.SubID,
		Name:      claims.Name,
		AvatarURL: claims.AvatarURL,
	}
	return s.userRepository.UpsertUser(ctx, user)
}

func (s *userService) SoftDeleteBySubID(ctx context.Context, subID string) error {
	return s.userRepository.SoftDeleteBySubID(ctx, subID)
}
