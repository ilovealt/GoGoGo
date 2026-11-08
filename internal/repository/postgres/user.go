package postgres

import (
	"context"
	"strings"

	"gogogo/internal/model"
)

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func (s *Store) CreateUser(ctx context.Context, user *model.User) error {
	username := normalizeUsername(user.Username)
	if username == "" {
		return invalid("username is required")
	}
	if user.Role == "" {
		user.Role = model.RoleViewer
	}

	user.Username = username
	at := now()
	user.CreatedAt = at
	user.UpdatedAt = at

	err := s.db.WithContext(ctx).Create(user).Error
	return duplicateAs(err, "username already exists")
}

func (s *Store) GetUserByID(ctx context.Context, id uint64) (model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return model.User{}, recordNotFound(err, "user not found")
	}
	return user, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).
		Where("lower(username) = ?", normalizeUsername(username)).
		First(&user).Error
	if err != nil {
		return model.User{}, recordNotFound(err, "user not found")
	}
	return user, nil
}

func (s *Store) UpdateUser(ctx context.Context, user model.User) error {
	existing, err := s.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}

	username := normalizeUsername(user.Username)
	if username == "" {
		return invalid("username is required")
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.User{}).
		Where("lower(username) = ? AND id <> ?", username, user.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("username already exists")
	}

	user.Username = username
	user.CreatedAt = existing.CreatedAt
	user.UpdatedAt = now()

	err = s.db.WithContext(ctx).Save(&user).Error
	return duplicateAs(err, "username already exists")
}

func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
