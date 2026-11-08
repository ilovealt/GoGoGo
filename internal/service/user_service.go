package service

import (
	"context"
	"strings"
	"time"

	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/pkg/utils"
	"gogogo/internal/repository"
)

type UserService struct {
	store repository.Store
}

func NewUserService(store repository.Store) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Get(ctx context.Context, id uint64) (model.User, error) {
	return s.store.GetUserByID(ctx, id)
}

func (s *UserService) UpdateProfile(ctx context.Context, id uint64, nickname string) (model.User, error) {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return model.User{}, apperrors.New(apperrors.ErrInvalid, "nickname is required")
	}
	if len(nickname) > 100 {
		return model.User{}, apperrors.New(apperrors.ErrInvalid, "nickname is too long")
	}
	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	user.Nickname = nickname
	if err := s.store.UpdateUser(ctx, user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uint64, oldPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return apperrors.New(apperrors.ErrInvalid, "new password must be at least 8 characters")
	}
	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if !utils.CheckPassword(user.PasswordHash, oldPassword) {
		return apperrors.New(apperrors.ErrUnauthorized, "current password is incorrect")
	}
	if oldPassword == newPassword {
		return apperrors.New(apperrors.ErrInvalid, "new password must be different from the current password")
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	user.TokenVersion++
	if err := s.store.UpdateUser(ctx, user); err != nil {
		return err
	}
	return s.store.RevokeUserRefreshTokens(ctx, id, time.Now().UTC())
}
