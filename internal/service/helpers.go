package service

import (
	"strings"

	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
)

func valueOr[T any](value *T, fallback T) T {
	if value != nil {
		return *value
	}
	return fallback
}

func stringOr(value *string, fallback string) string {
	if value != nil {
		return *value
	}
	return fallback
}

func validRole(role string) bool {
	switch role {
	case model.RoleAdmin, model.RoleOperator, model.RoleViewer:
		return true
	default:
		return false
	}
}

func normalizeRole(role string) string {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		return model.RoleViewer
	}
	return role
}

func requireStatus(status uint8) error {
	if status != model.StatusDisabled && status != model.StatusEnabled {
		return apperrors.New(apperrors.ErrInvalid, "status must be 0 or 1")
	}
	return nil
}
