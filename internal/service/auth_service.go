package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/pkg/utils"
	"gogogo/internal/repository"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthService struct {
	store      repository.Store
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(store repository.Store, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		store:      store,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, req request.Register) (model.User, error) {
	username := strings.ToLower(strings.TrimSpace(req.Username))
	if len(username) < 3 || len(username) > 64 {
		return model.User{}, apperrors.New(apperrors.ErrInvalid, "username must be 3-64 characters")
	}
	if len(req.Password) < 8 {
		return model.User{}, apperrors.New(apperrors.ErrInvalid, "password must be at least 8 characters")
	}
	role := normalizeRole(req.Role)
	if !validRole(role) {
		return model.User{}, apperrors.New(apperrors.ErrInvalid, "role must be admin, operator or viewer")
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		nickname = username
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return model.User{}, err
	}
	user := model.User{
		Username:     username,
		PasswordHash: hash,
		Nickname:     nickname,
		Role:         role,
		Status:       model.StatusEnabled,
	}
	if err := s.store.CreateUser(ctx, &user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, password, userAgent, ip string) (TokenPair, model.User, error) {
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return TokenPair{}, model.User{}, apperrors.New(apperrors.ErrUnauthorized, "invalid username or password")
	}
	if user.Status != model.StatusEnabled {
		return TokenPair{}, model.User{}, apperrors.New(apperrors.ErrForbidden, "user is disabled")
	}
	if !utils.CheckPassword(user.PasswordHash, password) {
		return TokenPair{}, model.User{}, apperrors.New(apperrors.ErrUnauthorized, "invalid username or password")
	}

	pair, err := s.issue(ctx, user, userAgent, ip, "")
	if err != nil {
		return TokenPair{}, model.User{}, err
	}
	return pair, user, nil
}

func (s *AuthService) Refresh(ctx context.Context, rawToken, userAgent, ip string) (TokenPair, error) {
	if strings.TrimSpace(rawToken) == "" {
		return TokenPair{}, apperrors.New(apperrors.ErrUnauthorized, "refresh token is required")
	}
	hash := utils.SHA256Hex(rawToken)
	old, err := s.store.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return TokenPair{}, apperrors.New(apperrors.ErrUnauthorized, "invalid refresh token")
	}
	if old.RevokedAt != nil {
		_ = s.store.RevokeRefreshFamily(ctx, old.FamilyID, time.Now().UTC())
		return TokenPair{}, apperrors.New(apperrors.ErrUnauthorized, "refresh token has been revoked")
	}
	if !old.ExpiresAt.After(time.Now().UTC()) {
		return TokenPair{}, apperrors.New(apperrors.ErrUnauthorized, "refresh token has expired")
	}

	user, err := s.store.GetUserByID(ctx, old.UserID)
	if err != nil || user.Status != model.StatusEnabled {
		return TokenPair{}, apperrors.New(apperrors.ErrUnauthorized, "user is not available")
	}

	raw, err := utils.NewOpaqueToken()
	if err != nil {
		return TokenPair{}, err
	}
	newHash := utils.SHA256Hex(raw)
	next := model.RefreshToken{
		UserID:    user.ID,
		TokenHash: newHash,
		FamilyID:  old.FamilyID,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL),
		UserAgent: userAgent,
		IP:        ip,
	}
	if err := s.store.RotateRefreshToken(ctx, old.ID, &next); err != nil {
		if errors.Is(err, apperrors.ErrTokenRevoked) {
			_ = s.store.RevokeRefreshFamily(ctx, old.FamilyID, time.Now().UTC())
		}
		return TokenPair{}, apperrors.New(apperrors.ErrUnauthorized, "invalid refresh token")
	}

	accessToken, expiresAt, err := utils.GenerateAccessToken(s.jwtSecret, s.accessTTL, user)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: raw,
		TokenType:    "Bearer",
		ExpiresIn:    int64(expiresAt.Sub(time.Now().UTC()).Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	if strings.TrimSpace(rawToken) == "" {
		return nil
	}
	token, err := s.store.GetRefreshTokenByHash(ctx, utils.SHA256Hex(rawToken))
	if err != nil {
		return nil
	}
	return s.store.RevokeRefreshToken(ctx, token.ID, time.Now().UTC())
}

func (s *AuthService) issue(ctx context.Context, user model.User, userAgent, ip, familyID string) (TokenPair, error) {
	accessToken, expiresAt, err := utils.GenerateAccessToken(s.jwtSecret, s.accessTTL, user)
	if err != nil {
		return TokenPair{}, err
	}
	raw, err := utils.NewOpaqueToken()
	if err != nil {
		return TokenPair{}, err
	}
	hash := utils.SHA256Hex(raw)
	if familyID == "" {
		familyID = utils.NewUUID()
	}
	token := model.RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		FamilyID:  familyID,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL),
		UserAgent: userAgent,
		IP:        ip,
	}
	if err := s.store.CreateRefreshToken(ctx, &token); err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: raw,
		TokenType:    "Bearer",
		ExpiresIn:    int64(expiresAt.Sub(time.Now().UTC()).Seconds()),
	}, nil
}
