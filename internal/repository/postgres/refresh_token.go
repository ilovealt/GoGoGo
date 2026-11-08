package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
)

func (s *Store) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	if token.CreatedAt.IsZero() {
		token.CreatedAt = now()
	}
	err := s.db.WithContext(ctx).Create(token).Error
	return duplicateAs(err, "refresh token already exists")
}

func (s *Store) GetRefreshTokenByHash(ctx context.Context, hash string) (model.RefreshToken, error) {
	var token model.RefreshToken
	err := s.db.WithContext(ctx).Where("token_hash = ?", hash).First(&token).Error
	if err != nil {
		return model.RefreshToken{}, recordNotFound(err, "refresh token not found")
	}
	return token, nil
}

// RotateRefreshToken 在同一个事务中锁定旧令牌、撤销旧令牌并写入新令牌。
func (s *Store) RotateRefreshToken(ctx context.Context, oldID uint64, next *model.RefreshToken) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var old model.RefreshToken
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&old, oldID).Error
		if err != nil {
			return recordNotFound(err, "refresh token not found")
		}
		if old.RevokedAt != nil {
			return apperrors.ErrTokenRevoked
		}

		at := now()
		if !old.ExpiresAt.After(at) {
			return apperrors.ErrTokenExpired
		}

		next.CreatedAt = at
		if err := tx.Create(next).Error; err != nil {
			return duplicateAs(err, "refresh token already exists")
		}

		return tx.Model(&model.RefreshToken{}).Where("id = ?", old.ID).Updates(map[string]any{
			"revoked_at":     at,
			"replaced_by_id": next.ID,
		}).Error
	})
}

func (s *Store) RevokeRefreshToken(ctx context.Context, id uint64, at time.Time) error {
	result := s.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("id = ?", id).
		UpdateColumn("revoked_at", gorm.Expr("COALESCE(revoked_at, ?)", at))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return notFound("refresh token not found")
	}
	return nil
}

func (s *Store) RevokeRefreshFamily(ctx context.Context, familyID string, at time.Time) error {
	return s.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("family_id = ? AND revoked_at IS NULL", familyID).
		UpdateColumn("revoked_at", at).Error
}

func (s *Store) RevokeUserRefreshTokens(ctx context.Context, userID uint64, at time.Time) error {
	return s.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		UpdateColumn("revoked_at", at).Error
}

func (s *Store) DeleteExpiredRefreshTokens(ctx context.Context, now time.Time) (int, error) {
	result := s.db.WithContext(ctx).Where("expires_at < ?", now).Delete(&model.RefreshToken{})
	return int(result.RowsAffected), result.Error
}
